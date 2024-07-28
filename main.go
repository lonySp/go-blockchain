package main

import (
	"bytes"
	"github.com/lonySp/go-blockchain/core"
	"github.com/lonySp/go-blockchain/crypto"
	"github.com/lonySp/go-blockchain/network"
	"github.com/sirupsen/logrus"
	"log"
	"math/rand"
	"strconv"
	"time"
)

func main() {
	// 创建本地和远程传输节点
	// Create local and remote transport nodes
	trLocal := network.NewLocalTransport("LOCAL")
	trRemote := network.NewLocalTransport("REMOTE")

	// 连接本地和远程传输节点
	// Connect local and remote transport nodes
	trLocal.Connect(trRemote)
	trRemote.Connect(trLocal)

	// 启动一个 goroutine 每秒发送一笔交易
	// Start a goroutine to send a transaction every second
	go func() {
		for {
			// 发送交易到本地传输节点
			// Send transaction to the local transport node
			if err := sendTransaction(trRemote, trLocal.Addr()); err != nil {
				logrus.Error(err)
			}
			time.Sleep(1 * time.Second)
		}
	}()

	// 创建服务器选项并启动服务器
	// Create server options and start the server
	privateKey := crypto.GeneratePrivateKey()
	opts := network.ServerOpts{PrivateKey: &privateKey, ID: "LOCAL", Transport: []network.Transport{trLocal}}
	server, err := network.NewServer(opts)
	if err != nil {
		log.Fatal(err)
	}
	server.Start()
}

// sendTransaction 函数生成并发送一笔交易
// sendTransaction function generates and sends a transaction
func sendTransaction(tr network.Transport, to network.NetAddr) error {
	// 生成私钥
	// Generate private key
	privateKey := crypto.GeneratePrivateKey()

	// 创建交易数据
	// Create transaction data
	data := []byte(strconv.FormatInt(int64(rand.Intn(1000)), 10))

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
