package core

import (
	"fmt"
	"github.com/lonySp/go-blockchain/crypto"
)

// Transaction 结构体表示区块链中的交易
// Transaction struct represents a transaction in the blockchain
type Transaction struct {
	Data      []byte
	From      crypto.PublicKey
	Signature *crypto.Signature
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
