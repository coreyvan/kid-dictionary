package transport

import (
	"context"
	"fmt"
)

type Server interface {
	Listen(addr string) error
	Ready() chan struct{}
	Shutdown(ctx context.Context) error
	Addr() (string, error)
}

type server struct {
	addr      string
	readyChan chan struct{}
}

func NewServer(bindAddr, port string) Server {
	return &server{
		addr:      fmt.Sprintf("%s:%s", bindAddr, port),
		readyChan: make(chan struct{}),
	}
}

func (s server) Listen(addr string) error {
	close(s.readyChan)
	return fmt.Errorf("done listening")
}

func (s server) Ready() chan struct{} {
	return s.readyChan
}

func (s server) Shutdown(ctx context.Context) error {
	return nil
}

func (s server) Addr() (string, error) {
	return s.addr, nil
}
