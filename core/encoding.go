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

// ProtobufBlockEncoder 结构体用于基于 Protobuf 的区块编码
// ProtobufBlockEncoder struct is used for Protobuf-based block encoding
type ProtobufBlockEncoder struct {
	w io.Writer // 用于写入编码数据的 io.Writer
}

// NewProtobufBlockEncoder 函数创建一个新的 ProtobufBlockEncoder 实例
// NewProtobufBlockEncoder function creates a new instance of ProtobufBlockEncoder
func NewProtobufBlockEncoder(w io.Writer) *ProtobufBlockEncoder {
	return &ProtobufBlockEncoder{w: w}
}

// Encode 方法将区块数据编码为字节流
// Encode method encodes block data into a byte stream
func (enc *ProtobufBlockEncoder) Encode(b *Block) error {
	pbBlock := &ProtoBlock{
		Header: &ProtoBlockHeader{
			Version:       b.Version,
			DataHash:      b.DataHash.ToSlice(),
			PrevBlockHash: b.PrevBlockHash.ToSlice(),
			Timestamp:     b.Timestamp,
			Height:        b.Height,
		},
		Validator:    b.Validator.ToSlice(),
		Signature:    b.Signature.ToBytes(),
		Hash:         b.hash.ToSlice(),
		Transactions: make([]*ProtoTransaction, len(b.Transactions)),
	}
	for i, tx := range b.Transactions {
		pbBlock.Transactions[i] = &ProtoTransaction{
			Data:      tx.Data,
			From:      tx.From.ToSlice(),
			Signature: tx.Signature.ToBytes(),
			Hash:      tx.hash.ToSlice(),
			FirstSeen: tx.firstSeen,
		}
	}
	data, err := proto.Marshal(pbBlock)
	if err != nil {
		return err
	}
	_, err = enc.w.Write(data)
	return err
}

// ProtobufBlockDecoder 结构体用于基于 Protobuf 的区块解码
// ProtobufBlockDecoder struct is used for Protobuf-based block decoding
type ProtobufBlockDecoder struct {
	r io.Reader // 用于读取编码数据的 io.Reader
}

// NewProtobufBlockDecoder 函数创建一个新的 ProtobufBlockDecoder 实例
// NewProtobufBlockDecoder function creates a new instance of ProtobufBlockDecoder
func NewProtobufBlockDecoder(r io.Reader) *ProtobufBlockDecoder {
	return &ProtobufBlockDecoder{r: r}
}

// Decode 方法从字节流解码区块数据
// Decode method decodes block data from a byte stream
func (dec *ProtobufBlockDecoder) Decode(b *Block) error {
	data, err := io.ReadAll(dec.r)
	if err != nil {
		return err
	}
	pbBlock := &ProtoBlock{}
	if err := proto.Unmarshal(data, pbBlock); err != nil {
		return err
	}
	b.Header = &Header{
		Version:       pbBlock.Header.Version,
		DataHash:      types.BytesToHash(pbBlock.Header.DataHash),
		PrevBlockHash: types.BytesToHash(pbBlock.Header.PrevBlockHash),
		Timestamp:     pbBlock.Header.Timestamp,
		Height:        pbBlock.Header.Height,
	}
	b.Validator = crypto.PublicKeyFromBytes(pbBlock.Validator)
	b.Signature = crypto.SignatureFromBytes(pbBlock.Signature)
	b.hash = types.BytesToHash(pbBlock.Hash)
	b.Transactions = make([]*Transaction, len(pbBlock.Transactions))
	for i, pbTx := range pbBlock.Transactions {
		b.Transactions[i] = &Transaction{
			Data:      pbTx.Data,
			From:      crypto.PublicKeyFromBytes(pbTx.From),
			Signature: crypto.SignatureFromBytes(pbTx.Signature),
			hash:      types.BytesToHash(pbTx.Hash),
			firstSeen: pbTx.FirstSeen,
		}
	}
	return nil
}
