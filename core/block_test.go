package core

import (
	"bytes"
	"fmt"
	"github.com/lonySp/go-blockchain/types"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// TestHeader_Encode_Decode 测试区块头的编码和解码
// TestHeader_Encode_Decode tests the encoding and decoding of a block header
func TestHeader_Encode_Decode(t *testing.T) {
	// 创建一个区块头
	// Create a block header
	h := &Header{
		Height:    10,
		Nonce:     989394,
		PrevBlock: types.RandomHash(),
		Timestamp: uint64(time.Now().UnixNano()),
		Version:   1,
	}
	// 将区块头进行编码
	// Encode the block header
	buf := &bytes.Buffer{}
	assert.Nil(t, h.EnCodeBinary(buf))

	// 进行解码
	// Decode the block header
	hDecode := &Header{}
	assert.Nil(t, hDecode.DecodeBinary(buf))

	// 判断解码前后是否一致
	// Check if the decoded header is the same as the original
	assert.Equal(t, h, hDecode)
}

// TestBlock_EncodeBinary 测试区块的编码
// TestBlock_EncodeBinary tests the encoding of a block
func TestBlock_EncodeBinary(t *testing.T) {
	// 将区块进行编码
	// Create and encode a block
	b := &Block{
		Header: Header{
			Height:    10,
			Nonce:     989394,
			PrevBlock: types.RandomHash(),
			Timestamp: uint64(time.Now().UnixNano()),
			Version:   1,
		},
		Transaction: nil,
	}
	buf := &bytes.Buffer{}
	assert.Nil(t, b.EncodeBinary(buf))

	// 在进行解码测试
	// Decode the block and test
	bDecode := &Block{}
	assert.Nil(t, bDecode.DecodeBinary(buf))

	// 判断解码前后是否一致
	// Check if the decoded block is the same as the original
	assert.Equal(t, b, bDecode)
}

// TestBlock_Hash 测试区块哈希
// TestBlock_Hash tests the hashing of a block
func TestBlock_Hash(t *testing.T) {
	b := &Block{
		Header: Header{
			Height:    10,
			Nonce:     989394,
			PrevBlock: types.RandomHash(),
			Timestamp: uint64(time.Now().UnixNano()),
			Version:   1,
		},
		Transaction: nil,
	}

	h := b.Hash()
	fmt.Println(h)
	// 检查哈希值是否非零
	// Check if the hash value is not zero
	assert.False(t, h.IsZero())
}
