package core

import (
	"encoding/binary"
)

// Instruction defines the type for VM instructions.
// 指令类型，用于定义虚拟机指令
type Instruction byte

// Constants representing the VM instructions.
// 虚拟机指令常量，每个指令对应一个唯一的 byte 值
const (
	InstrPushInt  Instruction = 0x0a // Push an integer onto the stack. 将整数推入栈中
	InstrAdd      Instruction = 0x0b // Add the top two integers on the stack. 栈顶两个整数相加
	InstrPushByte Instruction = 0x0c // Push a byte onto the stack. 将字节推入栈中
	InstrPack     Instruction = 0x0d // Pack multiple bytes into a byte array. 将多个字节打包成一个字节数组
	InstrSub      Instruction = 0x0e // Subtract the top two integers on the stack. 栈顶两个整数相减
	InstrStore    Instruction = 0x0f // Store the top data on the stack to the contract state. 将栈顶数据存储到合约状态
)

// Stack represents a stack data structure.
// 栈结构，表示一个栈数据结构
type Stack struct {
	data []any // Data stored in the stack. 栈中存储的数据
	sp   int   // Stack pointer indicating the top position of the stack. 栈指针，指示栈顶位置
}

// NewStack creates a new stack with a specified size.
// 创建一个指定大小的新栈
func NewStack(size int) *Stack {
	return &Stack{
		data: make([]any, size),
		sp:   0,
	}
}

// Push an element onto the top of the stack.
// 将元素推入栈顶
func (s *Stack) Push(v any) {
	s.data[s.sp] = v // 将元素放入栈顶位置 // Place the element at the top position
	s.sp++           // 栈指针上移 // Move the stack pointer up
}

// Pop an element from the top of the stack.
// 从栈顶弹出元素
func (s *Stack) Pop() any {
	value := s.data[0]                         // 获取栈顶元素 // Get the top element
	s.data = append(s.data[:0], s.data[1:]...) // 删除栈顶元素并调整栈 // Remove the top element and adjust the stack
	s.sp--                                     // 栈指针下移 // Move the stack pointer down
	return value                               // 返回栈顶元素 // Return the top element
}

// VM represents a virtual machine.
// 虚拟机结构，表示一个虚拟机
type VM struct {
	data          []byte // Contract data. 合约数据
	ip            int    // Instruction pointer. 指令指针
	stack         *Stack // Stack for the VM. 虚拟机的栈
	contractState *State // Contract state. 合约状态
}

// NewVM creates a new virtual machine with the given contract data and state.
// 创建一个包含指定合约数据和状态的新虚拟机
func NewVM(data []byte, contractState *State) *VM {
	return &VM{
		contractState: contractState,
		data:          data,
		ip:            0,
		stack:         NewStack(128),
	}
}

// Run executes the instructions in the virtual machine.
// 运行虚拟机中的指令
func (vm *VM) Run() error {
	for {
		instr := Instruction(vm.data[vm.ip])   // 获取当前指令 // Get the current instruction
		if err := vm.Exec(instr); err != nil { // 执行指令并检查错误 // Execute the instruction and check for errors
			return err // 返回错误 // Return the error
		}
		vm.ip++                     // 移动指令指针到下一条指令 // Move the instruction pointer to the next instruction
		if vm.ip > len(vm.data)-1 { // 检查是否到达指令末尾 // Check if the end of instructions is reached
			break // 退出循环 // Exit the loop
		}
	}
	return nil // 运行完成，返回 nil // Run completed, return nil
}

// Exec executes a single instruction in the virtual machine.
// 执行虚拟机中的单个指令
func (vm *VM) Exec(instr Instruction) error {
	switch instr {
	case InstrStore:
		var (
			key             = vm.stack.Pop().([]byte) // 从栈中弹出键 // Pop the key from the stack
			value           = vm.stack.Pop()          // 从栈中弹出值 // Pop the value from the stack
			serializedValue []byte                    // 序列化后的值 // Serialized value
		)
		switch v := value.(type) {
		case int:
			serializedValue = serializeInt64(int64(v)) // 将整数值序列化 // Serialize the integer value
		default:
			panic("TODO: unknown type") // 未知类型，抛出异常 // Unknown type, panic
		}
		vm.contractState.Put(key, serializedValue) // 将键值对存储到合约状态中 // Store the key-value pair in the contract state

	case InstrPushInt:
		vm.stack.Push(int(vm.data[vm.ip-1])) // 将数据推入栈中 // Push the data onto the stack

	case InstrPushByte:
		vm.stack.Push(byte(vm.data[vm.ip-1])) // 将字节推入栈中 // Push the byte onto the stack

	case InstrPack:
		n := vm.stack.Pop().(int) // 获取要打包的字节数 // Get the number of bytes to pack
		b := make([]byte, n)      // 创建一个字节数组 // Create a byte array
		for i := 0; i < n; i++ {  // 循环从栈中弹出字节并放入字节数组 // Loop to pop bytes from the stack and place them in the byte array
			b[i] = vm.stack.Pop().(byte) // 从栈中弹出字节 // Pop a byte from the stack
		}
		vm.stack.Push(b) // 将字节数组推入栈中 // Push the byte array onto the stack

	case InstrSub:
		a := vm.stack.Pop().(int) // 弹出第一个整数 // Pop the first integer
		b := vm.stack.Pop().(int) // 弹出第二个整数 // Pop the second integer
		c := a - b                // 执行减法运算 // Perform the subtraction
		vm.stack.Push(c)          // 将结果推入栈中 // Push the result onto the stack

	case InstrAdd:
		a := vm.stack.Pop().(int) // 弹出第一个整数 // Pop the first integer
		b := vm.stack.Pop().(int) // 弹出第二个整数 // Pop the second integer
		c := a + b                // 执行加法运算 // Perform the addition
		vm.stack.Push(c)          // 将结果推入栈中 // Push the result onto the stack
	}
	return nil // 执行完成，返回 nil // Execution completed, return nil
}

// serializeInt64 serializes an int64 to a byte slice.
// 序列化 int64 类型为字节切片
func serializeInt64(value int64) []byte {
	buf := make([]byte, 8)                            // 创建一个8字节的缓冲区 // Create an 8-byte buffer
	binary.LittleEndian.PutUint64(buf, uint64(value)) // 使用小端序将 int64 写入缓冲区 // Write the int64 into the buffer using little-endian
	return buf                                        // 返回序列化的字节切片 // Return the serialized byte slice
}

// deserializeInt64 deserializes a byte slice to an int64.
// 反序列化字节切片为 int64 类型
func deserializeInt64(b []byte) int64 {
	return int64(binary.LittleEndian.Uint64(b)) // 使用小端序从字节切片读取 int64 // Read the int64 from the byte slice using little-endian
}
