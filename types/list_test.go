package types

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewList 测试创建新列表
// TestNewList tests creating a new list
func TestNewList(t *testing.T) {
	l := NewList[int]()
	assert.Equal(t, l.Data, []int{})
}

// TestListClear 测试清空列表
// TestListClear tests clearing the list
func TestListClear(t *testing.T) {
	l := NewList[int]()
	n := 100

	for i := 0; i < n; i++ {
		l.Insert(i)
	}

	assert.Equal(t, n, l.Len())
	l.Clear()
	assert.Equal(t, 0, l.Len())
}

// TestListContains 测试列表是否包含特定元素
// TestListContains tests if the list contains a specific element
func TestListContains(t *testing.T) {
	l := NewList[int]()
	n := 100

	for i := 0; i < n; i++ {
		l.Insert(i)
		assert.True(t, l.Contains(i))
	}
}

// TestListGetIndex 测试获取元素索引
// TestListGetIndex tests getting the index of an element
func TestListGetIndex(t *testing.T) {
	l := NewList[string]()
	n := 100

	for i := 0; i < n; i++ {
		data := fmt.Sprintf("foo_%d", i)
		l.Insert(data)
		assert.Equal(t, l.GetIndex(data), i)
	}

	assert.Equal(t, l.GetIndex("bar"), -1)
}

// TestListRemove 测试从列表中移除元素
// TestListRemove tests removing an element from the list
func TestListRemove(t *testing.T) {
	l := NewList[string]()
	n := 100

	for i := 0; i < n; i++ {
		data := fmt.Sprintf("foo_%d", i)
		l.Insert(data)
		l.Remove(data)
		assert.Equal(t, l.Contains(data), false)
	}

	assert.Equal(t, l.Len(), 0)
}

// TestListGet 测试获取元素
// TestListGet tests getting an element from the list
func TestListGet(t *testing.T) {
	l := NewList[int]()
	n := 100

	for i := 0; i < n; i++ {
		l.Insert(i)
		assert.True(t, l.Contains(i))
		assert.Equal(t, l.Get(i), i)
	}
}

// TestRemoveAt 测试按索引移除元素
// TestRemoveAt tests removing an element at a specific index
func TestRemoveAt(t *testing.T) {
	l := NewList[int]()
	l.Insert(1)
	l.Insert(2)
	l.Insert(3)
	l.Insert(4)

	l.Pop(0)
	assert.Equal(t, l.Get(0), 2)
}

// TestListAdd 测试添加元素到列表
// TestListAdd tests adding elements to the list
func TestListAdd(t *testing.T) {
	l := NewList[int]()
	n := 100

	for i := 0; i < n; i++ {
		l.Insert(i)
	}

	assert.Equal(t, n, l.Len())
}

// TestListLast 测试获取列表的最后一个元素
// TestListLast tests getting the last element of the list
func TestListLast(t *testing.T) {
	l := NewList[int]()
	l.Insert(1)
	l.Insert(2)
	l.Insert(3)

	assert.Equal(t, 3, l.Last())
}
