package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"github.com/lonySp/go-blockchain/crypto"
	"github.com/lonySp/go-blockchain/types"
	"time"
)

// Header 结构体表示区块头
// Header struct represents the block header
type Header struct {
	Version       uint32     // 区块版本号 // Block version number
	DataHash      types.Hash // 区块数据哈希值 // Hash of the block data
	PrevBlockHash types.Hash // 前一个区块的哈希值 // Hash of the previous block
	Timestamp     uint64     // 区块生成的时间戳 // Timestamp when the block was created
	Height        uint32     // 区块高度 // Block height
}

// Bytes 方法返回区块头的二进制数据
// Bytes method returns the binary data of the block header
func (h *Header) Bytes() []byte {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	enc.Encode(h) // 将区块头编码为二进制数据 // Encode the block header into binary data
	return buf.Bytes()
}

// Block 结构体表示区块
// Block struct represents the block
type Block struct {
	*Header                        // 区块头 // Block header
	Transactions []*Transaction    // 区块包含的交易列表 // List of transactions included in the block
	Validator    crypto.PublicKey  // 验证者公钥 // Validator's public key
	Signature    *crypto.Signature // 区块签名 // Block signature

	// 缓存的区块头哈希 // Cached version of the header hash
	hash types.Hash
}

// NewBlock 创建一个新的区块
// NewBlock creates a new block
func NewBlock(h *Header, txx []*Transaction) (*Block, error) {
	return &Block{
		Header:       h,
		Transactions: txx,
	}, nil
}

// NewBlockFromPrevHeader 基于前一个区块头创建一个新的区块
// NewBlockFromPrevHeader creates a new block based on the previous block header
func NewBlockFromPrevHeader(prevHeader *Header, txx []*Transaction) (*Block, error) {
	// 计算数据哈希 // Calculate data hash
	dataHash, err := CalculateDataHash(txx)
	if err != nil {
		return nil, err
	}
	// 创建新的区块头 // Create a new block header
	header := &Header{
		Version:       1,
		Height:        prevHeader.Height + 1,
		DataHash:      dataHash,
		PrevBlockHash: BlockHasher{}.Hash(prevHeader),
		Timestamp:     uint64(time.Now().UnixNano()),
	}
	// 创建并返回新的区块 // Create and return the new block
	return NewBlock(header, txx)
}

// AddTransaction 添加交易到区块
// AddTransaction adds a transaction to the block
func (b *Block) AddTransaction(tx *Transaction) {
	b.Transactions = append(b.Transactions, tx)
}

// Sign 方法使用私钥对区块头数据进行签名
// Sign method signs the block header data using the private key
func (b *Block) Sign(privateKey crypto.PrivateKey) error {
	// 使用私钥对区块头进行签名 // Sign the block header using the private key
	sig, err := privateKey.Sign(b.Header.Bytes())
	if err != nil {
		return err
	}

	// 设置验证者公钥和区块签名 // Set the validator's public key and block signature
	b.Validator = privateKey.PublicKey()
	b.Signature = sig
	return nil
}

// Verify 方法验证区块签名的有效性
// Verify method verifies the validity of the block signature
func (b *Block) Verify() error {
	// 如果签名为空，则返回错误 // Return an error if the signature is nil
	if b.Signature == nil {
		return fmt.Errorf("block has no signature")
	}

	// 验证签名，如果无效则返回错误 // Verify the signature, return an error if invalid
	if !b.Signature.Verify(b.Validator, b.Header.Bytes()) {
		return fmt.Errorf("block has invalid signature")
	}

	// 验证区块中的每个交易 // Verify each transaction in the block
	for _, tx := range b.Transactions {
		if err := tx.Verify(); err != nil {
			return err
		}
	}

	// 重新计算数据哈希并验证 // Recalculate data hash and verify
	dataHash, err := CalculateDataHash(b.Transactions)
	if err != nil {
		return nil
	}
	if dataHash != b.DataHash {
		return fmt.Errorf("block (%s) has an invalid data hash", b.DataHash)
	}

	return nil
}

// Decode 方法解码区块
// Decode method decodes the block
func (b *Block) Decode(dec Decoder[*Block]) error {
	return dec.Decode(b)
}

// Encode 方法编码区块
// Encode method encodes the block
func (b *Block) Encode(enc Encoder[*Block]) error {
	return enc.Encode(b)
}

// Hash 方法计算区块的哈希值
// Hash method calculates the hash of the block
func (b *Block) Hash(hasher Hasher[*Header]) types.Hash {
	// 如果哈希值为空，则计算哈希值 // Compute the hash if it is zero
	if b.hash.IsZero() {
		b.hash = hasher.Hash(b.Header)
	}
	return b.hash
}

// CalculateDataHash 计算交易列表的数据哈希值
// CalculateDataHash calculates the data hash of the transaction list
func CalculateDataHash(txx []*Transaction) (hash types.Hash, err error) {
	// 创建一个字节缓冲区 // Create a byte buffer
	buf := &bytes.Buffer{}
	// 将每个交易编码为二进制数据并写入缓冲区 // Encode each transaction into binary data and write to the buffer
	for _, tx := range txx {
		if err := tx.Encode(NewProtobufTxEncoder(buf)); err != nil {
			return types.Hash{}, err
		}
	}

	// 计算缓冲区内容的 SHA-256 哈希值 // Calculate the SHA-256 hash of the buffer contents
	hash = sha256.Sum256(buf.Bytes())
	return
}
