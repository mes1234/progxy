package worker

import (
	"context"
	"fmt"
	"io"
	"net"

	"github.com/mes1234/progxy/internal/dto"
)

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

func createChannelFromReaderWriter(rw io.ReadWriter) (in chan []byte, out chan []byte) {

	out = make(chan []byte)
	in = make(chan []byte)

	go func() {
		buf := make([]byte, 2)
		for {
			n := 0
			for n == 0 {
				n, err := rw.Read(buf)

				if err != nil {
					panic(n)
				}
			}

			out <- buf
		}
	}()

	go func() {
		for {
			select {
			case data := <-in:
				rw.Write(data)
			}
		}
	}()

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
