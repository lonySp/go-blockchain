package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"go-blockchain/types"
	"io"
)

// Header 结构体表示区块头
// Header struct represents the block header
type Header struct {
	Version   uint32
	PrevBlock types.Hash
	Timestamp uint64
	Height    uint32
	Nonce     uint64
}

// EnCodeBinary 方法将区块头编码到 io.Writer
// EnCodeBinary method encodes the block header to an io.Writer
func (h *Header) EnCodeBinary(w io.Writer) error {
	if err := binary.Write(w, binary.LittleEndian, &h.Version); err != nil {
		return err
	}
	if err := binary.Write(w, binary.LittleEndian, &h.PrevBlock); err != nil {
		return err
	}
	if err := binary.Write(w, binary.LittleEndian, &h.Timestamp); err != nil {
		return err
	}
	if err := binary.Write(w, binary.LittleEndian, &h.Height); err != nil {
		return err
	}
	return binary.Write(w, binary.LittleEndian, &h.Nonce)
}

// DecodeBinary 方法从 io.Reader 解码区块头
// DecodeBinary method decodes the block header from an io.Reader
func (h *Header) DecodeBinary(r io.Reader) error {
	if err := binary.Read(r, binary.LittleEndian, &h.Version); err != nil {
		return err
	}
	if err := binary.Read(r, binary.LittleEndian, &h.PrevBlock); err != nil {
		return err
	}
	if err := binary.Read(r, binary.LittleEndian, &h.Timestamp); err != nil {
		return err
	}
	if err := binary.Read(r, binary.LittleEndian, &h.Height); err != nil {
		return err
	}
	return binary.Read(r, binary.LittleEndian, &h.Nonce)
}

// Block 结构体表示区块
// Block struct represents the block
type Block struct {
	Header
	Transaction []Transaction

	// 缓存的区块头哈希
	// Cached version of the header hash
	hash types.Hash
}

// Hash 方法计算并返回区块头的哈希值
// Hash method calculates and returns the hash of the block header
func (b *Block) Hash() types.Hash {
	buf := &bytes.Buffer{}
	b.Header.EnCodeBinary(buf)

	if b.hash.IsZero() {
		b.hash = types.Hash(sha256.Sum256(buf.Bytes()))
	}

	return b.hash
}

// EncodeBinary 方法将区块编码到 io.Writer
// EncodeBinary method encodes the block to an io.Writer
func (b *Block) EncodeBinary(w io.Writer) error {
	if err := b.Header.EnCodeBinary(w); err != nil {
		return err
	}
	for _, tx := range b.Transaction {
		if err := tx.EncodeBinary(w); err != nil {
			return err
		}
	}
	return nil
}

// DecodeBinary 方法从 io.Reader 解码区块
// DecodeBinary method decodes the block from an io.Reader
func (b *Block) DecodeBinary(r io.Reader) error {
	if err := b.Header.DecodeBinary(r); err != nil {
		return err
	}
	for _, tx := range b.Transaction {
		if err := tx.DecodeBinary(r); err != nil {
			return err
		}
	}
	return nil
}
