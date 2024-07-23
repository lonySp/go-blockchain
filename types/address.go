package types

import (
	"encoding/hex"
	"fmt"
)

// Address 结构体表示一个20字节的地址
// Address struct represents a 20-byte address
type Address [20]uint8

// ToSlice 方法将地址转换为字节切片
// ToSlice method converts the address to a byte slice
func (a Address) ToSlice() []byte {
	b := make([]byte, 20)
	for i := 0; i < 20; i++ {
		b[i] = a[i]
	}
	return b
}

// String 方法将地址转换为字符串
// String method converts the address to a string
func (a Address) String() string {
	return hex.EncodeToString(a.ToSlice())
}

// NewAddressFromBytes 函数从字节切片创建一个地址
// NewAddressFromBytes function creates an address from a byte slice
func NewAddressFromBytes(b []byte) Address {
	if len(b) != 20 {
		msg := fmt.Sprintf("given bytes with length %d should be 20", len(b))
		panic(msg)
	}

	var value [20]uint8
	for i := 0; i < 20; i++ {
		value[i] = b[i]
	}
	return value
}
