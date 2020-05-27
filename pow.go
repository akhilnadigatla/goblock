package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
)

var (
	maxNonce = math.MaxInt64
)

// Header value that indicates the difficulty at which the block was mined
const targetBits = 24

// Holds a proof-of-work for a given block
type pow struct {
	block *block
	target *big.Int
}

// Create the proof-of-work for a given block
func newProof(block *block) *pow {
	target := big.NewInt(1)
	target.Lsh(target, uint(256 - targetBits))
	
	pow := &pow{block, target}
	return pow
}

// Compile block data for proof inspection
func (pow *pow) compileData(nonce int) []byte {
	data := bytes.Join([][]byte{pow.block.prevHash, 
				    pow.block.data, 
				    intToHex(pow.block.timeStamp),
				    intToHex(int64(targetBits)),
				    intToHex(int64(nonce))}, []byte{})
	return data
}

// Run a proof-of-work check
func (pow *pow) checkPOW() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	fmt.Printf("Mining block with data '%s'\n", pow.block.data)	
	for nonce < maxNonce {
		data := pow.compileData(nonce)
		hash = sha256.Sum256(data)
		
		fmt.Printf("\r%x", hash)
		hashInt.SetBytes(hash[:])
		
		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce ++
		}
	}
	fmt.Printf("\n\n")

	return nonce, hash[:]
}

// Validate a block's proof-of-work
func (pow *pow) validatePOW() bool {
	var hashInt big.Int
	
	data := pow.compileData(pow.block.nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(pow.target) == -1

	return isValid	
}
