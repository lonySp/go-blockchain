package core

import (
	"crypto/sha256"
	"github.com/lonySp/go-blockchain/types"
)

// Hasher 接口定义了哈希器
// Hasher interface defines a hasher
type Hasher[T any] interface {
	Hash(T) types.Hash
}

// BlockHasher 实现了 Hasher 接口
// BlockHasher implements the Hasher interface
type BlockHasher struct{}

// Hash 方法计算区块的哈希值
// Hash method calculates the hash of the block
func (BlockHasher) Hash(b *Block) types.Hash {
	h := sha256.Sum256(b.HeaderData())
	return h
}
