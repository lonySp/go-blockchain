package network

import (
	"github.com/lonySp/go-blockchain/core"
	"github.com/lonySp/go-blockchain/types"
	"sort"
)

// TxMapSorter 结构体用于对交易映射进行排序
// TxMapSorter struct is used to sort a map of transactions
type TxMapSorter struct {
	transactions []*core.Transaction // 交易列表 // List of transactions
}

// NewTxMapSorter 创建并返回一个新的 TxMapSorter 实例
// NewTxMapSorter creates and returns a new TxMapSorter instance
func NewTxMapSorter(txMap map[types.Hash]*core.Transaction) *TxMapSorter {
	txx := make([]*core.Transaction, len(txMap))

	i := 0
	for _, val := range txMap {
		txx[i] = val
		i++
	}

	s := &TxMapSorter{txx}

	// 对交易列表进行排序
	// Sort the list of transactions
	sort.Sort(s)

	return s
}

// Len 方法返回交易列表的长度
// Len method returns the length of the transactions list
func (s *TxMapSorter) Len() int {
	return len(s.transactions)
}

// Swap 方法交换交易列表中两个元素的位置
// Swap method swaps the positions of two elements in the transactions list
func (s *TxMapSorter) Swap(i, j int) {
	s.transactions[i], s.transactions[j] = s.transactions[j], s.transactions[i]
}

// Less 方法比较两个交易的先后顺序
// Less method compares the order of two transactions
func (s *TxMapSorter) Less(i, j int) bool {
	return s.transactions[i].FirstSeen() < s.transactions[j].FirstSeen()
}

// TxPool 结构体表示交易池
// TxPool struct represents the transaction pool
type TxPool struct {
	transactions map[types.Hash]*core.Transaction // 交易映射 // Map of transactions
}

// NewTxPool 创建并返回一个新的 TxPool 实例
// NewTxPool creates and returns a new TxPool instance
func NewTxPool() *TxPool {
	return &TxPool{
		transactions: make(map[types.Hash]*core.Transaction),
	}
}

// Transactions 方法返回排序后的交易列表
// Transactions method returns a sorted list of transactions
func (p *TxPool) Transactions() []*core.Transaction {
	s := NewTxMapSorter(p.transactions)
	return s.transactions
}

// Add 方法向交易池添加一个交易，调用者负责验证交易
// Add method adds a transaction to the pool, the caller is responsible for verifying the transaction
func (p *TxPool) Add(tx *core.Transaction) error {
	hash := tx.Hash(core.TxHasher{})
	p.transactions[hash] = tx
	return nil
}

// Has 方法检查交易池中是否包含指定哈希的交易
// Has method checks if the pool contains a transaction with the specified hash
func (p *TxPool) Has(hash types.Hash) bool {
	_, ok := p.transactions[hash]
	return ok
}

// Len 方法返回交易池中的交易数量
// Len method returns the number of transactions in the pool
func (p *TxPool) Len() int {
	return len(p.transactions)
}

// Flush 方法清空交易池
// Flush method clears the transaction pool
func (p *TxPool) Flush() {
	p.transactions = make(map[types.Hash]*core.Transaction)
}
