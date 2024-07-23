package core

// Storage 接口定义了存储方法
// Storage interface defines storage methods
type Storage interface {
	Put(block *Block) error
}

// MemoryStore 结构体表示内存存储
// MemoryStore struct represents an in-memory storage
type MemoryStore struct {
}

// NewMemoryStore 创建一个新的内存存储
// NewMemoryStore creates a new in-memory storage
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{}
}

// Put 方法将区块存储在内存中
// Put method stores the block in memory
func (s *MemoryStore) Put(b *Block) error {
	return nil
}
