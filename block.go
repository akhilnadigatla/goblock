package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"
)

type Block struct {
	Nonce     int
	Timestamp int64
	Data      []byte
	PrevHash  []byte
	CurrHash  []byte
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

func NewBlock(data string, prevHash []byte) *Block {
	block := &Block{0, time.Now().Unix(), []byte(data), prevHash, []byte{}}
	pow := NewProof(block)
	nonce, hash := pow.CheckPOW()

	block.Nonce = nonce
	block.CurrHash = hash[:]

	return block
}

func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}
