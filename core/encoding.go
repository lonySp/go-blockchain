package core

import (
	"github.com/golang/protobuf/proto"
	"github.com/lonySp/go-blockchain/crypto"
	"github.com/lonySp/go-blockchain/types"
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

// ProtobufTxEncoder 结构体用于基于 Protobuf 的交易编码
// ProtobufTxEncoder struct is used for Protobuf-based transaction encoding
type ProtobufTxEncoder struct {
	w io.Writer // 用于写入编码数据的 io.Writer
}

// NewProtobufTxEncoder 函数创建一个新的 ProtobufTxEncoder 实例
// NewProtobufTxEncoder function creates a new instance of ProtobufTxEncoder
func NewProtobufTxEncoder(w io.Writer) *ProtobufTxEncoder {
	return &ProtobufTxEncoder{w: w}
}

// Encode 方法将交易数据编码为字节流
// Encode method encodes transaction data into a byte stream
func (enc *ProtobufTxEncoder) Encode(tx *Transaction) error {
	pbTx := &ProtoTransaction{
		Data:      tx.Data,
		From:      tx.From.ToSlice(),
		Signature: tx.Signature.ToBytes(),
		Hash:      tx.hash.ToSlice(),
		FirstSeen: tx.firstSeen,
	}
	data, err := proto.Marshal(pbTx)
	if err != nil {
		return err
	}
	_, err = enc.w.Write(data)
	return err
}

// ProtobufTxDecoder 结构体用于基于 Protobuf 的交易解码
// ProtobufTxDecoder struct is used for Protobuf-based transaction decoding
type ProtobufTxDecoder struct {
	r io.Reader // 用于读取编码数据的 io.Reader
}

// NewProtobufTxDecoder 函数创建一个新的 ProtobufTxDecoder 实例
// NewProtobufTxDecoder function creates a new instance of ProtobufTxDecoder
func NewProtobufTxDecoder(r io.Reader) *ProtobufTxDecoder {
	return &ProtobufTxDecoder{r: r}
}

// Decode 方法从字节流解码交易数据
// Decode method decodes transaction data from a byte stream
func (dec *ProtobufTxDecoder) Decode(tx *Transaction) error {
	data, err := io.ReadAll(dec.r)
	if err != nil {
		return err
	}
	pbTx := &ProtoTransaction{}
	if err := proto.Unmarshal(data, pbTx); err != nil {
		return err
	}
	tx.Data = pbTx.Data
	tx.From = crypto.PublicKeyFromBytes(pbTx.From)
	tx.Signature = crypto.SignatureFromBytes(pbTx.Signature)
	tx.hash = types.BytesToHash(pbTx.Hash)
	tx.firstSeen = pbTx.FirstSeen
	return nil
}
