package core

import (
	"fmt"
	"github.com/lonySp/go-blockchain/crypto"
	"github.com/lonySp/go-blockchain/types"
)

// Transaction 结构体表示区块链中的交易
// Transaction struct represents a transaction in the blockchain
type Transaction struct {
	Data       []byte
	From       crypto.PublicKey
	Signature  *crypto.Signature
	PrivateKey SerializablePrivateKey // 使用序列化私钥

	// cached version of the tx data hash
	hash types.Hash
	//firstSeen is the timestamp of when this tx is first seen locally
	firstSeen int64
}

func NewTransaction(data []byte) *Transaction {
	return &Transaction{
		Data: data,
	}
}

func (tx *Transaction) Hash(hasher Hasher[*Transaction]) types.Hash {
	if tx.hash.IsZero() {
		tx.hash = hasher.Hash(tx)
	}
	return tx.hash
}

// Sign 方法使用私钥对交易数据进行签名
// Sign method signs the transaction data using the private key
func (tx *Transaction) Sign(privateKey crypto.PrivateKey) error {
	sig, err := privateKey.Sign(tx.Data)
	if err != nil {
		return err
	}
	tx.From = privateKey.PublicKey()
	tx.Signature = sig
	return nil
}

// Verify 方法验证交易签名的有效性
// Verify method verifies the validity of the transaction signature
func (tx *Transaction) Verify() error {
	if tx.Signature == nil {
		return fmt.Errorf("transaction has no signature")
	}
	if !tx.Signature.Verify(tx.From, tx.Data) {
		return fmt.Errorf("invalid signature")
	}
	return nil
}

func (tx *Transaction) Decode(dec Decoder[*Transaction]) error {
	return dec.Decode(tx)
}
func (tx *Transaction) Encode(enc Encoder[*Transaction]) error {
	return enc.Encode(tx)
}
func (tx *Transaction) SetFirstSeen(t int64) {
	tx.firstSeen = t
}

func (tx *Transaction) FirstSeen() int64 {
	return tx.firstSeen
}
