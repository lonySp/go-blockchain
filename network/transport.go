package network

type NetAddr string

/*
RPC is a message that is sent between nodes.
*/
type RPC struct {
	From    NetAddr
	Payload []byte
}

/*
Transport is a network abstraction that allows us to send and receive RPCs
between nodes.
*/
type Transport interface {
	Consume() <-chan RPC
	Connect(Transport) error
	SendMessage(NetAddr, []byte) error
	Addr() NetAddr
}
