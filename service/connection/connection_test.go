package connection_test

import (
	"net"
	"net/url"
	"testing"

	"github.com/mes1234/progxy/service/connection"
)

func TestConnectionServiceReturnsClientOnConnection(t *testing.T) {
	// Arrange
	cs := connection.NewConnectionSerivce()
	// Act
	_, err := cs.Attach(url.URL{
		Host: "localhost:5002",
	}, &net.IPConn{})
	// Assert
	if err != nil {
		t.Fatalf("Error during attaching client: %v", err)
	}

}
