package crypto

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestKeypairSignVerifySuccess 测试签名和验证成功的情况
// TestKeypairSignVerifySuccess tests the case where signing and verification are successful
func TestKeypairSignVerifySuccess(t *testing.T) {
	privateKey := GeneratePrivateKey()
	publicKey := privateKey.PublicKey()
	msg := []byte("hello world")
	signature, err := privateKey.Sign(msg)
	assert.Nil(t, err)
	assert.True(t, signature.Verify(publicKey, msg))
}

// TestKeypairSignVerifyFail 测试签名和验证失败的情况
// TestKeypairSignVerifyFail tests the case where signing and verification fail
func TestKeypairSignVerifyFail(t *testing.T) {
	privateKey := GeneratePrivateKey()
	PublicKey := privateKey.PublicKey()
	msg := []byte("hello world")
	signature, err := privateKey.Sign(msg)
	assert.Nil(t, err)

	otherPrivateKey := GeneratePrivateKey()
	otherPublicKey := otherPrivateKey.PublicKey()

	assert.False(t, signature.Verify(otherPublicKey, msg))
	assert.False(t, signature.Verify(PublicKey, []byte("xxxxxx")))
}
