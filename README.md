# go-blockchain

go-blockchain is a blockchain implementation in Golang.

## Roadmap

1. Build the underlying network transport layer(P2P).
2. Add block scaffolding to realize data encoding, decoding and hash generation.
3. A key pair is generated using an elliptic curve encryption algorithm and is used for signing and verifying signatures.
4. Implement transaction handling with signing and verification to ensure secure and valid transactions within the blockchain.
5. Block signing and verification, blockchain structure and storage, and block addition and validation functionality in a blockchain system.
6. Added validator, block validation, and transaction validation.
7. Adding Mutex Locks & Expanding Transaction Pool with Tests.
8. The added RPC functionality includes processing decoded messages and broadcasting transactions.
9. Change the encoding mode to protobuf
10. Connect local server transport nodes, start the server, create a transaction every second, and create a block every five seconds.
11. Reconstructing the transaction pool, adding block broadcasting
12. impl VM basic (supporting basic arithmetic operations and stack manipulation)
13. Implement virtual machine contract state transition



## News
1.log change => kitlog
2.Change the encoding mode to protobuf


## bug:
    code => transaction_test.go => TestTxEncodeDecode => BUG(waiting for a solution)
    type not registered for interface: elliptic.p256Curve 
    solution:
    Change the encoding mode to protobuf