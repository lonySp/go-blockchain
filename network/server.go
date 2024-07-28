package network

import (
	"bytes"
	"github.com/go-kit/log"
	"github.com/lonySp/go-blockchain/core"
	"github.com/lonySp/go-blockchain/crypto"
	"github.com/lonySp/go-blockchain/types"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

// defaultBlockTime 定义区块生成的默认时间间隔
// defaultBlockTime defines the default interval for block creation
var defaultBlockTime = 5 * time.Second

// ServerOpts 结构体包含服务器的传输选项
// ServerOpts struct contains server transport options
type ServerOpts struct {
	ID            string             // 服务器的唯一标识符 // Unique identifier for the server
	Logger        log.Logger         // 日志记录器 // Logger
	RPCDecodeFunc RPCDecodeFunc      // RPC 解码函数 // RPC decode function
	RPCProcessor  RPCProcessor       // RPC 处理器 // RPC processor
	Transport     []Transport        // 传输选项 // Transport options
	BlockTime     time.Duration      // 区块生成时间间隔 // Block creation time interval
	PrivateKey    *crypto.PrivateKey // 私钥，用于签名 // Private key for signing
}

// Server 结构体表示服务器
// Server struct represents the server
type Server struct {
	ServerOpts                   // 嵌入 ServerOpts 结构体 // Embedding ServerOpts struct
	chain       *core.Blockchain // 区块链实例 // Blockchain instance
	memPool     *TxPool          // 交易池 // Transaction pool
	isValidator bool             // 是否是验证者 // Whether the server is a validator
	rpcCh       chan RPC         // 接收 RPC 消息的通道 // Channel for receiving RPC messages
	quitCh      chan struct{}    // 关闭服务器的通道 // Channel for shutting down the server
}

// NewServer 创建并返回一个新的 Server 实例
// NewServer creates and returns a new Server instance
func NewServer(opts ServerOpts) (*Server, error) {
	// 设置区块生成时间间隔的默认值 // Set default block creation time interval
	if opts.BlockTime == time.Duration(0) {
		opts.BlockTime = defaultBlockTime
	}
	// 设置默认的 RPC 解码函数 // Set default RPC decode function
	if opts.RPCDecodeFunc == nil {
		opts.RPCDecodeFunc = DefaultRPCDecodeFunc
	}
	// 设置默认的日志记录器 // Set default logger
	if opts.Logger == nil {
		opts.Logger = log.NewLogfmtLogger(os.Stderr)
		opts.Logger = log.With(opts.Logger, "ID", opts.ID)
	}

	// 创建区块链实例 // Create blockchain instance
	chain, err := core.NewBlockchain(opts.Logger, genesisBlock())
	if err != nil {
		return nil, err
	}
	s := &Server{
		ServerOpts:  opts,
		chain:       chain,
		memPool:     NewTxPool(1000),
		isValidator: opts.PrivateKey != nil,
		quitCh:      make(chan struct{}, 1),
		rpcCh:       make(chan RPC),
	}
	// 如果未提供 RPC 处理器，使用服务器本身作为默认处理器
	// If no RPC processor is provided, use the server itself as the default processor
	if s.RPCProcessor == nil {
		s.RPCProcessor = s
	}

	// 如果服务器是验证者，启动验证者循环
	// If the server is a validator, start the validator loop
	if s.isValidator {
		go s.validatorLoop()
	}
	return s, nil
}

// Start 方法启动服务器
// Start method starts the server
func (s *Server) Start() {
	// 初始化传输选项 // Initialize transport options
	s.initTransport()
free:
	for {
		select {
		case rpc := <-s.rpcCh:
			// 解码 RPC 消息 // Decode the RPC message
			msg, err := s.RPCDecodeFunc(rpc)
			if err != nil {
				logrus.Error("Error", err)
			}

			// 处理解码后的消息 // Process the decoded message
			if err := s.RPCProcessor.ProcessMessage(msg); err != nil {
				logrus.Error("Error", err)
			}
		case <-s.quitCh:
			// 处理服务器关闭信号 // Handle server shutdown signal
			break free
		}
	}
	// 记录服务器关闭日志 // Log server shutdown
	s.Logger.Log("msg", "Server is shutting down")
}

// validatorLoop 是验证者的主循环，负责定期创建新区块
// validatorLoop is the main loop for the validator, responsible for creating new blocks periodically
func (s *Server) validatorLoop() {
	ticker := time.NewTicker(s.BlockTime)
	s.Logger.Log("msg", "Starting validator loop", "blockTime", s.BlockTime)
	for {
		<-ticker.C
		s.createNewBlock()
	}
}

// ProcessMessage 方法处理解码后的消息
// ProcessMessage method processes the decoded message
func (s *Server) ProcessMessage(msg *DecodedMessage) error {
	switch t := msg.Data.(type) {
	case *core.Transaction:
		return s.processTransaction(t)
	case *core.Block:
		return s.processBlock(t)
	}
	return nil
}

// broadcast 方法将消息广播到所有已连接的传输节点
// broadcast method broadcasts a message to all connected transport nodes
func (s *Server) broadcast(payload []byte) error {
	for _, tr := range s.Transport {
		if err := tr.Broadcast(payload); err != nil {
			return err
		}
	}
	return nil
}

// processBlock 方法处理区块
// processBlock method processes a block
func (s *Server) processBlock(b *core.Block) error {
	if err := s.chain.AddBlock(b); err != nil {
		return err
	}
	go s.broadcastBlock(b)
	return nil
}

// processTransaction 方法处理交易
// processTransaction method processes a transaction
func (s *Server) processTransaction(tx *core.Transaction) error {
	// 计算交易的哈希值 // Calculate the transaction hash
	hash := tx.Hash(core.TxHasher{})

	// 如果交易池中已经有该交易，则返回 // Return if the transaction already exists in the transaction pool
	if s.memPool.Contains(hash) {
		return nil
	}

	// 验证交易 // Verify the transaction
	if err := tx.Verify(); err != nil {
		return err
	}

	// 记录日志 // Log the transaction addition to the mempool
	// s.Logger.Log(
	//	"msg", "adding new tx to mempool",
	//	"hash", hash,
	//	"mempoolLength", s.memPool.PendingCount())

	// 广播交易 // Broadcast the transaction
	go s.broadcastTx(tx)

	s.memPool.Add(tx)

	return nil
}

// broadcastBlock 方法广播区块
// broadcastBlock method broadcasts a block
func (s *Server) broadcastBlock(b *core.Block) error {
	buf := &bytes.Buffer{}
	if err := b.Encode(core.NewProtobufBlockEncoder(buf)); err != nil {
		return err
	}
	msg := NewMessage(MessageTypeBlock, buf.Bytes())
	return s.broadcast(msg.Bytes())
}

// broadcastTx 方法广播交易
// broadcastTx method broadcasts a transaction
func (s *Server) broadcastTx(tx *core.Transaction) error {
	buf := &bytes.Buffer{}
	if err := tx.Encode(core.NewProtobufTxEncoder(buf)); err != nil {
		return err
	}

	// 创建消息并广播 // Create a message and broadcast it
	msg := NewMessage(MessageTypeTx, buf.Bytes())
	return s.broadcast(msg.Bytes())
}

// createNewBlock 方法创建一个新块
// createNewBlock method creates a new block
func (s *Server) createNewBlock() error {
	// 获取当前区块头
	// Get the current block header
	currentHeader, err := s.chain.GetHeader(s.chain.Height())
	if err != nil {
		return err
	}

	// 目前我们将使用所有在交易池中的交易
	// For now, we are going to use all transactions that are in the mempool
	// 后续我们了解交易的内部结构后
	// Later on, when we know the internal structure of our transaction
	// 我们将实现某种复杂度函数来确定一个区块中可以包含多少交易
	// We will implement some kind of complexity function to determine how many transactions can be included in a block.
	txx := s.memPool.Pending()

	// 创建新的区块
	// Create a new block
	block, err := core.NewBlockFromPrevHeader(currentHeader, txx)
	if err != nil {
		return err
	}

	// 使用私钥签名区块
	// Sign the block with the private key
	if err := block.Sign(*s.PrivateKey); err != nil {
		return err
	}

	// 将区块添加到区块链
	// Add the block to the blockchain
	if err := s.chain.AddBlock(block); err != nil {
		return err
	}

	// 清除已包含在区块中的待处理交易
	// Clear pending transactions that have been included in the block
	s.memPool.ClearPending()

	// 异步广播新创建的区块
	// Asynchronously broadcast the newly created block
	go s.broadcastBlock(block)

	return nil
}

// initTransport 方法初始化所有的传输选项
// initTransport method initializes all transport options
func (s *Server) initTransport() {
	for _, tr := range s.Transport {
		go func(tr Transport) {
			for rpc := range tr.Consume() {
				// 处理 RPC 消息 // Handle RPC messages
				s.rpcCh <- rpc
			}
		}(tr)
	}
}

// genesisBlock 创建并返回创世区块
// genesisBlock creates and returns the genesis block
func genesisBlock() *core.Block {
	header := &core.Header{
		Version:   1,            // 区块版本号 // Block version number
		Height:    0,            // 区块高度 // Block height
		Timestamp: 000000,       // 区块生成的时间戳 // Timestamp when the block was created
		DataHash:  types.Hash{}, // 区块数据哈希值 // Hash of the block data
	}

	// 创建并返回创世区块 // Create and return the genesis block
	b, _ := core.NewBlock(header, nil)
	return b
}
