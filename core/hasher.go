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

// BlockHasher 实现了 Hasher 接口，用于计算区块头的哈希值
// BlockHasher implements the Hasher interface for calculating the hash of block headers
type BlockHasher struct{}

// Hash 方法计算区块头的哈希值
// Hash method calculates the hash of the block header
func (BlockHasher) Hash(b *Header) types.Hash {
	h := sha256.Sum256(b.Bytes())
	return h
}

// TxHasher 实现了 Hasher 接口，用于计算交易的哈希值
// TxHasher implements the Hasher interface for calculating the hash of transactions
type TxHasher struct{}

// Hash 方法计算交易的哈希值
// Hash method calculates the hash of the transaction
func (TxHasher) Hash(tx *Transaction) types.Hash {
	return sha256.Sum256(tx.Data)
}
