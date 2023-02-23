package worker

import (
	"context"
	"fmt"
	"net"
	"sync"

	"github.com/mes1234/progxy/internal/dto"
	"github.com/sirupsen/logrus"
)

const bufferSize = 1024

// TcpAdapter listen to connections and
type TcpAdapter interface {
}

type tcpAdapter struct {
	listner       net.Listener
	dialFunc      DialFunc
	clientCounter int
	logger        *logrus.Logger
}

type DialFunc func(network string, raddr string) (net.Conn, error)
type ListenFunc func(network, address string) (net.Listener, error)
type DnsLookupFunc func(host string) ([]net.IP, error)

func (tcpA *tcpAdapter) start(clientsListner net.Listener, proxiedAddr string, port int, ctx context.Context) {

	tcpA.listner = clientsListner

	tcpA.logger.WithFields(logrus.Fields{
		"proxied": proxiedAddr,
		"port":    port,
	}).Info("tcp Adapter started")

	tcpA.handle(proxiedAddr, ctx)

}

func (tcpA *tcpAdapter) handle(proxiedAddr string, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			{
				// connection to talk with client
				clientConn, err := tcpA.listner.Accept()

				if err == nil {
					tcpA.clientCounter++
					tcpA.logger.WithField("counter", tcpA.clientCounter).Info("started handling client")
					go start(tcpA, proxiedAddr, clientConn, ctx)
				}

			}
		}
	}
}

func start(tcpA *tcpAdapter, proxiedAddr string, clientConn net.Conn, ctx context.Context) {

	cancallableContext, cancelFunc := context.WithCancel(ctx)
	logger := tcpA.logger

	// connection to talk with proxied
	proxiedConn, _ := tcpA.dialFunc("tcp", proxiedAddr)

	// clientInChan allow to write to client
	// clientOutChan gets data from client
	clientInChan, clientOutChan, clientWg := CreateChannelFromReaderWriter("client", clientConn, logger, cancallableContext)

	// proxiedInChan allow to write to proxied
	// proxiedOutChan gets data from proxied
	proxiedInChan, proxiedOutChan, proxiedWg := CreateChannelFromReaderWriter("proxied", proxiedConn, logger, cancallableContext)

	//Shuffler which will process data from client -> proxied
	clientShuffler := NewShuffler(clientOutChan, cancallableContext)

	// Shuffler which will process data from proxied -> client
	proxiedShuffler := NewShuffler(proxiedOutChan, cancallableContext)

	// Pass data from client to proxied
	clientShuffler.Attach(CreateWriteToConsoleProcessorFunc("client -> proxied", logger))
	clientShuffler.Attach(CreateWriteToChannelProcessorFunc(proxiedInChan))

	// Pass data from proxied to client
	proxiedShuffler.Attach(CreateWriteToConsoleProcessorFunc("proxied -> client", logger))
	//proxiedShuffler.Attach(CreateMuddlingProcessorFunc())
	proxiedShuffler.Attach(CreateWriteToChannelProcessorFunc(clientInChan))

	go WaitToClose("client", clientWg, clientConn, logger, cancelFunc)
	go WaitToClose("proxied", proxiedWg, proxiedConn, logger, cancelFunc)
}

// Should be run as goroutine otherwise will block
func WaitToClose(who string, waiter *sync.WaitGroup, conn net.Conn, logger *logrus.Logger, cancelFunc context.CancelFunc) {
	waiter.Wait()
	logger.WithField("who", who).Info("Closed connection")
	if conn != nil {
		conn.Close()
	}
	cancelFunc()
}

func NewTcpAdaper(
	proxied dto.Proxied,
	port int,
	listenFunc ListenFunc,
	dialFunc DialFunc,
	dnsLookupFunc DnsLookupFunc,
	logger *logrus.Logger,
	ctx context.Context) TcpAdapter {

	// based on provided proxied address and port get address to proxied service
	proxiedAddr := getProxiedAddr(proxied, dnsLookupFunc)

	// start accepting clients on provided common port
	listen, _ := listenFunc("tcp", "localhost"+":"+fmt.Sprint(port))

	// create tcpAdapter which will setup pipelines for clients
	adapter := tcpAdapter{
		dialFunc:      dialFunc,
		clientCounter: 0,
		logger:        logger,
	}

	// bootstrap bounding and shuffling data between client and proxied
	go adapter.start(listen, proxiedAddr, port, ctx)

	// wait to close listner when context closes
	go waitToClose(proxied, port, listen, logger, ctx)

	return &adapter
}

func waitToClose(proxied dto.Proxied, port int, listen net.Listener, logger *logrus.Logger, ctx context.Context) {
	<-ctx.Done()
	if listen != nil {
		listen.Close()
	}

	logger.WithFields(logrus.Fields{
		"proxied": fmt.Sprintf("%v:%v", proxied.Host, proxied.Port),
		"port":    port,
	}).Info("tcp Adapter finished")
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
