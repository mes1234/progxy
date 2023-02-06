package dto

import "github.com/google/uuid"

type client struct {
	Id uuid.UUID
}

type Client interface {
	Close() error
}

func (c *client) Close() error {
	return nil
}

func NewClient() Client {
	return &client{}
}
