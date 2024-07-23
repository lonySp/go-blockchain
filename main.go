package main

import (
	"github.com/lonySp/go-blockchain/network"
	"time"
)

func main() {
	trLocal := network.NewLocalTransport("LOCAL")
	trRemote := network.NewLocalTransport("REMOTE")

	trLocal.Connect(trRemote)
	trRemote.Connect(trLocal)

	go func() {
		for {
			trRemote.SendMessage(trLocal.Addr(), []byte("hello World!"))
			time.Sleep(1 * time.Second)
		}

	}()

	opts := network.ServerOpts{Transport: []network.Transport{trLocal}}

	server := network.NewServer(opts)
	server.Start()
}
