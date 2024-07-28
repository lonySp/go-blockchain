package core

import (
	"github.com/lonySp/go-blockchain/types"
	"testing"

	"github.com/go-kit/log"
	"github.com/stretchr/testify/assert"
)

// TestAddBlock 测试添加区块到区块链
// TestAddBlock tests adding blocks to the blockchain
func TestAddBlock(t *testing.T) {
	// 创建带有创世区块的区块链 // Create a blockchain with a genesis block
	bc := newBlockchainWithGenesis(t)

	lenBlocks := 1000
	for i := 0; i < lenBlocks; i++ {
		block := randomBlock(t, uint32(i+1), getPrevBlockHash(t, bc, uint32(i+1)))
		assert.Nil(t, bc.AddBlock(block))
	}

	// 验证区块链高度和区块头数量是否正确 // Verify the blockchain height and number of headers
	assert.Equal(t, bc.Height(), uint32(lenBlocks))
	assert.Equal(t, len(bc.headers), lenBlocks+1)
	assert.NotNil(t, bc.AddBlock(randomBlock(t, 89, types.Hash{})))
}

// TestNewBlockchain 测试创建新的区块链
// TestNewBlockchain tests creating a new blockchain
func TestNewBlockchain(t *testing.T) {
	bc := newBlockchainWithGenesis(t)
	assert.NotNil(t, bc.validator)
	assert.Equal(t, bc.Height(), uint32(0))
}

// TestHasBlock 测试区块链是否包含指定高度的区块
// TestHasBlock tests if the blockchain contains a block at a given height
func TestHasBlock(t *testing.T) {
	bc := newBlockchainWithGenesis(t)
	assert.True(t, bc.HasBlock(0))
	assert.False(t, bc.HasBlock(1))
	assert.False(t, bc.HasBlock(100))
}

// TestGetHeader 测试获取区块头
// TestGetHeader tests getting a block header
func TestGetHeader(t *testing.T) {
	bc := newBlockchainWithGenesis(t)
	lenBlocks := 1000

	for i := 0; i < lenBlocks; i++ {
		block := randomBlock(t, uint32(i+1), getPrevBlockHash(t, bc, uint32(i+1)))
		assert.Nil(t, bc.AddBlock(block))
		header, err := bc.GetHeader(block.Height)
		assert.Nil(t, err)
		assert.Equal(t, header, block.Header)
	}
}

// TestAddBlockToHigh 测试添加高度过高的区块
// TestAddBlockToHigh tests adding a block with too high height
func TestAddBlockToHigh(t *testing.T) {
	bc := newBlockchainWithGenesis(t)

	assert.Nil(t, bc.AddBlock(randomBlock(t, 1, getPrevBlockHash(t, bc, uint32(1)))))
	assert.NotNil(t, bc.AddBlock(randomBlock(t, 3, types.Hash{})))
}

// newBlockchainWithGenesis 创建带有创世区块的区块链
// newBlockchainWithGenesis creates a blockchain with a genesis block
func newBlockchainWithGenesis(t *testing.T) *Blockchain {
	bc, err := NewBlockchain(log.NewNopLogger(), randomBlock(t, 0, types.Hash{}))
	assert.Nil(t, err)

	return bc
}

// getPrevBlockHash 获取前一个区块的哈希值
// getPrevBlockHash gets the hash of the previous block
func getPrevBlockHash(t *testing.T, bc *Blockchain, height uint32) types.Hash {
	prevHeader, err := bc.GetHeader(height - 1)
	assert.Nil(t, err)
	return BlockHasher{}.Hash(prevHeader)
}
