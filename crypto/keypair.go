package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"github.com/lonySp/go-blockchain/types"
	"math/big"
)

// PrivateKey 结构体表示一个私钥
// PrivateKey struct represents a private key
type PrivateKey struct {
	key *ecdsa.PrivateKey
}

// Sign 方法使用私钥对数据进行签名
// Sign method signs the data using the private key
func (k PrivateKey) Sign(data []byte) (*Signature, error) {
	r, s, err := ecdsa.Sign(rand.Reader, k.key, data)
	if err != nil {
		panic(err)
	}
	return &Signature{r, s}, nil
}

// GeneratePrivateKey 函数生成一个新的私钥
// GeneratePrivateKey function generates a new private key
func GeneratePrivateKey() PrivateKey {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	return PrivateKey{key}
}

// PublicKey 方法返回与私钥对应的公钥
// PublicKey method returns the public key corresponding to the private key
func (k PrivateKey) PublicKey() PublicKey {
	return PublicKey{Key: &k.key.PublicKey}
}

// PublicKey 结构体表示一个公钥
// PublicKey struct represents a public key
type PublicKey struct {
	Key *ecdsa.PublicKey
}

// ToSlice 方法将公钥转换为字节切片
// ToSlice method converts the public key to a byte slice
func (k PublicKey) ToSlice() []byte {
	return elliptic.MarshalCompressed(k.Key, k.Key.X, k.Key.Y)
}

// Address 方法生成与公钥对应的地址
// Address method generates the address corresponding to the public key
func (k PublicKey) Address() types.Address {
	h := sha256.Sum256(k.ToSlice())
	return types.NewAddressFromBytes(h[len(h)-20:])
}

// Signature 结构体表示一个签名
// Signature struct represents a signature
type Signature struct {
	R, S *big.Int
}

// Verify 方法验证签名是否有效
// Verify method verifies if the signature is valid
func (sig Signature) Verify(publicKey PublicKey, data []byte) bool {
	return ecdsa.Verify(publicKey.Key, data, sig.R, sig.S)
}
