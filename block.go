package main

import (
	"bytes"
	"crypto/sha256"
	"strconv"
	"time"
)

// Block structure
type block struct {
	timeStamp int64   // Time of block creation
	data      []byte  // Information to be contained by block
	prevHash  []byte  // Hash of previous block
	currHash  []byte  // Hash of defined block
}

// Set hash value for a block
func (block *block) setHash() {
	timeStamp := []byte(strconv.FormatInt(b.timeStamp, 10))
	headers := bytes.Join([][]byte{timeStamp, b.data, b.prevHash}, []byte{})
	hash := sha256.Sum256(headers)

	b.Hash = hash[:]
} 

// Create a new block
func newBlock(data string, prevHash []byte) *block {
	block := &block{time.Now().Unix(), []byte(data), prevHash, []byte{}}
	block.setHash()
	return block 
}

// Create first block in a chain
func newFirstBlock() *block {
	return newBlock("First block", []byte{})
}
