package main

import (
	"time"
)

// Block structure
type block struct {
	timeStamp int64  // Time of block creation
	data      []byte // Information to be contained by block
	prevHash  []byte // Hash of previous block
	currHash  []byte // Hash of defined block
	nonce	  int	 // Adjusted to ensure hash is less than target
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
