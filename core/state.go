package core

import "fmt"

// State 结构体表示合约的状态存储
// State struct represents the state storage for contracts
type State struct {
	data map[string][]byte // 存储键值对的映射 // Map for storing key-value pairs
}

// NewState 创建一个新的 State 实例
// NewState creates a new instance of State
func NewState() *State {
	return &State{
		data: make(map[string][]byte), // 初始化映射 // Initialize the map
	}
}

// Put 将键值对存储到状态中
// Put stores a key-value pair in the state
func (s *State) Put(k, v []byte) error {
	s.data[string(k)] = v // 将值存储在映射中 // Store the value in the map
	return nil
}

// Delete 从状态中删除指定键的值
// Delete removes the value for the specified key from the state
func (s *State) Delete(k []byte) error {
	delete(s.data, string(k)) // 从映射中删除键值对 // Remove the key-value pair from the map
	return nil
}

// Get 获取指定键的值
// Get retrieves the value for the specified key from the state
func (s *State) Get(k []byte) ([]byte, error) {
	key := string(k)         // 将键转换为字符串 // Convert the key to a string
	value, ok := s.data[key] // 从映射中获取值 // Get the value from the map
	if !ok {                 // 如果键不存在，返回错误 // Return an error if the key does not exist
		return nil, fmt.Errorf("given key %s not found", key)
	}
	return value, nil // 返回值 // Return the value
}
