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
	listner net.Listener
	proxied net.Conn
}

func (tcpA *tcpAdapter) start(clientsListner net.Listener, proxiedAddr *net.TCPAddr) {

	tcpA.listner = clientsListner

	for {
		// create context to close this proxy
		ctx := context.Background()

		// connection to talk with client
		clientConn, _ := tcpA.listner.Accept()

		// connection to talk with proxied
		proxiedConn, _ := net.DialTCP("tcp", nil, proxiedAddr)

		clientInChan, clientOutChan := createChannelFromReaderWriter(clientConn)

		proxiedInChan, proxiedOutChan := createChannelFromReaderWriter(proxiedConn)

		clientShuffler, _ := NewShuffler(clientOutChan, ctx)

		clientShuffler.Attach(CreateWriteToChannelProcessorFunc(proxiedInChan))

		proxiedShuffler, _ := NewShuffler(proxiedOutChan, ctx)

		proxiedShuffler.Attach(CreateWriteToChannelProcessorFunc(clientInChan))
	}
}

func createChannelFromReaderWriter(rw io.ReadWriter) (in chan []byte, out chan []byte) {
	return make(chan []byte), make(chan []byte)
}

func NewTcpAdaper(proxied dto.Proxied, port int) TcpAdapter {

	// based on provided proxied address and port get address to proxied service
	proxiedAddr := GetProxiedAddr(proxied)

	// start accepting clients on provided common port
	listen, _ := net.Listen("tcp", "localhost"+":"+fmt.Sprint(port))

	// create tcpAdapter which will setup pipelines for clients
	adapter := tcpAdapter{}

	// bootstrap bounding and shuffling data between client and proxied
	go adapter.start(listen, proxiedAddr)

	return &adapter
}

func GetProxiedAddr(proxied dto.Proxied) *net.TCPAddr {

	// get IP for HOST
	proxiedIp, _ := net.LookupIP(proxied.Host)

	// create IPAddr based on first found match
	addr := net.IPAddr{
		IP: proxiedIp[0],
	}
	return &net.TCPAddr{IP: addr.IP, Port: proxied.Port}

}
