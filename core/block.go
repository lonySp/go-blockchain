package core

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/lonySp/go-blockchain/crypto"
	"github.com/lonySp/go-blockchain/types"
	"io"
)

// Header 结构体表示区块头
// Header struct represents the block header
type Header struct {
	Version       uint32
	PrevBlockHash types.Hash
	Timestamp     uint64
	Height        uint32
	Nonce         uint64
}

// Bytes 方法返回区块头的二进制数据
// Bytes method returns the binary data of the block header
func (h *Header) Bytes() []byte {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	enc.Encode(h)
	return buf.Bytes()
}

// Block 结构体表示区块
// Block struct represents the block
type Block struct {
	*Header
	Transaction []Transaction
	Validator   crypto.PublicKey
	Signature   *crypto.Signature

	// 缓存的区块头哈希
	// Cached version of the header hash
	hash types.Hash
}

// NewBlock 创建一个新的区块
// NewBlock creates a new block
func NewBlock(h *Header, txx []Transaction) *Block {
	return &Block{
		Header:      h,
		Transaction: txx,
	}
}

// AddTransaction 添加交易到区块
// AddTransaction adds a transaction to the block
func (b *Block) AddTransaction(tx *Transaction) {
	b.Transaction = append(b.Transaction, *tx)
}

// Sign 方法使用私钥对区块头数据进行签名
// Sign method signs the block header data using the private key
func (b *Block) Sign(privateKey crypto.PrivateKey) error {
	sig, err := privateKey.Sign(b.Header.Bytes())
	if err != nil {
		return err
	}

	b.Validator = privateKey.PublicKey()
	b.Signature = sig
	return nil
}

// Verify 方法验证区块签名的有效性
// Verify method verifies the validity of the block signature
func (b *Block) Verify() error {
	if b.Signature == nil {
		return fmt.Errorf("block has no signature")
	}

	if !b.Signature.Verify(b.Validator, b.Header.Bytes()) {
		return fmt.Errorf("block has invalid signature")
	}

	for _, tx := range b.Transaction {
		if err := tx.Verify(); err != nil {
			return err
		}
	}

	return nil
}

// Decode 方法解码区块
// Decode method decodes the block
func (b *Block) Decode(r io.Reader, dec Decoder[*Block]) error {
	return dec.Decode(r, b)
}

// Encode 方法编码区块
// Encode method encodes the block
func (b *Block) Encode(w io.Writer, enc Encoder[*Block]) error {
	return enc.Encode(w, b)
}

// Hash 方法计算区块的哈希值
// Hash method calculates the hash of the block
func (b *Block) Hash(hasher Hasher[*Header]) types.Hash {
	if b.hash.IsZero() {
		b.hash = hasher.Hash(b.Header)
	}
	return b.hash
}
