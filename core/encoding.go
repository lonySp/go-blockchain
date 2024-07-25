package core

import (
	"encoding/gob"
	"io"
)

// Encoder 接口定义了编码器
// Encoder interface defines an encoder
type Encoder[T any] interface {
	Encode(T) error
}

// Decoder 接口定义了解码器
// Decoder interface defines a decoder
type Decoder[T any] interface {
	Decode(T) error
}

// GobTxEncoder 结构体用于基于 gob 的交易编码
// GobTxEncoder struct is used for gob-based transaction encoding
type GobTxEncoder struct {
	w io.Writer // 用于写入编码数据的 io.Writer
}

// NewGobTxEncoder 函数创建一个新的 GobTxEncoder 实例
// NewGobTxEncoder function creates a new instance of GobTxEncoder
func NewGobTxEncoder(w io.Writer) *GobTxEncoder {
	return &GobTxEncoder{w: w}
}

// Encode 方法将交易数据编码为字节流
// Encode method encodes transaction data into a byte stream
func (enc *GobTxEncoder) Encode(tx *Transaction) error {
	return gob.NewEncoder(enc.w).Encode(tx)
}

// GobTxDecoder 结构体用于基于 gob 的交易解码
// GobTxDecoder struct is used for gob-based transaction decoding
type GobTxDecoder struct {
	r io.Reader // 用于读取编码数据的 io.Reader
}

// NewGobTxDecoder 函数创建一个新的 GobTxDecoder 实例
// NewGobTxDecoder function creates a new instance of GobTxDecoder
func NewGobTxDecoder(r io.Reader) *GobTxDecoder {
	return &GobTxDecoder{r: r}
}

// Decode 方法从字节流解码交易数据
// Decode method decodes transaction data from a byte stream
func (dec *GobTxDecoder) Decode(tx *Transaction) error {
	return gob.NewDecoder(dec.r).Decode(tx)
}
