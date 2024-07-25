# go-blockchain

go-blockchain is a blockchain implementation in Golang.

## Roadmap

1. Build the underlying network transport layer(P2P).
2. Add block scaffolding to realize data encoding, decoding and hash generation.
3. A key pair is generated using an elliptic curve encryption algorithm and is used for signing and verifying signatures.
4. Implement transaction handling with signing and verification to ensure secure and valid transactions within the blockchain.
5. Block signing and verification, blockchain structure and storage, and block addition and validation functionality in a blockchain system.
6. Added validator, block validation, and transaction validation.

TODO : 7. 添加互斥锁&Expected nil, but got: &errors.errorString{s:"gob: type not registered for interface: elliptic.p256Curve"}
