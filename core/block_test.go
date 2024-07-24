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
	b := randomBlock(0, types.Hash{})
	assert.Nil(t, b.Sign(privateKey))
	assert.NotNil(t, b.Signature)
}

// TestVerifyBlock 测试区块签名验证功能
// TestVerifyBlock tests the block signature verification functionality
func TestVerifyBlock(t *testing.T) {
	privateKey := crypto.GeneratePrivateKey()
	b := randomBlockWithSignature(t, 0, types.Hash{})

	assert.Nil(t, b.Sign(privateKey))
	assert.Nil(t, b.Verify())

	otherPrivateKey := crypto.GeneratePrivateKey()
	b.Validator = otherPrivateKey.PublicKey()
	assert.NotNil(t, b.Verify())

	b.Height = 100
	assert.NotNil(t, b.Verify())
}

// randomBlock 创建一个随机区块
// randomBlock creates a random block
func randomBlock(height uint32, prevBlockHash types.Hash) *Block {
	header := &Header{
		Version:       1,
		PrevBlockHash: prevBlockHash,
		Height:        height,
		Timestamp:     uint64(time.Now().UnixNano()),
	}

	return NewBlock(header, []Transaction{})
}

// randomBlockWithSignature 创建一个带签名的随机区块
// randomBlockWithSignature creates a random block with a signature
func randomBlockWithSignature(t *testing.T, height uint32, prevBlockHash types.Hash) *Block {
	privateKey := crypto.GeneratePrivateKey()
	bc := randomBlock(height, prevBlockHash)
	tx := randomTxWithSignature(t)
	bc.AddTransaction(tx)
	assert.Nil(t, bc.Sign(privateKey))
	return bc
}
