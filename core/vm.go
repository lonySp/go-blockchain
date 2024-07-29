package core

import (
	"encoding/binary"
)

// Instruction 定义字节码指令
// Instruction defines bytecode instructions
type Instruction byte

const (
	InstrPushInt  Instruction = 0x0a // 10 // 推入整数 // Push integer
	InstrAdd      Instruction = 0x0b // 11 // 加法操作 // Addition operation
	InstrPushByte Instruction = 0x0c // 12 // 推入字节 // Push byte
	InstrPack     Instruction = 0x0d // 13 // 打包字节数组 // Pack byte array
	InstrSub      Instruction = 0x0e // 14 // 减法操作 // Subtraction operation
	instrStore    Instruction = 0x0f // 15 // 存储 // Store
)

// Stack 结构体表示一个栈
// Stack struct represents a stack
type Stack struct {
	data []any // 栈数据 // Stack data
	sp   int   // 栈指针 // Stack pointer
}

// NewStack 创建一个新的栈实例
// NewStack creates a new stack instance
func NewStack(size int) *Stack {
	return &Stack{
		data: make([]any, size),
		sp:   0,
	}
}

// Push 将元素推入栈
// Push pushes an element onto the stack
func (s *Stack) Push(v any) {
	s.data[s.sp] = v
	s.sp++
}

// Pop 从栈中弹出元素
// Pop pops an element from the stack
func (s *Stack) Pop() any {
	value := s.data[0]
	s.data = append(s.data[:0], s.data[1:]...)
	s.sp--
	return value
}

// VM 结构体表示一个虚拟机
// VM struct represents a virtual machine
type VM struct {
	data          []byte // 字节码数据 // Bytecode data
	ip            int    // 指令指针 // Instruction pointer
	stack         *Stack // 栈实例 // Stack instance
	contractState *State // 合约状态 // Contract state
}

// NewVM 创建一个新的虚拟机实例
// NewVM creates a new VM instance
func NewVM(data []byte, contractState *State) *VM {
	return &VM{
		contractState: contractState,
		data:          data,
		ip:            0,
		stack:         NewStack(128),
	}
}

// Run 运行虚拟机
// Run runs the virtual machine
func (vm *VM) Run() error {
	for {
		instr := Instruction(vm.data[vm.ip]) // 获取当前指令 // Get the current instruction

		if err := vm.Exec(instr); err != nil { // 执行指令 // Execute the instruction
			return err
		}

		vm.ip++                     // 移动指令指针 // Move the instruction pointer
		if vm.ip > len(vm.data)-1 { // 检查是否到达字节码末尾 // Check if the end of the bytecode is reached
			break
		}
	}
	return nil
}

// Exec 执行单个指令
// Exec executes a single instruction
func (vm *VM) Exec(instr Instruction) error {
	switch instr {
	case instrStore:
		// 存储指令 // Store instruction
		var (
			key             = vm.stack.Pop().([]byte)
			value           = vm.stack.Pop()
			serializedValue []byte
		)
		switch v := value.(type) {
		case int:
			serializedValue = serializeInt64(int64(v))
		default:
			panic("TODO: unknown type")
		}
		vm.contractState.Put(key, serializedValue)
	case InstrPushInt:
		vm.stack.Push(int(vm.data[vm.ip-1])) // 推入整数到栈 // Push an integer onto the stack
	case InstrPushByte:
		vm.stack.Push(vm.data[vm.ip-1]) // 推入字节到栈 // Push a byte onto the stack
	case InstrPack:
		n := vm.stack.Pop().(int) // 从栈中弹出整数 // Pop an integer from the stack
		b := make([]byte, n)
		for i := 0; i < n; i++ {
			b[i] = vm.stack.Pop().(byte) // 从栈中弹出字节并放入数组 // Pop bytes from the stack and place in array
		}
		vm.stack.Push(b) // 将字节数组推入栈 // Push the byte array onto the stack
	case InstrSub:
		a := vm.stack.Pop().(int) // 从栈中弹出第一个整数 // Pop the first integer from the stack
		b := vm.stack.Pop().(int) // 从栈中弹出第二个整数 // Pop the second integer from the stack
		c := a - b                // 计算差值 // Calculate the difference
		vm.stack.Push(c)          // 将结果推入栈 // Push the result onto the stack
	case InstrAdd:
		a := vm.stack.Pop().(int) // 从栈中弹出第一个整数 // Pop the first integer from the stack
		b := vm.stack.Pop().(int) // 从栈中弹出第二个整数 // Pop the second integer from the stack
		c := a + b                // 计算和 // Calculate the sum
		vm.stack.Push(c)          // 将结果推入栈 // Push the result onto the stack
	}
	return nil
}

// serializeInt64 序列化 int64
// serializeInt64 serializes int64
func serializeInt64(value int64) []byte {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, uint64(value))
	return buf
}

// deserializeInt64 反序列化 int64
// deserializeInt64 deserializes int64
func deserializeInt64(buf []byte) int64 {
	return int64(binary.LittleEndian.Uint64(buf))
}
