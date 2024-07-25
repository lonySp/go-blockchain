package core

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/gob"
	"io"
	"math/big"
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

type GobTxEncoder struct {
	w io.Writer
}

func NewGobTxEncoder(w io.Writer) *GobTxEncoder {
	gob.Register(elliptic.P256())
	return &GobTxEncoder{w: w}
}

func (e *GobTxEncoder) Encode(tx *Transaction) error {
	return gob.NewEncoder(e.w).Encode(tx)
}

type GobTxDecoder struct {
	r io.Reader
}

func NewGobTxDecoder(r io.Reader) *GobTxDecoder {
	gob.Register(elliptic.P256())
	return &GobTxDecoder{r: r}
}

func (e *GobTxDecoder) Decode(tx *Transaction) error {
	return gob.NewDecoder(e.r).Decode(tx)
}

// SerializablePrivateKey 用于序列化和反序列化 ecdsa.PrivateKey
type SerializablePrivateKey struct {
	D, X, Y *big.Int
	Curve   string
}

func (spk *SerializablePrivateKey) FromECDSA(privateKey *ecdsa.PrivateKey) {
	spk.D = privateKey.D
	spk.X = privateKey.PublicKey.X
	spk.Y = privateKey.PublicKey.Y
	spk.Curve = "P256" // 假设只处理 P256 曲线
}

func (spk *SerializablePrivateKey) ToECDSA() *ecdsa.PrivateKey {
	var curve elliptic.Curve
	if spk.Curve == "P256" {
		curve = elliptic.P256()
	}
	return &ecdsa.PrivateKey{
		D: spk.D,
		PublicKey: ecdsa.PublicKey{
			Curve: curve,
			X:     spk.X,
			Y:     spk.Y,
		},
	}
}
