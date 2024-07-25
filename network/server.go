package network

import (
	"fmt"
	"github.com/lonySp/go-blockchain/core"
	"github.com/lonySp/go-blockchain/crypto"
	"github.com/sirupsen/logrus"
	"time"
)

var defaultBlockTime = 5 * time.Second

// ServerOpts 结构体包含服务器的传输选项
// ServerOpts struct contains server transport options
type ServerOpts struct {
	Transport  []Transport
	BlockTime  time.Duration
	PrivateKey *crypto.PrivateKey
}

// Server 结构体表示服务器
// Server struct represents the server
type Server struct {
	ServerOpts
	blockTime   time.Duration
	memPool     *TxPool
	isValidator bool
	rpcCh       chan RPC
	quitCh      chan struct{}
}

// NewServer 创建并返回一个新的 Server 实例
// NewServer creates and returns a new Server instance
func NewServer(opts ServerOpts) *Server {
	if opts.BlockTime == time.Duration(0) {
		opts.BlockTime = defaultBlockTime
	}
	return &Server{
		ServerOpts:  opts,
		blockTime:   opts.BlockTime,
		memPool:     NewTxPool(),
		isValidator: opts.PrivateKey != nil,
		quitCh:      make(chan struct{}),
		rpcCh:       make(chan RPC),
	}
}

// Start 方法启动服务器
// Start method starts the server
func (s *Server) Start() {
	s.initTransport()
	ticker := time.NewTicker(s.blockTime)

free:
	for {
		select {
		case rpc := <-s.rpcCh:
			fmt.Printf("%+v\n", rpc)
		case <-s.quitCh:
			break free
		case <-ticker.C:
			if s.isValidator {
				fmt.Println("creating new block")
			}
		}
	}
	fmt.Println("Server shutdown")
}

func (s *Server) handleTransaction(tx *core.Transaction) error {
	if err := tx.Verify(); err != nil {
		return err
	}

	hash := tx.Hash(core.TxHasher{})
	if s.memPool.Has(hash) {
		logrus.WithFields(logrus.Fields{
			"hash": tx.Hash(core.TxHasher{}),
		}).Info("transaction already in mempool")

		return nil
	}

	logrus.WithFields(logrus.Fields{
		"hash": tx.Hash(core.TxHasher{}),
	}).Info("adding new tx to the mempool")
	return s.memPool.Add(tx)
}

func (s *Server) createNewBlock() error {
	fmt.Println("creating a new block")
	return nil
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
