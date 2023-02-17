package worker

import (
	"context"
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/mes1234/progxy/internal/dto"
)

const bufferSize = 1024

// TcpAdapter listen to connections and
type TcpAdapter interface {
}

type tcpAdapter struct {
	listner  net.Listener
	dialFunc DialFunc
}

type DialFunc func(network string, raddr string) (net.Conn, error)
type ListenFunc func(network, address string) (net.Listener, error)
type DnsLookupFunc func(host string) ([]net.IP, error)

func (tcpA *tcpAdapter) start(clientsListner net.Listener, proxiedAddr string) {

	tcpA.listner = clientsListner

	tcpA.handle(proxiedAddr)
}

func (tcpA *tcpAdapter) handle(proxiedAddr string) {
	for {
		// create context to close this proxy
		ctx := context.Background()

		// connection to talk with client
		clientConn, _ := tcpA.listner.Accept()

		// connection to talk with proxied
		proxiedConn, _ := tcpA.dialFunc("tcp", proxiedAddr)

		clientWg := new(sync.WaitGroup)
		clientWg.Add(1)

		// clientInChan allow to write to client
		// clientOutChan gets data from client
		clientInChan, clientOutChan := createChannelFromReaderWriter(clientConn, clientWg)

		proxiedWg := new(sync.WaitGroup)
		proxiedWg.Add(1)

		// proxiedInChan allow to write to proxied
		// proxiedOutChan gets data from proxied
		proxiedInChan, proxiedOutChan := createChannelFromReaderWriter(proxiedConn, proxiedWg)

		//Shuffler which will process data from client -> proxied
		clientShuffler, _ := NewShuffler(clientOutChan, ctx)

		// Shuffler which will process data from proxied -> client
		proxiedShuffler, _ := NewShuffler(proxiedOutChan, ctx)

		// Pass data from client to proxied
		clientShuffler.Attach(CreateWriteToConsoleProcessorFunc("client -> proxied"))
		clientShuffler.Attach(CreateWriteToChannelProcessorFunc(proxiedInChan))

		// Pass data from proxied to client
		proxiedShuffler.Attach(CreateWriteToConsoleProcessorFunc("proxied -> client"))
		proxiedShuffler.Attach(CreateWriteToChannelProcessorFunc(clientInChan))

		//go WaitToClose("client", clientWg, clientConn)
		//go WaitToClose("proxied", proxiedWg, proxiedConn)
	}
}

// Should be run as goroutine otherwise will block
func WaitToClose(who string, waiter *sync.WaitGroup, client net.Conn) {
	waiter.Wait()
	fmt.Printf("Closed connection by %v\n", who)
	client.Close()
}

// Should be used as goroutine otherwise it will never release thread
func readAndForward(out chan<- []byte, reader io.Reader, wg *sync.WaitGroup) {
	defer wg.Done()
	readBuf := make([]byte, bufferSize)
	for {

		n, err := reader.Read(readBuf)
		if err != nil {
			return
		}
		if n != 0 {
			outBuf := make([]byte, n)
			copy(outBuf, readBuf)
			out <- outBuf
		}

	}
}

// Should be used as goroutine otherwise it will never release thread
func ForwardToWriter(in <-chan []byte, writer io.Writer) {
	for {
		data := <-in
		n, err := writer.Write(data)
		if err != nil {
			fmt.Printf("read %v data and failed", n)
			return
		}
	}
}

func createChannelFromReaderWriter(rw io.ReadWriter, wg *sync.WaitGroup) (in chan []byte, out chan []byte) {

	out = make(chan []byte, 1024)
	in = make(chan []byte, 1024)

	// Read data from Reader and push to channel
	go readAndForward(out, rw, wg)

	// Forward data from in channel to Writer
	go ForwardToWriter(in, rw)

	return in, out
}

func NewTcpAdaper(
	proxied dto.Proxied,
	port int,
	listenFunc ListenFunc,
	dialFunc DialFunc,
	dnsLookupFunc DnsLookupFunc) TcpAdapter {

	// based on provided proxied address and port get address to proxied service
	proxiedAddr := getProxiedAddr(proxied, dnsLookupFunc)

	// start accepting clients on provided common port
	listen, _ := listenFunc("tcp", "localhost"+":"+fmt.Sprint(port))

	// create tcpAdapter which will setup pipelines for clients
	adapter := tcpAdapter{
		dialFunc: dialFunc,
	}

	// bootstrap bounding and shuffling data between client and proxied
	go adapter.start(listen, proxiedAddr)

	return &adapter
}

func getProxiedAddr(proxied dto.Proxied, dnsLookupFunc DnsLookupFunc) string {

	// get IP for HOST
	proxiedIp, _ := dnsLookupFunc(proxied.Host)

	// create IPAddr based on first found match
	addr := net.IPAddr{
		IP: proxiedIp[0],
	}
	return addr.IP.String() + ":" + fmt.Sprint(proxied.Port)

}
