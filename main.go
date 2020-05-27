package main

import (
	"fmt"
	"strconv"
)

func main() {
	chain := newChain()

	chain.addBlock("Test 1")
	chain.addBlock("Test 2")

	for _, block := range chain.blocks {
		fmt.Printf("Previous Hash: %x\n", block.prevHash)
		fmt.Printf("Data: %s\n", block.data)
		fmt.Printf("Current Hash: %x\n", block.currHash)
		pow := newProof(block)
		fmt.Printf("Proof-of-Work: %s\n", strconv.FormatBool(pow.validatePOW()))
		fmt.Println()
	}
}
