package main

import (
	"time"
	"bytes"
	"encoding/gob"
	"log"
)

// Block structure
type block struct {
	timeStamp int64  // Time of block creation
	data      []byte // Information to be contained by block
	prevHash  []byte // Hash of previous block
	currHash  []byte // Hash of defined block
	nonce	  int	 // Adjusted to ensure hash is less than target
}

func (block *block) blockToBytes() []byte {
	var bytes bytes.Buffer
	encoder := gob.NewEncoder(&bytes)

	err := encoder.Encode(block)
	if err != nil {
		log.Panic(err)
	}

	return bytes.Bytes()
}

func bytesToBlock(seq []byte) *block {
	var block block
	decoder := gob.NewDecoder(bytes.NewReader(seq))

	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return &block
}

// Create a new block
func newBlock(data string, prevHash []byte) *block {
	block := &block{time.Now().Unix(), []byte(data), prevHash, []byte{}, 0}
	pow := newProof(block)
	nonce, hash := pow.checkPOW()	

	block.currHash = hash[:]
	block.nonce = nonce

	return block
}

// Create first block in a chain
func newFirstBlock() *block {
	return newBlock("First block", []byte{})
}
