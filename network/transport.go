package network

type NetAddr string

// RPC 是在节点之间传递的消息
// RPC is a message that is sent between nodes
type RPC struct {
	From    NetAddr
	Payload []byte
}

// Transport 是一个网络抽象，允许我们在节点之间发送和接收 RPC
// Transport is a network abstraction that allows us to send and receive RPCs between nodes
type Transport interface {
	Consume() <-chan RPC
	Connect(Transport) error
	SendMessage(NetAddr, []byte) error
	Addr() NetAddr
}
