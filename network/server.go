package network

import (
	"fmt"
	"time"
)

// ServerOpts 结构体包含服务器的传输选项
// ServerOpts struct contains server transport options
type ServerOpts struct {
	Transport []Transport
}

// Server 结构体表示服务器
// Server struct represents the server
type Server struct {
	ServerOpts
	rpcCh  chan RPC
	quitCh chan struct{}
}

// NewServer 创建并返回一个新的 Server 实例
// NewServer creates and returns a new Server instance
func NewServer(opts ServerOpts) *Server {
	return &Server{
		ServerOpts: opts,
		rpcCh:      make(chan RPC),
	}
}

// Start 方法启动服务器
// Start method starts the server
func (s *Server) Start() {
	s.initTransport()
	ticker := time.NewTicker(5 * time.Second)

free:
	for {
		select {
		case rpc := <-s.rpcCh:
			fmt.Printf("%+v\n", rpc)
		case <-s.quitCh:
			break free
		case <-ticker.C:
			fmt.Println("Tick")
		}
	}
	fmt.Println("Server shutdown")
}

// initTransport 方法初始化所有的传输选项
// initTransport method initializes all transport options
func (s *Server) initTransport() {
	for _, tr := range s.Transport {
		go func(tr Transport) {
			for rpc := range tr.Consume() {
				// 处理 RPC 消息
				// handle RPC messages
				s.rpcCh <- rpc
			}
		}(tr)
	}
}
