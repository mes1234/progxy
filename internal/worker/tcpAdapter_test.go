package worker_test

import (
	"net"
	"testing"
	"time"

	"github.com/mes1234/progxy/internal/dto"
	"github.com/mes1234/progxy/internal/worker"
	"github.com/sirupsen/logrus"
)

func TestTcpAdapterShallShuffleDataFromProxiedToClient(t *testing.T) {

	//Arrange
	// Buffer which shall be written from client to proxied
	bufOut := make([]byte, 1024)

	valueToSendFromProxiedToClient := byte(0xAA)

	_ = worker.NewTcpAdaper(dto.Proxied{}, 1111, generateListenFuncWithBuf(&bufOut), generateDialFuncWithValue(valueToSendFromProxiedToClient), dummydnsLookupFunc, logrus.New())

	//Act
	// Allow to shuffling start
	time.Sleep(1 * time.Second)

	//Assert
	for i := range bufOut {
		if bufOut[i] != valueToSendFromProxiedToClient {
			t.Fatalf("Expected data %v got %v", valueToSendFromProxiedToClient, bufOut[i])
		}
	}
}

func generateListenFuncWithBuf(bufOut *[]byte) worker.ListenFunc {
	return func(network, address string) (net.Listener, error) {
		return &dummyListner{
			amount: 1,
			bufOut: *bufOut,
		}, nil
	}
}

func dummydnsLookupFunc(host string) ([]net.IP, error) {
	ips := make([]net.IP, 1)
	ips = append(ips, net.IPv4(0xFF, 0xFF, 0xFF, 0xFF))

	return ips, nil
}

func generateDialFuncWithValue(value byte) worker.DialFunc {
	return func(network string, raddr string) (net.Conn, error) {

		ticker := time.NewTicker(1 * time.Millisecond)

		return &dummyConn{
			ticker: *ticker,
			value:  value,
		}, nil
	}
}

type dummyListner struct {
	amount int
	bufOut []byte
	value  byte
}

type dummyConn struct {
	bufOut []byte
	value  byte
	ticker time.Ticker
}

func (dc *dummyConn) Read(b []byte) (n int, err error) {

	<-dc.ticker.C
	for i := range b {
		b[i] = dc.value
	}

	return len(b), nil
}

func (dc *dummyConn) Write(b []byte) (n int, err error) {
	copy(dc.bufOut, b)

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
		return &dummyConn{
			bufOut: dl.bufOut,
			value:  dl.value,
		}, nil
	}

	// When amount is exceeded wait long for next accept to not overhelm handler
	time.Sleep(1000 * time.Second)
	return &dummyConn{}, nil
}

func (dl *dummyListner) Addr() net.Addr {
	panic("Not Needed")
}

func (dl *dummyListner) Close() error {
	panic("Not Needed")
}
