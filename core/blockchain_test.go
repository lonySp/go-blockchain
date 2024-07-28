package core

import (
	"fmt"
	"github.com/lonySp/go-blockchain/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

// newBlockchainWithGenesis 创建一个带有创世区块的区块链
// newBlockchainWithGenesis creates a blockchain with a genesis block
func newBlockchainWithGenesis(t *testing.T) *Blockchain {
	bc, err := NewBlockchain(randomBlock(t, 0, types.Hash{}))
	assert.Nil(t, err)
	return bc
}

// TestBlockchain 测试区块链初始化功能
// TestBlockchain tests the blockchain initialization functionality
func TestBlockchain(t *testing.T) {
	bc, err := NewBlockchain(randomBlock(t, 0, types.Hash{}))
	assert.Nil(t, err)
	assert.NotNil(t, bc.validator)
	assert.Equal(t, bc.Height(), uint32(0))

	fmt.Println(bc.Height())
}

// TestAddBlock 测试区块链添加区块功能
// TestAddBlock tests the block addition functionality of the blockchain
func TestAddBlock(t *testing.T) {
	bc := newBlockchainWithGenesis(t)
	lenBlocks := 1000
	for i := 0; i < lenBlocks; i++ {
		block := randomBlock(t, uint32(i+1), getPrevBlockHash(t, bc, uint32(i+1)))
		assert.Nil(t, bc.AddBlock(block))
	}

	assert.Equal(t, bc.Height(), uint32(lenBlocks))
	assert.Equal(t, len(bc.headers), lenBlocks+1)
	assert.NotNil(t, bc.AddBlock(randomBlock(t, 89, types.Hash{})))
}

// TestHashBlock 测试区块链是否包含某个高度的区块
// TestHashBlock tests if the blockchain contains a block of a certain height
func TestHashBlock(t *testing.T) {
	bc := newBlockchainWithGenesis(t)
	assert.True(t, bc.HasBlock(0))
	assert.False(t, bc.HasBlock(1))
	assert.False(t, bc.HasBlock(100))
}

// TestGetHeader 测试获取指定高度的区块头
// TestGetHeader tests getting the block header at a specific height
func TestGetHeader(t *testing.T) {
	bc := newBlockchainWithGenesis(t)
	lenBlocks := 1000
	for i := 0; i < lenBlocks; i++ {
		block := randomBlock(t, uint32(i+1), getPrevBlockHash(t, bc, uint32(i+1)))
		assert.Nil(t, bc.AddBlock(block))
		header, err := bc.GetHeader(uint32(i + 1))
		assert.Nil(t, err)
		assert.Equal(t, header, block.Header)
	}
}

// TestAddBlockToHigh 测试添加高度过高的区块
// TestAddBlockToHigh tests adding a block with too high height
func TestAddBlockToHigh(t *testing.T) {
	bc := newBlockchainWithGenesis(t)
	assert.NotNil(t, bc.AddBlock(randomBlock(t, 3, types.Hash{})))
}

// getPrevBlockHash 获取前一个区块的哈希值
// getPrevBlockHash gets the previous block’s hash
func getPrevBlockHash(t *testing.T, bc *Blockchain, height uint32) types.Hash {
	prevHeader, err := bc.GetHeader(height - 1)
	assert.Nil(t, err)
	return BlockHasher{}.Hash(prevHeader)
}
