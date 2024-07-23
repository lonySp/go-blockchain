package core

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

// newBlockchainWithGenesis 创建一个带有创世区块的区块链
// newBlockchainWithGenesis creates a blockchain with a genesis block
func newBlockchainWithGenesis(t *testing.T) *Blockchain {
	bc, err := NewBlockchain(randomBlock(0))
	assert.Nil(t, err)
	return bc
}

// TestBlockchain 测试区块链初始化功能
// TestBlockchain tests the blockchain initialization functionality
func TestBlockchain(t *testing.T) {
	bc, err := NewBlockchain(randomBlock(0))
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
		block := randomBlockWithSignature(t, uint32(i+1))
		assert.Nil(t, bc.AddBlock(block))
	}

	assert.Equal(t, bc.Height(), uint32(lenBlocks))
	assert.Equal(t, len(bc.headers), lenBlocks+1)
	assert.NotNil(t, bc.AddBlock(randomBlock(89)))
}

// TestHashBlock 测试区块链是否包含某个高度的区块
// TestHashBlock tests if the blockchain contains a block of a certain height
func TestHashBlock(t *testing.T) {
	bc := newBlockchainWithGenesis(t)
	assert.True(t, bc.HasBlock(0))
}
