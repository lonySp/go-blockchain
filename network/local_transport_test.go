package network

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestConnect 测试 Connect 方法
// TestConnect tests the Connect method
func TestConnect(t *testing.T) {
	tra := NewLocalTransport("A")
	trb := NewLocalTransport("B")
	tra.Connect(trb)
	trb.Connect(tra)
	//assert.Equal(t, tra.peers[trb.Addr()], trb)
	//assert.Equal(t, trb.peers[tra.Addr()], tra)
}

// TestSendMessage 测试 SendMessage 方法
// TestSendMessage tests the SendMessage method
func TestSendMessage(t *testing.T) {
	tra := NewLocalTransport("A")
	trb := NewLocalTransport("B")
	tra.Connect(trb)
	trb.Connect(tra)

	msg := []byte("hello world")
	assert.Nil(t, tra.SendMessage(trb.Addr(), msg))

	rpc := <-trb.Consume()
	assert.Equal(t, rpc.Payload, msg)
	assert.Equal(t, rpc.From, tra.Addr())
}
