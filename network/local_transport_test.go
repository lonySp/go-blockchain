package network

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

// TestConnect 测试 Connect 方法
// TestConnect tests the Connect method
func TestConnect(t *testing.T) {
	tra := NewLocalTransport("A") // 创建传输节点 A // Create transport node A
	trb := NewLocalTransport("B") // 创建传输节点 B // Create transport node B
	tra.Connect(trb)              // 连接 A 到 B // Connect A to B
	trb.Connect(tra)              // 连接 B 到 A // Connect B to A
	//assert.Equal(t, tra.peers[trb.Addr()], trb)
	//assert.Equal(t, trb.peers[tra.Addr()], tra)
}

// TestSendMessage 测试 SendMessage 方法
// TestSendMessage tests the SendMessage method
func TestSendMessage(t *testing.T) {
	tra := NewLocalTransport("A") // 创建传输节点 A // Create transport node A
	trb := NewLocalTransport("B") // 创建传输节点 B // Create transport node B
	tra.Connect(trb)              // 连接 A 到 B // Connect A to B
	trb.Connect(tra)              // 连接 B 到 A // Connect B to A

	msg := []byte("hello world")                    // 创建消息 // Create a message
	assert.Nil(t, tra.SendMessage(trb.Addr(), msg)) // 发送消息 // Send the message

	rpc := <-trb.Consume()                // 从 B 接收消息 // Receive the message from B
	b, err := ioutil.ReadAll(rpc.Payload) // 读取消息内容 // Read the message content
	assert.Nil(t, err)                    // 确认无错误 // Ensure no error
	assert.Equal(t, b, msg)               // 确认消息内容一致 // Ensure the message content is correct
	assert.Equal(t, rpc.From, tra.Addr()) // 确认发送者地址一致 // Ensure the sender address is correct
}

// TestBroadcast 测试 Broadcast 方法
// TestBroadcast tests the Broadcast method
func TestBroadcast(t *testing.T) {
	tra := NewLocalTransport("A") // 创建传输节点 A // Create transport node A
	trb := NewLocalTransport("B") // 创建传输节点 B // Create transport node B
	trc := NewLocalTransport("C") // 创建传输节点 C // Create transport node C

	tra.Connect(trb) // 连接 A 到 B // Connect A to B
	tra.Connect(trc) // 连接 A 到 C // Connect A to C

	msg := []byte("foo")              // 创建广播消息 // Create a broadcast message
	assert.Nil(t, tra.Broadcast(msg)) // 广播消息 // Broadcast the message

	rpcb := <-trb.Consume()                // 从 B 接收广播消息 // Receive the broadcast message from B
	b, err := ioutil.ReadAll(rpcb.Payload) // 读取消息内容 // Read the message content
	assert.Nil(t, err)                     // 确认无错误 // Ensure no error
	assert.Equal(t, b, msg)                // 确认消息内容一致 // Ensure the message content is correct

	rpcc := <-trc.Consume()               // 从 C 接收广播消息 // Receive the broadcast message from C
	b, err = ioutil.ReadAll(rpcc.Payload) // 读取消息内容 // Read the message content
	assert.Nil(t, err)                    // 确认无错误 // Ensure no error
	assert.Equal(t, b, msg)               // 确认消息内容一致 // Ensure the message content is correct
}
