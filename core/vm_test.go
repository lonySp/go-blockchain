package core

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestStack 测试栈功能
// TestStack tests the stack functionality
func TestStack(t *testing.T) {
	s := NewStack(128) // 创建栈实例 // Create a stack instance
	s.Push(1)          // 推入元素1 // Push element 1
	s.Push(2)          // 推入元素2 // Push element 2

	value := s.Pop()          // 弹出元素 // Pop element
	assert.Equal(t, value, 1) // 验证弹出的元素是否为1 // Verify the popped element is 1

	value = s.Pop()           // 再次弹出元素 // Pop another element
	assert.Equal(t, value, 2) // 验证弹出的元素是否为2 // Verify the popped element is 2
}

// TestVM 测试虚拟机功能
// TestVM tests the VM functionality
func TestVM(t *testing.T) {
	data := []byte{0x03, 0x0a, 0x46, 0x0c, 0x4f, 0x0c, 0x4f, 0x0c, 0x0d} // 示例字节码 // Example bytecode
	vm := NewVM(data)                                                    // 创建虚拟机实例 // Create a VM instance
	assert.Nil(t, vm.Run())                                              // 运行虚拟机并验证是否没有错误 // Run the VM and verify no errors

	result := vm.stack.Pop().([]byte)      // 获取栈顶结果 // Get the result from the stack
	assert.Equal(t, "FOO", string(result)) // 验证结果是否为"FOO" // Verify the result is "FOO"
}
