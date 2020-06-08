package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

type CLI struct {}

func (cli *CLI) CLCreateChain(address string) {
	chain := CreateChain(address)
	chain.db.Close()
	fmt.Println("Created chain.")
}

func (cli *CLI) GetBalance(address string) {
	chain := NewChain(address)
	defer chain.db.Close()
	
	balance := 0
	UTxnsO := chain.FindUTxnsO(address)

	for _, out := range UTxnsO {
		balance += out.Value
	}

	fmt.Printf("Balance of '%s': %d\n", address, balance)
}

func (cli *CLI) PrintUsage() {
	fmt.Println("Usage:")
	fmt.Println("	print -> print data of all blocks on the chain")
	fmt.Println("	balance -addr <address> -> get balance of specified address")
	fmt.Println("	create -addr <address> -> create chain and send genesis reward to specified address")
	fmt.Println("	send -from <from_address> -to <to_address> -amt <amount> -> send amount from address to another")
}

func (cli *CLI) ValidateArgs() {
	if len(os.Args) < 2 {
		cli.PrintUsage()
		os.Exit(1)
	}
}

func (cli *CLI) PrintChain() {
	chain := NewChain("")
	defer chain.db.Close()
	
	iter := chain.Iterator()

	for {
		block := iter.Next()

		fmt.Printf("Previous Hash: %x\n", block.PrevHash)
		fmt.Printf("Current Hash: %x\n", block.CurrHash)
		pow := NewProof(block)
		fmt.Printf("Proof-of-Work Valid: %s\n", strconv.FormatBool(pow.ValidatePOW()))
		fmt.Println()

		if len(block.PrevHash) == 0 {
			break
		}
	}
}

func (cli *CLI) Send(from, to string, amount int) {
	chain := NewChain(from)
	defer chain.db.Close()

	tx := NewUTOTxn(from, to, amount, chain)
	chain.MineBlock([]*Txn{tx})
	fmt.Println("Success.")
}

func (cli *CLI) Run() {
	cli.ValidateArgs()

	balanceCMD := flag.NewFlagSet("balance", flag.ExitOnError)
	createCMD := flag.NewFlagSet("create", flag.ExitOnError)
	sendCMD := flag.NewFlagSet("send", flag.ExitOnError)
	printCMD := flag.NewFlagSet("print", flag.ExitOnError)

	balanceAddr := balanceCMD.String("addr", "", "Address to get balance for.")
	createAddr:= createCMD.String("addr", "", "Address to send genesis reward to.")
	sendFrom := sendCMD.String("from", "", "Source wallet address.")
	sendTo := sendCMD.String("to", "", "Destination wallet address.")
	sendAmt := sendCMD.Int("amt", 0, "Amount to send.")	

	switch os.Args[1] {
	case "balance":
		err := balanceCMD.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}	
	case "create":
		err := createCMD.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "print":
		err := printCMD.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "send":
		err := sendCMD.Parse(os.Args[2:])
		if err != nil {	
			log.Panic(err)
		}
	default:
		cli.PrintUsage()
		os.Exit(1)
	}
	
	if balanceCMD.Parsed() {
		if *balanceAddr == "" {
			balanceCMD.Usage()
			os.Exit(1)
		}
		cli.GetBalance(*balanceAddr)
	}

	if createCMD.Parsed() {
		if *createAddr == "" {
			createCMD.Usage()
			os.Exit(1)
		}
		cli.CLCreateChain(*createAddr)
	}

	if printCMD.Parsed() {
		cli.PrintChain()
	}

	if sendCMD.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmt <= 0 {
			sendCMD.Usage()
			os.Exit(1)
		}
		cli.Send(*sendFrom, *sendTo, *sendAmt)
	}
}
