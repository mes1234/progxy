package worker

import (
	"context"
	"fmt"
	"io"
	"net"

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

		clientInChan, clientOutChan := createChannelFromReaderWriter(clientConn)

		proxiedInChan, proxiedOutChan := createChannelFromReaderWriter(proxiedConn)

		clientShuffler, _ := NewShuffler(clientOutChan, ctx)

		clientShuffler.Attach(CreateWriteToChannelProcessorFunc(proxiedInChan))

		proxiedShuffler, _ := NewShuffler(proxiedOutChan, ctx)

		proxiedShuffler.Attach(CreateWriteToChannelProcessorFunc(clientInChan))
	}
}

// Should be used as goroutine otherwise it will never release thread
func readAndForward(out chan<- []byte, reader io.Reader) {
	buf := make([]byte, bufferSize)
	for {
		n := 0
		for n == 0 {
			n, _ = reader.Read(buf)
			//TODO error handling
		}
		out <- buf
	}
}

// Should be used as goroutine otherwise it will never release thread
func ForwardToWriter(in <-chan []byte, writer io.Writer) {
	for {
		data := <-in
		writer.Write(data)
	}
}

func createChannelFromReaderWriter(rw io.ReadWriter) (in chan []byte, out chan []byte) {

	out = make(chan []byte, bufferSize)
	in = make(chan []byte, bufferSize)

	// Read data from Reader and push to channel
	go readAndForward(out, rw)

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
