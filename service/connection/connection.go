package connection

import (
	"net"
	"net/url"

	"github.com/google/uuid"
)

type client struct {
	id uuid.UUID
}

func (c *client) Close() error {
	return nil
}

type Client interface {
	Close() error
}

type ConnectionService interface {
	Attach(address url.URL, handler net.Conn) (Client, error)
}

type connectionService struct {
	clients []client
}

func NewConnectionSerivce() ConnectionService {
	return &connectionService{}
}

func (cs *connectionService) Attach(address url.URL, handler net.Conn) (Client, error) {
	return &client{}, nil
}
