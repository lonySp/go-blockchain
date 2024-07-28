package types

import (
	"fmt"
	"reflect"
)

// List 结构体表示一个泛型列表
// List struct represents a generic list
type List[T any] struct {
	Data []T // 列表数据 // List data
}

// NewList 创建并返回一个新的 List 实例
// NewList creates and returns a new instance of List
func NewList[T any]() *List[T] {
	return &List[T]{
		Data: []T{},
	}
}

// Get 返回列表中指定索引处的元素，如果索引越界则抛出错误
// Get returns the element at the specified index in the list, and panics if the index is out of range
func (l *List[T]) Get(index int) T {
	if index > len(l.Data)-1 {
		err := fmt.Sprintf("the given index (%d) is higher than the length (%d)", index, len(l.Data))
		panic(err)
	}
	return l.Data[index]
}

// Insert 将元素插入列表
// Insert inserts an element into the list
func (l *List[T]) Insert(v T) {
	l.Data = append(l.Data, v)
}

// Clear 清空列表
// Clear clears the list
func (l *List[T]) Clear() {
	l.Data = []T{}
}

// GetIndex 返回元素在列表中的索引，如果元素不存在则返回 -1
// GetIndex returns the index of an element in the list, or -1 if the element does not exist
func (l *List[T]) GetIndex(v T) int {
	for i := 0; i < l.Len(); i++ {
		if reflect.DeepEqual(v, l.Data[i]) {
			return i
		}
	}
	return -1
}

// Remove 从列表中移除指定元素
// Remove removes the specified element from the list
func (l *List[T]) Remove(v T) {
	index := l.GetIndex(v)
	if index == -1 {
		return
	}
	l.Pop(index)
}

// Pop 移除指定索引处的元素
// Pop removes the element at the specified index
func (l *List[T]) Pop(index int) {
	l.Data = append(l.Data[:index], l.Data[index+1:]...)
}

// Contains 检查列表是否包含指定元素
// Contains checks if the list contains the specified element
func (l *List[T]) Contains(v T) bool {
	for i := 0; i < len(l.Data); i++ {
		if reflect.DeepEqual(l.Data[i], v) {
			return true
		}
	}
	return false
}

// Last 返回列表中的最后一个元素
// Last returns the last element in the list
func (l List[T]) Last() T {
	return l.Data[l.Len()-1]
}

// Len 返回列表的长度
// Len returns the length of the list
func (l *List[T]) Len() int {
	return len(l.Data)
}
