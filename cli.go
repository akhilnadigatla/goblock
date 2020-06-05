package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

type CLI struct {
	chain *Chain
}

func (cli *CLI) PrintUsage() {
	fmt.Println("Usage:")
	fmt.Println("	add -data <block data> -> adds a block to the chain")
	fmt.Println("	print -> print data of all blocks on the chain")
}

func (cli *CLI) ValidateArgs() {
	if len(os.Args) < 2 {
		cli.PrintUsage()
		os.Exit(1)
	}
}

func (cli *CLI) AddBlock(data string) {
	cli.chain.AddBlock(data)
	fmt.Println("Successfully added block.")
}

func (cli *CLI) PrintChain() {
	iter := cli.chain.Iterator()

	for {
		block := iter.Next()

		fmt.Printf("Previous Hash: %x\n", block.PrevHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Current Hash: %x\n", block.CurrHash)
		pow := NewProof(block)
		fmt.Printf("Proof-of-Work Valid: %s\n", strconv.FormatBool(pow.ValidatePOW()))
		fmt.Println()

		if len(block.PrevHash) == 0 {
			break
		}
	}
}

func (cli *CLI) Run() {
	cli.ValidateArgs()

	addCMD := flag.NewFlagSet("add", flag.ExitOnError)
	printCMD := flag.NewFlagSet("print", flag.ExitOnError)
	
	addData := addCMD.String("data", "", "Block Data")

	switch os.Args[1] {
	case "add":
		err := addCMD.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "print":
		err := printCMD.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.PrintUsage()
		os.Exit(1)
	}

	if addCMD.Parsed() {
		if *addData == "" {
			addCMD.Usage()
			os.Exit(1)
		}
		cli.AddBlock(*addData)
	}

	if printCMD.Parsed() {
		cli.PrintChain()
	}
}
