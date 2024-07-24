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
	bc *Blockchain
}

// NewBlockValidator 创建一个新的区块验证器
// NewBlockValidator creates a new block validator
func NewBlockValidator(bc *Blockchain) *BlockValidator {
	return &BlockValidator{bc: bc}
}

// ValidateBlock 方法验证区块的有效性
// ValidateBlock method validates the validity of the block
func (v *BlockValidator) ValidateBlock(b *Block) error {
	if v.bc.HasBlock(b.Height) {
		return fmt.Errorf("chain already contains block (%d) with hash (%s)", b.Height, b.Hash(BlockHasher{}))
	}

	if b.Height != v.bc.Height()+1 {
		return fmt.Errorf("block (%s) too high", b.Hash(BlockHasher{}))
	}

	prevHeader, err := v.bc.GetHeader(b.Height - 1)
	if err != nil {
		return fmt.Errorf("block (%s) has no previous block", b.Hash(BlockHasher{}))
	}

	hash := BlockHasher{}.Hash(prevHeader)
	if hash != b.PrevBlockHash {
		return fmt.Errorf("block (%s) has invalid previous block hash", b.PrevBlockHash)
	}

	if err := b.Verify(); err != nil {
		return err
	}

	return nil
}
