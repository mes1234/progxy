package connection

import (
	"github.com/mes1234/progxy/internal/dto"
)

// ConnectionService
type ConnectionService interface {
	Attach(address dto.Proxied) (dto.Client, error)
}

type connectionService struct {
	clients []dto.Client
}

func NewConnectionSerivce() ConnectionService {
	return &connectionService{}
}

func (cs *connectionService) Attach(address dto.Proxied) (dto.Client, error) {
	return dto.NewClient(), nil
}
