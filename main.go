package main

import (
	"go-blockchain/network"
)

func main() {
	trLocal := network.NewLocalTransport("LOCAL")
	opts := network.ServerOpts{Transport: []network.Transport{trLocal}}

	server := network.NewServer(opts)
	server.Start()
}
