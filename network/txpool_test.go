package network

import (
	"github.com/lonySp/go-blockchain/core"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"strconv"
	"testing"
)

// TestTxPool 测试交易池的初始化
// TestTxPool tests the initialization of the transaction pool
func TestTxPool(t *testing.T) {
	p := NewTxPool()
	// 确认交易池长度为 0
	// Assert that the transaction pool length is 0
	assert.Equal(t, 0, p.Len())
}

// TestTxPoolAddTx 测试向交易池添加交易
// TestTxPoolAddTx tests adding a transaction to the transaction pool
func TestTxPoolAddTx(t *testing.T) {
	p := NewTxPool()
	tx := core.NewTransaction([]byte("foo"))
	// 向交易池添加交易并确认没有错误
	// Add a transaction to the pool and assert no error
	assert.Nil(t, p.Add(tx))
	// 确认交易池长度为 1
	// Assert that the transaction pool length is 1
	assert.Equal(t, 1, p.Len())

	// 添加相同的交易并确认交易池长度仍为 1
	// Add the same transaction and assert that the pool length is still 1
	_ = core.NewTransaction([]byte("foo"))
	assert.Equal(t, p.Len(), 1)

	// 清空交易池并确认长度为 0
	// Flush the transaction pool and assert the length is 0
	p.Flush()
	assert.Equal(t, 0, p.Len())
}

// TestSortTransactions 测试交易的排序功能
// TestSortTransactions tests the sorting of transactions
func TestSortTransactions(t *testing.T) {
	p := NewTxPool()
	txLen := 1000
	for i := 0; i < txLen; i++ {
		tx := core.NewTransaction([]byte(strconv.FormatInt(int64(i), 10)))
		// 设置交易的首次看到时间
		// Set the first seen time of the transaction
		tx.SetFirstSeen(int64(i * rand.Intn(10000)))
		// 向交易池添加交易并确认没有错误
		// Add the transaction to the pool and assert no error
		assert.Nil(t, p.Add(tx))
	}

	// 确认交易池长度为 txLen
	// Assert that the transaction pool length is txLen
	assert.Equal(t, txLen, p.Len())

	// 获取排序后的交易列表
	// Get the sorted list of transactions
	txx := p.Transactions()
	for i := 0; i < len(txx)-1; i++ {
		// 确认每个交易的首次看到时间按升序排序
		// Assert that each transaction's first seen time is in ascending order
		assert.True(t, txx[i].FirstSeen() < txx[i+1].FirstSeen())
	}
}
