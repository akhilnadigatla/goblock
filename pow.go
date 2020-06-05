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

const targetBits = 24

type POW struct {
	block  *Block
	target *big.Int
}

func NewProof(block *Block) *POW {
	target := big.NewInt(1)
	target.Lsh(target, uint(256 - targetBits))

	pow := &POW{block, target}

	return pow
}

func (pow *POW) CompileData(nonce int) []byte {
	data := bytes.Join([][]byte{pow.block.PrevHash, 
				    pow.block.Data,
				    IntToHex(pow.block.Timestamp),
				    IntToHex(int64(targetBits)),
				    IntToHex(int64(nonce)),
				   }, []byte{},)
	return data
}

func (pow *POW) CheckPOW() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0
	
	fmt.Printf("Mining block with data '%s'\n", pow.block.Data)
	for nonce < maxNonce {
		data := pow.CompileData(nonce)
		
		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)
		hashInt.SetBytes(hash[:]) 
		
		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce ++
		}
	}
	fmt.Print("\n\n")
	
	return nonce, hash[:]
}

func (pow *POW) ValidatePOW() bool {
	var hashInt big.Int

	data := pow.CompileData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(pow.target) == -1

	return isValid
}
