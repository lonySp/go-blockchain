package core

import "io"

// Encoder 接口定义了编码器
// Encoder interface defines an encoder
type Encoder[T any] interface {
	Encode(io.Writer, T) error
}

// Decoder 接口定义了解码器
// Decoder interface defines a decoder
type Decoder[T any] interface {
	Decode(io.Reader, T) error
}
