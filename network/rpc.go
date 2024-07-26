package network

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/lonySp/go-blockchain/core"
	"github.com/sirupsen/logrus"
	"io"
)

// MessageType 表示消息类型
// MessageType represents the type of a message
type MessageType byte

const (
	MessageTypeTx    MessageType = 0x1 // 交易消息类型 // Transaction message type
	MessageTypeBlock                   // 区块消息类型 // Block message type
)

// RPC 结构体表示一个远程过程调用
// RPC struct represents a remote procedure call
type RPC struct {
	From    NetAddr   // 消息发送者的地址 // Address of the message sender
	Payload io.Reader // 消息的有效负载 // Payload of the message
}

// Message 结构体表示一个网络消息
// Message struct represents a network message
type Message struct {
	Header MessageType // 消息的类型 // Type of the message
	Data   []byte      // 消息的数据 // Data of the message
}

// NewMessage 创建并返回一个新的 Message 实例
// NewMessage creates and returns a new instance of Message
func NewMessage(t MessageType, data []byte) *Message {
	return &Message{
		Header: t,
		Data:   data,
	}
}

// Bytes 方法将消息编码为字节切片
// Bytes method encodes the message into a byte slice
func (msg *Message) Bytes() []byte {
	buf := &bytes.Buffer{}
	gob.NewEncoder(buf).Encode(msg) // 使用 gob 编码消息 // Encode the message using gob
	return buf.Bytes()
}

// DecodedMessage 结构体表示一个解码后的消息
// DecodedMessage struct represents a decoded message
type DecodedMessage struct {
	From NetAddr // 消息发送者的地址 // Address of the message sender
	Data any     // 解码后的数据 // Decoded data
}

// RPCDecodeFunc 是一个用于解码 RPC 消息的函数类型
// RPCDecodeFunc is a function type for decoding RPC messages
type RPCDecodeFunc func(rpc RPC) (*DecodedMessage, error)

// DefaultRPCDecodeFunc 是默认的 RPC 解码函数
// DefaultRPCDecodeFunc is the default RPC decode function
func DefaultRPCDecodeFunc(rpc RPC) (*DecodedMessage, error) {
	msg := Message{}
	if err := gob.NewDecoder(rpc.Payload).Decode(&msg); err != nil {
		return nil, fmt.Errorf("failed to decode message from %s: %s", rpc.From, err)
	}

	logrus.WithFields(logrus.Fields{
		"type": msg.Header,
		"from": rpc.From,
	}).Debug("new incoming message")

	switch msg.Header {
	case MessageTypeTx:
		tx := new(core.Transaction)
		if err := tx.Decode(core.NewGobTxDecoder(bytes.NewReader(msg.Data))); err != nil {
			return nil, err
		}

		return &DecodedMessage{
			From: rpc.From,
			Data: tx,
		}, nil
	default:
		return nil, fmt.Errorf("invalid message type %x", msg.Header)
	}
}

// RPCProcessor 接口定义了处理解码消息的方法
// RPCProcessor interface defines a method for processing decoded messages
type RPCProcessor interface {
	ProcessMessage(message *DecodedMessage) error
}
