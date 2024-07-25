package core

import "fmt"

// Validator 接口定义了区块验证方法
// Validator interface defines block validation methods
type Validator interface {
	ValidateBlock(b *Block) error
}

// BlockValidator 结构体实现了 Validator 接口
// BlockValidator struct implements the Validator interface
type BlockValidator struct {
	bc *Blockchain // 区块链实例，用于获取区块链状态 // Blockchain instance to get blockchain state
}

// NewBlockValidator 创建一个新的区块验证器
// NewBlockValidator creates a new block validator
func NewBlockValidator(bc *Blockchain) *BlockValidator {
	return &BlockValidator{bc: bc}
}

// ValidateBlock 方法验证区块的有效性
// ValidateBlock method validates the validity of the block
func (v *BlockValidator) ValidateBlock(b *Block) error {
	// 检查区块链中是否已经包含该高度的区块
	// Check if the blockchain already contains a block at this height
	if v.bc.HasBlock(b.Height) {
		return fmt.Errorf("chain already contains block (%d) with hash (%s)", b.Height, b.Hash(BlockHasher{}))
	}

	// 检查区块的高度是否是当前区块链高度的下一个高度
	// Check if the block height is the next height of the current blockchain height
	if b.Height != v.bc.Height()+1 {
		return fmt.Errorf("block (%s) too high", b.Hash(BlockHasher{}))
	}

	// 获取前一个区块头
	// Get the previous block header
	prevHeader, err := v.bc.GetHeader(b.Height - 1)
	if err != nil {
		return fmt.Errorf("block (%s) has no previous block", b.Hash(BlockHasher{}))
	}

	// 检查前一个区块头的哈希是否匹配
	// Check if the hash of the previous block header matches
	hash := BlockHasher{}.Hash(prevHeader)
	if hash != b.PrevBlockHash {
		return fmt.Errorf("block (%s) has invalid previous block hash", b.PrevBlockHash)
	}

	// 验证区块签名和交易
	// Verify the block signature and transactions
	if err := b.Verify(); err != nil {
		return err
	}

	return nil
}
