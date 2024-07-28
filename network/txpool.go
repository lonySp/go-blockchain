package network

import (
	"github.com/lonySp/go-blockchain/core"
	"github.com/lonySp/go-blockchain/types"
	"sync"
)

// TxPool 结构体表示交易池
// TxPool struct represents a transaction pool
type TxPool struct {
	all       *TxSortedMap // 所有交易的有序映射 // Sorted map of all transactions
	pending   *TxSortedMap // 待处理交易的有序映射 // Sorted map of pending transactions
	maxLength int          // 交易池的最大长度 // The maximum length of the transaction pool
}

// NewTxPool 创建并返回一个新的 TxPool 实例
// NewTxPool creates and returns a new instance of TxPool
func NewTxPool(maxLength int) *TxPool {
	return &TxPool{
		all:       NewTxSortedMap(),
		pending:   NewTxSortedMap(),
		maxLength: maxLength,
	}
}

// Add 方法向交易池中添加交易
// Add method adds a transaction to the transaction pool
func (p *TxPool) Add(tx *core.Transaction) {
	// 如果交易池已满，移除最早的交易 // Prune the oldest transaction when the pool is full
	if p.all.Count() == p.maxLength {
		oldest := p.all.First()
		p.all.Remove(oldest.Hash(core.TxHasher{}))
	}

	// 如果交易池中不包含该交易，则添加到交易池中 // Add the transaction to the pool if it does not already exist
	if !p.all.Contains(tx.Hash(core.TxHasher{})) {
		p.all.Add(tx)
		p.pending.Add(tx)
	}
}

// Contains 方法检查交易池中是否包含某个交易哈希
// Contains method checks if the transaction pool contains a given transaction hash
func (p *TxPool) Contains(hash types.Hash) bool {
	return p.all.Contains(hash)
}

// Pending 返回待处理交易的切片
// Pending returns a slice of transactions that are in the pending pool
func (p *TxPool) Pending() []*core.Transaction {
	return p.pending.txx.Data
}

// ClearPending 清除待处理交易
// ClearPending clears all pending transactions
func (p *TxPool) ClearPending() {
	p.pending.Clear()
}

// PendingCount 返回待处理交易的数量
// PendingCount returns the count of pending transactions
func (p *TxPool) PendingCount() int {
	return p.pending.Count()
}

// TxSortedMap 结构体表示有序的交易映射
// TxSortedMap struct represents a sorted map of transactions
type TxSortedMap struct {
	lock   sync.RWMutex                     // 读写锁 // Read-write lock
	lookup map[types.Hash]*core.Transaction // 交易哈希到交易的映射 // Map from transaction hash to transaction
	txx    *types.List[*core.Transaction]   // 交易列表 // List of transactions
}

// NewTxSortedMap 创建并返回一个新的 TxSortedMap 实例
// NewTxSortedMap creates and returns a new instance of TxSortedMap
func NewTxSortedMap() *TxSortedMap {
	return &TxSortedMap{
		lookup: make(map[types.Hash]*core.Transaction),
		txx:    types.NewList[*core.Transaction](),
	}
}

// First 返回交易映射中的第一个交易
// First returns the first transaction in the map
func (t *TxSortedMap) First() *core.Transaction {
	t.lock.RLock()
	defer t.lock.RUnlock()

	first := t.txx.Get(0)
	return t.lookup[first.Hash(core.TxHasher{})]
}

// Get 根据交易哈希获取交易
// Get retrieves a transaction by its hash
func (t *TxSortedMap) Get(h types.Hash) *core.Transaction {
	t.lock.RLock()
	defer t.lock.RUnlock()

	return t.lookup[h]
}

// Add 方法向交易映射中添加交易
// Add method adds a transaction to the map
func (t *TxSortedMap) Add(tx *core.Transaction) {
	hash := tx.Hash(core.TxHasher{})

	t.lock.Lock()
	defer t.lock.Unlock()

	if _, ok := t.lookup[hash]; !ok {
		t.lookup[hash] = tx
		t.txx.Insert(tx)
	}
}

// Remove 方法从交易映射中移除交易
// Remove method removes a transaction from the map
func (t *TxSortedMap) Remove(h types.Hash) {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.txx.Remove(t.lookup[h])
	delete(t.lookup, h)
}

// Count 返回交易映射中的交易数量
// Count returns the number of transactions in the map
func (t *TxSortedMap) Count() int {
	t.lock.RLock()
	defer t.lock.RUnlock()

	return len(t.lookup)
}

// Contains 方法检查交易映射中是否包含某个交易哈希
// Contains method checks if the map contains a given transaction hash
func (t *TxSortedMap) Contains(h types.Hash) bool {
	t.lock.RLock()
	defer t.lock.RUnlock()

	_, ok := t.lookup[h]
	return ok
}

// Clear 方法清空交易映射
// Clear method clears all transactions from the map
func (t *TxSortedMap) Clear() {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.lookup = make(map[types.Hash]*core.Transaction)
	t.txx.Clear()
}
