A basic blockchain prototype built with Golang.

A blockchain is a dynamic collection of records that are referred to as blocks. Each block is identified by a header, which consists of:
- A timestamp
- Transaction data
- Hash value for the previous block in the chain

Therefore, a blockchain represents a distributed ledger of transactions that resides in public domain. This means that adding a new block requires significant work and permission of the keepers of this growing list.

This implementation uses the SHA-256 hashing algorithm to generate hash values for the blocks.
