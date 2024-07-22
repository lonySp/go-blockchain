package network

import (
	"fmt"
	"time"
)

type ServerOpts struct {
	Transport []Transport
}

type Server struct {
	ServerOpts
	rpcCh  chan RPC
	quitCh chan struct{}
}

func NewServer(opts ServerOpts) *Server {
	return &Server{
		ServerOpts: opts,
		rpcCh:      make(chan RPC),
	}
}

func (s *Server) Start() {
	s.initTransport()
	ticker := time.NewTicker(5 * time.Second)

free:
	for {
		select {
		case rpc := <-s.rpcCh:
			fmt.Printf("Got RPC: %v\n", rpc)
		case <-s.quitCh:
			break free
		case <-ticker.C:
			fmt.Println("Tick")
		}
	}
	fmt.Println("Server shutdown")
}

func (s *Server) initTransport() {
	for _, tr := range s.Transport {
		go func(tr Transport) {
			for rpc := range tr.Consume() {
				// handle
				s.rpcCh <- rpc
			}
		}(tr)
	}
}
