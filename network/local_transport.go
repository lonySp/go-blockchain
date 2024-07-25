package network

import (
	"fmt"
	"sync"
)

// LocalTransport 结构体表示本地传输
// LocalTransport struct represents local transport
type LocalTransport struct {
	addr       NetAddr                     // 传输节点的地址 // Address of the transport node
	consumerCh chan RPC                    // 用于接收 RPC 消息的通道 // Channel for receiving RPC messages
	lock       sync.RWMutex                // 读写锁，用于并发访问 // Read-write lock for concurrent access
	peers      map[NetAddr]*LocalTransport // 已连接的传输节点 // Connected transport nodes
}

// NewLocalTransport 创建并返回一个新的 LocalTransport 实例
// NewLocalTransport creates and returns a new LocalTransport instance
func NewLocalTransport(addr NetAddr) Transport {
	return &LocalTransport{
		addr:       addr,
		consumerCh: make(chan RPC, 1024),              // 创建一个带缓冲区的通道 // Create a buffered channel
		peers:      make(map[NetAddr]*LocalTransport), // 初始化已连接节点的映射 // Initialize the map of connected nodes
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

	// 将另一个传输节点添加到 peers 映射中
	// Add the other transport node to the peers map
	t.peers[tr.Addr()] = tr.(*LocalTransport)
	return nil
}

// SendMessage 方法发送消息到指定地址
// SendMessage method sends a message to the specified address
func (t *LocalTransport) SendMessage(to NetAddr, payLoad []byte) error {
	t.lock.RLock()
	defer t.lock.RUnlock()

	// 查找目标地址的传输节点
	// Find the transport node for the target address
	peer, ok := t.peers[to]
	if !ok {
		return fmt.Errorf("%s: could not message to %s", t.addr, to)
	}

	// 发送 RPC 消息到目标节点的 consumerCh 通道
	// Send the RPC message to the target node's consumerCh channel
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
