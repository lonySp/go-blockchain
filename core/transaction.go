package core

import (
	"fmt"
	"github.com/lonySp/go-blockchain/crypto"
	"github.com/lonySp/go-blockchain/types"
)

// Transaction 结构体表示区块链中的交易
// Transaction struct represents a transaction in the blockchain
type Transaction struct {
	Data      []byte            // 交易数据
	From      crypto.PublicKey  // 发送方公钥
	Signature *crypto.Signature // 交易签名

	// tx 数据哈希的缓存版本
	// cached version of the tx data hash
	hash types.Hash
	// firstSeen 是本地首次看到该交易的时间戳
	// firstSeen is the timestamp of when this tx is first seen locally
	firstSeen int64
}

// NewTransaction 函数创建一个新的交易实例
// NewTransaction function creates a new instance of Transaction
func NewTransaction(data []byte) *Transaction {
	return &Transaction{
		Data: data,
	}
}

// Hash 方法计算并返回交易数据的哈希值
// Hash method computes and returns the hash of the transaction data
func (tx *Transaction) Hash(hasher Hasher[*Transaction]) types.Hash {
	if tx.hash.IsZero() {
		// 如果哈希值为空，则计算哈希值
		// Compute the hash if it is zero
		tx.hash = hasher.Hash(tx)
	}
	return tx.hash
}

// Sign 方法使用私钥对交易数据进行签名
// Sign method signs the transaction data using the private key
func (tx *Transaction) Sign(privateKey crypto.PrivateKey) error {
	// 使用私钥对交易数据进行签名
	// Sign the transaction data using the private key
	sig, err := privateKey.Sign(tx.Data)
	if err != nil {
		return err
	}
	// 设置发送方公钥和交易签名
	// Set the sender's public key and transaction signature
	tx.From = privateKey.PublicKey()
	tx.Signature = sig
	return nil
}

// Verify 方法验证交易签名的有效性
// Verify method verifies the validity of the transaction signature
func (tx *Transaction) Verify() error {
	// 如果签名为空，则返回错误
	// Return an error if the signature is nil
	if tx.Signature == nil {
		return fmt.Errorf("transaction has no signature")
	}
	// 验证签名，如果无效则返回错误
	// Verify the signature, return an error if invalid
	if !tx.Signature.Verify(tx.From, tx.Data) {
		return fmt.Errorf("invalid signature")
	}
	return nil
}

// Decode 方法从解码器中解码交易数据
// Decode method decodes the transaction data from the decoder
func (tx *Transaction) Decode(dec Decoder[*Transaction]) error {
	return dec.Decode(tx)
}

// Encode 方法将交易数据编码到编码器中
// Encode method encodes the transaction data into the encoder
func (tx *Transaction) Encode(enc Encoder[*Transaction]) error {
	return enc.Encode(tx)
}

// SetFirstSeen 方法设置首次看到交易的时间戳
// SetFirstSeen method sets the timestamp of when the transaction is first seen
func (tx *Transaction) SetFirstSeen(t int64) {
	tx.firstSeen = t
}

// FirstSeen 方法返回首次看到交易的时间戳
// FirstSeen method returns the timestamp of when the transaction was first seen
func (tx *Transaction) FirstSeen() int64 {
	return tx.firstSeen
}
