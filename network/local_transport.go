package network

import (
	"fmt"
	"sync"
)

// LocalTransport 结构体表示本地传输
// LocalTransport struct represents local transport
type LocalTransport struct {
	addr       NetAddr
	consumerCh chan RPC
	lock       sync.RWMutex
	peers      map[NetAddr]*LocalTransport
}

// NewLocalTransport 创建并返回一个新的 LocalTransport 实例
// NewLocalTransport creates and returns a new LocalTransport instance
func NewLocalTransport(addr NetAddr) Transport {
	return &LocalTransport{
		addr:       addr,
		consumerCh: make(chan RPC, 1024),
		peers:      make(map[NetAddr]*LocalTransport),
	}
}

// Consume 方法返回一个 RPC 消费通道
// Consume method returns a RPC consumption channel
func (t *LocalTransport) Consume() <-chan RPC {
	return t.consumerCh
}

// Connect 方法连接两个传输节点
// Connect method connects two transport nodes
func (t *LocalTransport) Connect(tr Transport) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.peers[tr.Addr()] = tr.(*LocalTransport)
	return nil
}

// SendMessage 方法发送消息到指定地址
// SendMessage method sends a message to the specified address
func (t *LocalTransport) SendMessage(to NetAddr, payLoad []byte) error {
	t.lock.RLock()
	defer t.lock.RUnlock()

	peer, ok := t.peers[to]
	if !ok {
		return fmt.Errorf("%s: could not message to %s", t.addr, to)
	}

	peer.consumerCh <- RPC{
		From:    t.addr,
		Payload: payLoad,
	}

	return nil
}

// Addr 方法返回传输节点的地址
// Addr method returns the address of the transport node
func (t *LocalTransport) Addr() NetAddr {
	return t.addr
}
