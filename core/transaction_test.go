package core

import (
	"bytes"
	"github.com/lonySp/go-blockchain/crypto"
	"github.com/lonySp/go-blockchain/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestSignTransaction 测试交易签名功能
// TestSignTransaction tests the transaction signing functionality
func TestSignTransaction(t *testing.T) {
	privateKey := crypto.GeneratePrivateKey()

	data := []byte("foo")
	tx := &Transaction{
		Data: data,
	}

	// 测试交易签名
	// Test the transaction signing
	assert.Nil(t, tx.Sign(privateKey))
	assert.NotNil(t, tx.Signature)
}

// TestVerifyTransaction 测试交易签名验证功能
// TestVerifyTransaction tests the transaction signature verification functionality
func TestVerifyTransaction(t *testing.T) {
	privateKey := crypto.GeneratePrivateKey()
	tx := &Transaction{Data: []byte("foo")}

	// 签名交易并验证签名
	// Sign the transaction and verify the signature
	assert.Nil(t, tx.Sign(privateKey))
	assert.Nil(t, tx.Verify())

	// 测试用其他私钥验证签名的情况
	// Test verifying the signature with another private key
	otherPrivateKey := crypto.GeneratePrivateKey()
	tx.From = otherPrivateKey.PublicKey()
	assert.NotNil(t, tx.Verify())
}

// TestTxEncodeDecode 测试交易的编码和解码
// TestTxEncodeDecode tests the encoding and decoding of a transaction
func TestTxEncodeDecode(t *testing.T) {
	tx := randomTxWithSignature(t)
	buf := &bytes.Buffer{}

	// 编码交易
	// Encode the transaction
	assert.Nil(t, tx.Encode(NewProtobufTxEncoder(buf)))

	// 清空交易哈希
	// Clear the transaction hash
	tx.hash = types.Hash{}

	txDecoded := new(Transaction)

	// 解码交易
	// Decode the transaction
	assert.Nil(t, txDecoded.Decode(NewProtobufTxDecoder(buf)))

	// 验证解码后的交易与原交易相等
	// Verify that the decoded transaction is equal to the original transaction
	assert.Equal(t, &tx, txDecoded)
}

// randomTxWithSignature 创建一个带签名的随机交易
// randomTxWithSignature creates a random transaction with a signature
func randomTxWithSignature(t *testing.T) Transaction {
	privateKey := crypto.GeneratePrivateKey()
	tx := Transaction{
		Data: []byte("foo"),
	}
	assert.Nil(t, tx.Sign(privateKey))
	return tx
}
