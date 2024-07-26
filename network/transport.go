package network

type NetAddr string

// Transport 是一个网络抽象，允许我们在节点之间发送和接收 RPC
// Transport is a network abstraction that allows us to send and receive RPCs between nodes
type Transport interface {
	Consume() <-chan RPC
	Connect(Transport) error
	SendMessage(NetAddr, []byte) error
	Broadcast([]byte) error
	Addr() NetAddr
}
