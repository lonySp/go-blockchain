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
	Transport  []Transport        // 传输选项 // Transport options
	BlockTime  time.Duration      // 区块生成时间间隔 // Block creation time interval
	PrivateKey *crypto.PrivateKey // 私钥，用于签名 // Private key for signing
}

// Server 结构体表示服务器
// Server struct represents the server
type Server struct {
	ServerOpts                // 嵌入 ServerOpts 结构体 // Embedding ServerOpts struct
	blockTime   time.Duration // 区块生成时间间隔 // Block creation time interval
	memPool     *TxPool       // 交易池 // Transaction pool
	isValidator bool          // 是否是验证者 // Whether the server is a validator
	rpcCh       chan RPC      // 接收 RPC 消息的通道 // Channel for receiving RPC messages
	quitCh      chan struct{} // 关闭服务器的通道 // Channel for shutting down the server
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
			// 处理接收到的 RPC 消息
			// Handle received RPC messages
			fmt.Printf("%+v\n", rpc)
		case <-s.quitCh:
			// 处理服务器关闭信号
			// Handle server shutdown signal
			break free
		case <-ticker.C:
			// 如果是验证者，则创建新块
			// If the server is a validator, create a new block
			if s.isValidator {
				fmt.Println("creating new block")
			}
		}
	}
	fmt.Println("Server shutdown")
}

// handleTransaction 方法处理接收到的交易
// handleTransaction method handles received transactions
func (s *Server) handleTransaction(tx *core.Transaction) error {
	// 验证交易
	// Verify the transaction
	if err := tx.Verify(); err != nil {
		return err
	}

	hash := tx.Hash(core.TxHasher{})
	// 检查交易池中是否已有该交易
	// Check if the transaction is already in the transaction pool
	if s.memPool.Has(hash) {
		logrus.WithFields(logrus.Fields{
			"hash": tx.Hash(core.TxHasher{}),
		}).Info("transaction already in mempool")
		return nil
	}

	// 添加新交易到交易池
	// Add new transaction to the transaction pool
	logrus.WithFields(logrus.Fields{
		"hash": tx.Hash(core.TxHasher{}),
	}).Info("adding new tx to the mempool")
	return s.memPool.Add(tx)
}

// createNewBlock 方法创建一个新块
// createNewBlock method creates a new block
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
				// Handle RPC messages
				s.rpcCh <- rpc
			}
		}(tr)
	}
}
