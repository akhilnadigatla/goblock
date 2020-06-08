package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"time"
)

type Block struct {
	Nonce        int
	Timestamp    int64
	Transactions []*Txn
	PrevHash     []byte
	CurrHash     []byte
}

func (block *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range block.Transactions {
		txHashes = append(txHashes, tx.ID)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]
}

func (block *Block) BlockToBytes() []byte {
	var bytes bytes.Buffer
	encoder := gob.NewEncoder(&bytes)

	err := encoder.Encode(block)
	if err != nil {
		log.Panic(err)
	}

	return bytes.Bytes()
}

func BytesToBlock(seq []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(seq))

	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return &block
}

func NewBlock(transactions []*Txn, prevHash []byte) *Block {
	block := &Block{0, time.Now().Unix(), transactions, prevHash, []byte{}}
	pow := NewProof(block)
	nonce, hash := pow.CheckPOW()

	block.Nonce = nonce
	block.CurrHash = hash[:]

	return block
}

func NewGenesisBlock(coinbase *Txn) *Block {
	return NewBlock([]*Txn{coinbase}, []byte{})
}
