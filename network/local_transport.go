package network

import (
	"fmt"
	"sync"
)

type LocalTransport struct {
	addr       NetAddr
	consumerCh chan RPC
	lock       sync.RWMutex
	peers      map[NetAddr]*LocalTransport
}

func NewLocalTransport(addr NetAddr) Transport {
	return &LocalTransport{
		addr:       addr,
		consumerCh: make(chan RPC, 1024),
		peers:      make(map[NetAddr]*LocalTransport),
	}
}

func (t *LocalTransport) Consume() <-chan RPC {
	return t.consumerCh
}

func (t *LocalTransport) Connect(tr Transport) error {
	// ps: Lock 表示写操作 Write_Lock
	t.lock.Lock()
	defer t.lock.Unlock()

	t.peers[tr.Addr()] = tr.(*LocalTransport)
	return nil
}

func (t *LocalTransport) SendMessage(to NetAddr, payLoad []byte) error {
	// ps: RLock 表示读操作 Read_Lock
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

func (t *LocalTransport) Addr() NetAddr {
	return t.addr
}
