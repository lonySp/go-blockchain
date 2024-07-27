package types

import (
	"encoding/hex"
	"fmt"
	"math/rand"
)

// Hash 结构体表示固定长度的哈希值
// Hash struct represents a fixed-length hash value
type Hash [32]uint8

// IsZero 方法检查哈希值是否为零
// IsZero method checks if the hash value is zero
func (h Hash) IsZero() bool {
	for i := 0; i < 32; i++ {
		if h[i] != 0 {
			return false
		}
	}
	return true
}

// ToSlice 方法将哈希值转换为字节切片
// ToSlice method converts the hash value to a byte slice
func (h Hash) ToSlice() []byte {
	b := make([]byte, 32)
	copy(b, h[:])
	return b
}

// String 方法将哈希值转换为字符串
// String method converts the hash value to a string
func (h Hash) String() string {
	return hex.EncodeToString(h.ToSlice())
}

// HashFromBytes 函数从字节切片创建一个哈希值
// HashFromBytes function creates a hash value from a byte slice
func HashFromBytes(b []byte) Hash {
	if len(b) != 32 {
		msg := fmt.Sprintf("given bytes with length %d should be 32", len(b))
		panic(msg)
	}

	var value [32]uint8
	copy(value[:], b)
	return Hash(value)
}

// RandomBytes 函数生成指定长度的随机字节切片
// RandomBytes function generates a random byte slice of the specified length
func RandomBytes(size int) []byte {
	token := make([]byte, size)
	rand.Read(token)
	return token
}

// RandomHash 函数生成一个随机哈希值
// RandomHash function generates a random hash value
func RandomHash() Hash {
	return HashFromBytes(RandomBytes(32))
}

// BytesToHash 函数将字节数组转换为哈希值
func BytesToHash(b []byte) Hash {
	return HashFromBytes(b)
}
