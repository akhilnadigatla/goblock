package main

// Chain composed of a list of blocks
type chain struct {
	blocks []*block
}

// Add a new block to the end of the chain
func (chain *chain) addBlock(data string) {
	lastBlock := chain.blocks[len(chain.blocks) - 1]
	newBlock := newBlock(data, lastBlock.hash)
	chain.blocks = append(chain.blocks, newBlock)
}

// Create a new blockchain
fun newChain() *chain {
	return &chain{[]*block{newFirstBlock()}}
}
