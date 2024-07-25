package core

import (
	"github.com/lonySp/go-blockchain/crypto"
	"github.com/lonySp/go-blockchain/types"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// TestSignBlock 测试区块签名功能
// TestSignBlock tests the block signing functionality
func TestSignBlock(t *testing.T) {
	privateKey := crypto.GeneratePrivateKey()

	// 创建一个随机区块
	// Create a random block
	b := randomBlock(0, types.Hash{})

	// 测试区块签名
	// Test the block signing
	assert.Nil(t, b.Sign(privateKey))
	assert.NotNil(t, b.Signature)
}

// TestVerifyBlock 测试区块签名验证功能
// TestVerifyBlock tests the block signature verification functionality
func TestVerifyBlock(t *testing.T) {
	privateKey := crypto.GeneratePrivateKey()

	// 创建一个带签名的随机区块
	// Create a random block with a signature
	b := randomBlockWithSignature(t, 0, types.Hash{})

	// 签名区块并验证签名
	// Sign the block and verify the signature
	assert.Nil(t, b.Sign(privateKey))
	assert.Nil(t, b.Verify())

	// 使用另一个私钥验证签名
	// Verify the signature with another private key
	otherPrivateKey := crypto.GeneratePrivateKey()
	b.Validator = otherPrivateKey.PublicKey()
	assert.NotNil(t, b.Verify())

	// 更改区块高度并验证签名
	// Change the block height and verify the signature
	b.Height = 100
	assert.NotNil(t, b.Verify())
}

// randomBlock 创建一个随机区块
// randomBlock creates a random block
func randomBlock(height uint32, prevBlockHash types.Hash) *Block {
	// 创建区块头
	// Create the block header
	header := &Header{
		Version:       1,
		PrevBlockHash: prevBlockHash,
		Height:        height,
		Timestamp:     uint64(time.Now().UnixNano()),
	}

	// 创建区块并返回
	// Create and return the block
	return NewBlock(header, []Transaction{})
}

// randomBlockWithSignature 创建一个带签名的随机区块
// randomBlockWithSignature creates a random block with a signature
func randomBlockWithSignature(t *testing.T, height uint32, prevBlockHash types.Hash) *Block {
	privateKey := crypto.GeneratePrivateKey()

	// 创建一个随机区块
	// Create a random block
	bc := randomBlock(height, prevBlockHash)

	// 创建一个带签名的随机交易并添加到区块中
	// Create a random transaction with a signature and add it to the block
	tx := randomTxWithSignature(t)
	bc.AddTransaction(tx)

	// 签名区块
	// Sign the block
	assert.Nil(t, bc.Sign(privateKey))

	return bc
}
