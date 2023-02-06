package connection

import (
	"net"
	"net/url"

	"github.com/mes1234/progxy/internal/dto"
)

type ConnectionService interface {
	Attach(address url.URL, handler net.Conn) (dto.Client, error)
}

type connectionService struct {
	clients []dto.Client
}

func NewConnectionSerivce() ConnectionService {
	return &connectionService{}
}

func (cs *connectionService) Attach(address url.URL, handler net.Conn) (dto.Client, error) {
	return dto.NewClient(), nil
}
