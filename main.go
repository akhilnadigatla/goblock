package main

func main() {
	chain := NewChain()
	defer chain.db.Close()

	cli := CLI{chain}
	cli.Run()
}
