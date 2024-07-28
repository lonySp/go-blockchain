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
	// 生成私钥 // Generate a private key
	privateKey := crypto.GeneratePrivateKey()

	// 创建一个随机区块 // Create a random block
	b := randomBlock(t, 0, types.Hash{})

	// 测试区块签名 // Test the block signing
	assert.Nil(t, b.Sign(privateKey)) // 验证签名操作是否成功 // Verify if the signing operation is successful
	assert.NotNil(t, b.Signature)     // 验证签名是否存在 // Verify if the signature is present
}

// TestVerifyBlock 测试区块签名验证功能
// TestVerifyBlock tests the block signature verification functionality
func TestVerifyBlock(t *testing.T) {
	// 生成私钥 // Generate a private key
	privateKey := crypto.GeneratePrivateKey()

	// 创建一个带签名的随机区块 // Create a random block with a signature
	b := randomBlock(t, 0, types.Hash{})

	// 签名区块并验证签名 // Sign the block and verify the signature
	assert.Nil(t, b.Sign(privateKey)) // 验证签名操作是否成功 // Verify if the signing operation is successful
	assert.Nil(t, b.Verify())         // 验证签名是否有效 // Verify if the signature is valid

	// 使用另一个私钥验证签名 // Verify the signature with another private key
	otherPrivateKey := crypto.GeneratePrivateKey()
	b.Validator = otherPrivateKey.PublicKey() // 更换验证者的公钥 // Change the validator's public key
	assert.NotNil(t, b.Verify())              // 验证签名是否无效 // Verify if the signature is invalid

	// 更改区块高度并验证签名 // Change the block height and verify the signature
	b.Height = 100               // 更改区块高度 // Change the block height
	assert.NotNil(t, b.Verify()) // 验证签名是否无效 // Verify if the signature is invalid
}

// randomBlock 创建一个随机区块
// randomBlock creates a random block
func randomBlock(t *testing.T, height uint32, prevBlockHash types.Hash) *Block {
	// 生成私钥 // Generate a private key
	privateKey := crypto.GeneratePrivateKey()

	// 创建一个随机交易并签名 // Create a random transaction with a signature
	tx := randomTxWithSignature(t)

	// 创建区块头 // Create the block header
	header := &Header{
		Version:       1,
		PrevBlockHash: prevBlockHash,
		Height:        height,
		Timestamp:     uint64(time.Now().UnixNano()),
	}

	// 创建区块并返回 // Create and return the block
	b, err := NewBlock(header, []*Transaction{tx})
	assert.Nil(t, err) // 验证区块创建是否成功 // Verify if the block creation is successful

	// 计算并设置区块数据哈希 // Calculate and set the block data hash
	dataHash, err := CalculateDataHash(b.Transactions)
	assert.Nil(t, err) // 验证数据哈希计算是否成功 // Verify if the data hash calculation is successful
	b.Header.DataHash = dataHash

	// 签名区块 // Sign the block
	assert.Nil(t, b.Sign(privateKey)) // 验证签名操作是否成功 // Verify if the signing operation is successful

	return b
}
