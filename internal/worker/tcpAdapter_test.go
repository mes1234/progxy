package worker_test

import (
	"net"
	"testing"
	"time"

	"github.com/mes1234/progxy/internal/dto"
	"github.com/mes1234/progxy/internal/worker"
)

func TestTcpAdapterShallShuffleDataFromClientToProxyTest(t *testing.T) {

	//Arrange
	_ = worker.NewTcpAdaper(dto.Proxied{}, 1111, dummylistenFunc, dummyDialFunc, dummydnsLookupFunc)
	//Act

	//Assert
	//TODO do some assertions
	time.Sleep(100 * time.Second)
}

func dummylistenFunc(network, address string) (net.Listener, error) {
	return &dummyListner{
		amount: 1,
	}, nil
}

func dummydnsLookupFunc(host string) ([]net.IP, error) {
	ips := make([]net.IP, 1)
	ips = append(ips, net.IPv4(0xFF, 0xFF, 0xFF, 0xFF))

	return ips, nil
}

func dummyDialFunc(network string, raddr string) (net.Conn, error) {

	ticker := time.NewTicker(10000 * time.Millisecond)

	return &dummyConn{
		ticker: *ticker,
	}, nil
}

type dummyListner struct {
	amount int
}

type dummyConn struct {
	ticker time.Ticker
}

func (dc *dummyConn) Read(b []byte) (n int, err error) {

	<-dc.ticker.C
	for i := range b {
		b[i] = 0xFF
	}

	return len(b), nil
}

func (dc *dummyConn) Write(b []byte) (n int, err error) {
	for i := range b {
		b[i] = 0xFF
	}

	return len(b), nil
}

func (dc *dummyConn) Close() error {
	panic("Not Needed")
}

func (dc *dummyConn) LocalAddr() net.Addr {
	panic("Not Needed")
}

func (dc *dummyConn) RemoteAddr() net.Addr {
	panic("Not Needed")
}

func (dc *dummyConn) SetDeadline(t time.Time) error {
	panic("Not Needed")
}

func (dc *dummyConn) SetReadDeadline(t time.Time) error {
	panic("Not Needed")
}

func (dc *dummyConn) SetWriteDeadline(t time.Time) error {
	panic("Not Needed")
}

func (dl *dummyListner) Accept() (net.Conn, error) {
	if dl.amount > 0 {

		dl.amount--
		return &dummyConn{}, nil
	}
	time.Sleep(1000 * time.Second)
	return &dummyConn{}, nil
}

func (dl *dummyListner) Addr() net.Addr {
	panic("Not Needed")
}

func (dl *dummyListner) Close() error {
	panic("Not Needed")
}
