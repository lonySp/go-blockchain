package core

import (
	"github.com/lonySp/go-blockchain/crypto"
	"github.com/lonySp/go-blockchain/types"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// randomBlock 创建一个随机区块
// randomBlock creates a random block
func randomBlock(height uint32) *Block {
	header := &Header{
		Version:   1,
		PrevBlock: types.RandomHash(),
		Height:    height,
		Timestamp: uint64(time.Now().UnixNano()),
	}
	tx := Transaction{Data: []byte("foo")}

	return NewBlock(header, []Transaction{tx})
}

// randomBlockWithSignature 创建一个带签名的随机区块
// randomBlockWithSignature creates a random block with a signature
func randomBlockWithSignature(t *testing.T, height uint32) *Block {
	privateKey := crypto.GeneratePrivateKey()
	bc := randomBlock(height)
	assert.Nil(t, bc.Sign(privateKey))
	return bc
}

// TestSignBlock 测试区块签名功能
// TestSignBlock tests the block signing functionality
func TestSignBlock(t *testing.T) {
	privateKey := crypto.GeneratePrivateKey()
	b := randomBlock(0)
	assert.Nil(t, b.Sign(privateKey))
	assert.NotNil(t, b.Signature)
}

// TestVerifyBlock 测试区块签名验证功能
// TestVerifyBlock tests the block signature verification functionality
func TestVerifyBlock(t *testing.T) {
	privateKey := crypto.GeneratePrivateKey()
	b := randomBlock(0)

	assert.Nil(t, b.Sign(privateKey))
	assert.Nil(t, b.Verify())

	otherPrivateKey := crypto.GeneratePrivateKey()
	b.Validator = otherPrivateKey.PublicKey()
	assert.NotNil(t, b.Verify())

	b.Height = 100
	assert.NotNil(t, b.Verify())
}
