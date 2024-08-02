package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestStack 用于测试 Stack 数据结构的功能
// TestStack is used to test the functionality of the Stack data structure
func TestStack(t *testing.T) {
	// 创建一个大小为128的栈
	// Create a stack with a size of 128
	s := NewStack(128)

	// 压入元素1和2
	// Push elements 1 and 2
	s.Push(1)
	s.Push(2)

	// 弹出一个元素并检查是否为1
	// Pop an element and check if it is 1
	value := s.Pop()
	assert.Equal(t, value, 1)

	// 再弹出一个元素并检查是否为2
	// Pop another element and check if it is 2
	value = s.Pop()
	assert.Equal(t, value, 2)
}

// TestVM 用于测试虚拟机的功能
// TestVM is used to test the functionality of the virtual machine (VM)
func TestVM(t *testing.T) {
	// 初始化合约数据
	// Initialize contract data
	data := []byte{0x03, 0x0a, 0x46, 0x0c, 0x4f, 0x0c, 0x4f, 0x0c, 0x0d, 0x05, 0x0a, 0x0f}
	// 创建合约状态
	// Create contract state
	contractState := NewState()
	// 创建虚拟机实例
	// Create VM instance
	vm := NewVM(data, contractState)
	// 运行虚拟机并确保没有错误
	// Run the VM and ensure no errors
	assert.Nil(t, vm.Run())

	// 从合约状态获取键为 "FOO" 的值并反序列化为 int64 类型
	// Get the value with key "FOO" from the contract state and deserialize it to int64
	valueBytes, err := contractState.Get([]byte("FOO"))
	value := deserializeInt64(valueBytes)
	// 确保没有错误
	// Ensure no errors
	assert.Nil(t, err)
	// 检查值是否为5
	// Check if the value is 5
	assert.Equal(t, value, int64(5))
}
