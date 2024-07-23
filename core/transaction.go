package core

import "io"

// Transaction 结构体表示区块链中的交易
// Transaction struct represents a transaction in the blockchain
type Transaction struct {
	Data []byte
}

// DecodeBinary 从 io.Reader 解码交易数据
// DecodeBinary decodes transaction data from an io.Reader
func (tx Transaction) DecodeBinary(r io.Reader) error {
	return nil
}

// EncodeBinary 将交易数据编码到 io.Writer
// EncodeBinary encodes transaction data to an io.Writer
func (tx Transaction) EncodeBinary(w io.Writer) error {
	return nil
}
