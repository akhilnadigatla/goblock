package main

import (
	"fmt"
	"strconv"
)

func main() {
	chain := newChain()
	defer chain.db.Close()

	cli := CLI{chain}
	cli.Run()
}
