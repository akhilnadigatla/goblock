package main

import (
	"fmt"
)

func main() {
	chain := newChain()

	chain.addBlock("Test 1")
	chain.addBlock("Test 2")

	for _, block := range chain.blocks {
		fmt.Printf("Previous Hash: %x\n", block.prevHash)
		fmt.Printf("Current Hash: %x\n", block.currHash)
		fmt.Printf("Data: %s\n", block.data)
		fmt.Println()
	}
}
