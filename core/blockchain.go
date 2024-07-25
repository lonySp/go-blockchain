package core

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"sync"
)

// Blockchain 结构体表示区块链
// Blockchain struct represents the blockchain
type Blockchain struct {
	store     Storage      // 存储区块链数据的存储接口 // Storage interface for blockchain data
	lock      sync.RWMutex // 读写锁，用于并发访问区块链 // Read-write lock for concurrent access to the blockchain
	headers   []*Header    // 区块链中的区块头列表 // List of block headers in the blockchain
	validator Validator    // 验证器，用于验证区块 // Validator for validating blocks
}

// NewBlockchain 创建一个新的区块链
// NewBlockchain creates a new blockchain
func NewBlockchain(genesis *Block) (*Blockchain, error) {
	bc := &Blockchain{
		headers: []*Header{},
		store:   NewMemoryStore(), // 使用内存存储 // Use in-memory storage
	}
	bc.validator = NewBlockValidator(bc)
	err := bc.addBlockWithoutValidation(genesis) // 添加创世区块 // Add the genesis block

	return bc, err
}

// SetValidator 设置区块链的验证器
// SetValidator sets the validator for the blockchain
func (bc *Blockchain) SetValidator(v Validator) {
	bc.validator = v
}

// AddBlock 添加一个区块到区块链
// AddBlock adds a block to the blockchain
func (bc *Blockchain) AddBlock(b *Block) error {
	// validate block
	// 验证区块
	if err := bc.validator.ValidateBlock(b); err != nil {
		return err
	}
	return bc.addBlockWithoutValidation(b)
}

// GetHeader 获取指定高度的区块头
// GetHeader gets the block header at the given height
func (bc *Blockchain) GetHeader(height uint32) (*Header, error) {
	if height > bc.Height() {
		return nil, fmt.Errorf("given height (%d) too high", height)
	}

	bc.lock.Lock()
	defer bc.lock.Unlock()

	return bc.headers[height], nil
}

// HasBlock 检查区块链是否包含某个高度的区块
// HasBlock checks if the blockchain contains a block of a certain height
func (bc *Blockchain) HasBlock(height uint32) bool {
	return height <= bc.Height()
}

// Height 返回区块链的高度
// Height returns the height of the blockchain
func (bc *Blockchain) Height() uint32 {
	bc.lock.RLock()
	defer bc.lock.RUnlock()

	return uint32(len(bc.headers) - 1)
}

// addBlockWithoutValidation 添加一个未验证的区块到区块链
// addBlockWithoutValidation adds a block to the blockchain without validation
func (bc *Blockchain) addBlockWithoutValidation(b *Block) error {
	bc.lock.Lock()
	defer bc.lock.Unlock()

	// 添加区块头到区块链
	// Add block header to the blockchain
	bc.headers = append(bc.headers, b.Header)

	// 记录区块添加信息
	// Log block addition information
	logrus.WithFields(logrus.Fields{
		"height": b.Height,
		"hash":   b.Hash(BlockHasher{}),
	}).Info("Adding block to blockchain")

	// 将区块存储到存储接口中
	// Store the block in the storage interface
	return bc.store.Put(b)
}
