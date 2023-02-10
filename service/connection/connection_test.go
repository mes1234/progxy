package connection_test

import (
	"testing"

	"github.com/mes1234/progxy/internal/dto"
	"github.com/mes1234/progxy/service/connection"
)

func TestConnectionServiceReturnsClientOnConnection(t *testing.T) {
	// Arrange
	cs := connection.NewConnectionSerivce()
	// Act
	_, err := cs.Attach(dto.Proxied{
		Host: "localhost",
		Port: 1234,
	})
	// Assert
	if err != nil {
		t.Fatalf("Error during attaching client: %v", err)
	}

}
