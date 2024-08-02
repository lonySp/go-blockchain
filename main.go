package main

import (
	"bytes"
	"fmt"
	"github.com/lonySp/go-blockchain/core"
	"github.com/lonySp/go-blockchain/crypto"
	"github.com/lonySp/go-blockchain/network"
	"github.com/sirupsen/logrus"
	"log"
	"time"
)

func main() {
	// 创建本地和远程传输节点
	// Create local and remote transport nodes
	trLocal := network.NewLocalTransport("LOCAL")
	trRemoteA := network.NewLocalTransport("REMOTE_A")
	trRemoteB := network.NewLocalTransport("REMOTE_B")
	trRemoteC := network.NewLocalTransport("REMOTE_C")

	// 连接本地和远程传输节点
	// Connect local and remote transport nodes
	trLocal.Connect(trRemoteA)
	trRemoteA.Connect(trRemoteB)
	trRemoteB.Connect(trRemoteC)
	trRemoteA.Connect(trLocal)

	// 初始化远程服务器
	// Initialize remote servers
	initRemoteServers([]network.Transport{trRemoteA, trRemoteB, trRemoteC})

	// 启动一个 goroutine 每秒发送一笔交易
	// Start a goroutine to send a transaction every second
	go func() {
		for {
			// 发送交易到本地传输节点
			// Send transaction to the local transport node
			if err := sendTransaction(trRemoteA, trLocal.Addr()); err != nil {
				logrus.Error(err)
			}
			time.Sleep(2 * time.Second)
		}
	}()

	// 延迟连接新的远程节点
	// Delay connection of a new remote node
	//go func() {
	//	time.Sleep(7 * time.Second)
	//
	//	trLate := network.NewLocalTransport("LATE_REMOTE")
	//	trRemoteC.Connect(trLate)
	//	lateServer := makeServer("LATE_REMOTE", trLate, nil)
	//
	//	go lateServer.Start()
	//}()

	// 创建服务器选项并启动服务器
	// Create server options and start the server
	privateKey := crypto.GeneratePrivateKey()
	localServer := makeServer("LOCAL", trLocal, &privateKey)
	localServer.Start()
}

// initRemoteServers 初始化远程服务器
// initRemoteServers initializes remote servers
func initRemoteServers(trs []network.Transport) {
	for i := 0; i < len(trs); i++ {
		id := fmt.Sprintf("REMOTE_%d", i)
		s := makeServer(id, trs[i], nil)
		go s.Start()
	}
}

// makeServer 创建并返回一个新的服务器实例
// makeServer creates and returns a new server instance
func makeServer(id string, tr network.Transport, pk *crypto.PrivateKey) *network.Server {
	opts := network.ServerOpts{PrivateKey: pk, ID: id, Transport: []network.Transport{tr}}
	s, err := network.NewServer(opts)
	if err != nil {
		log.Fatal(err)
	}
	return s
}

// sendTransaction 函数生成并发送一笔交易
// sendTransaction function generates and sends a transaction
func sendTransaction(tr network.Transport, to network.NetAddr) error {
	// 生成私钥
	// Generate private key
	privateKey := crypto.GeneratePrivateKey()

	// 创建交易数据
	// Create transaction data
	data := []byte{0x03, 0x0a, 0x46, 0x0c, 0x4f, 0x0c, 0x4f, 0x0c, 0x0d, 0x05, 0x0a, 0x0f}

	// 创建新交易
	// Create a new transaction
	tx := core.NewTransaction(data)

	// 使用私钥签名交易
	// Sign the transaction with the private key
	tx.Sign(privateKey)

	// 创建缓冲区并编码交易
	// Create a buffer and encode the transaction
	buf := &bytes.Buffer{}
	if err := tx.Encode(core.NewProtobufTxEncoder(buf)); err != nil {
		return err
	}

	// 创建消息并发送
	// Create a message and send it
	msg := network.NewMessage(network.MessageTypeTx, buf.Bytes())
	return tr.SendMessage(to, msg.Bytes())
}
