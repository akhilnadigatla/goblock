package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"github.com/boltdb/bolt"
)

const dbFile = "chain.db"
const bucket = "blocks"
const genesisCoinbaseData = "Genesis Data"

type Chain struct {
	tip []byte
	db  *bolt.DB
}

type Iterator struct {
	currHash []byte
	db	 *bolt.DB
}

func (chain *Chain) MineBlock(transactions []*Txn) {
	var lastHash []byte
	
	err := chain.db.View(
	func (tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		lastHash = b.Get([]byte("l"))
		
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	
	newBlock := NewBlock(transactions, lastHash)

	err = chain.db.Update(
	func (tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		
		err := b.Put(newBlock.CurrHash, newBlock.BlockToBytes())
		if err != nil {
			log.Panic(err)
		}
		
		err = b.Put([]byte("l"), newBlock.CurrHash)
		if err != nil {
			log.Panic(err)
		}

		chain.tip = newBlock.CurrHash
		
		return nil
	})
}

func (chain *Chain) FindUTxns(address string) []Txn {
	var unspent []Txn
	spent := make(map[string][]int)
	iter := chain.Iterator()

	for {
		block := iter.Next()
		
		for _, txn := range block.Transactions {
			txnID := hex.EncodeToString(txn.ID)
		
		Outputs:
			for outIdx, out := range txn.Vout {
				if spent[txnID] != nil {
					for _, spentOut := range spent[txnID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}
				
				if out.CanBeUnlockedWith(address) {
					unspent = append(unspent, *txn)
				}
			}

			if txn.IsCoinbase() == false {
				for _, in := range txn.Vin {
					if in.CanUnlockOutputWith(address) {
						inTxnID := hex.EncodeToString(in.Txid)
						spent[inTxnID] = append(spent[inTxnID], in.Vout)
					}
				}
			}
		
		}

		if len(block.PrevHash) == 0 {
			break
		}
	}	

	return unspent
}

func (chain *Chain) FindUTxnsO(address string) []TxOutput {
	var UTxnsO []TxOutput
	unspent := chain.FindUTxns(address)

	for _, txn := range unspent {
		for _, out := range txn.Vout {
			if out.CanBeUnlockedWith(address) {
				UTxnsO = append(UTxnsO, out)
			}
		}
	}

	return UTxnsO
}

func (chain *Chain) FindSTxnsO(address string, amount int) (int, map[string][]int) {
	UTxnsO := make(map[string][]int)
	UTxns := chain.FindUTxns(address)
	acc := 0

Work:
	for _, txn := range UTxns {
		txnID := hex.EncodeToString(txn.ID)
		
		for outIdx, out := range txn.Vout {
			if out.CanBeUnlockedWith(address) && acc < amount {
				acc += out.Value
				UTxnsO[txnID] = append(UTxnsO[txnID], outIdx)
			
				if acc >= amount {
					break Work
				}
			}
		}
	}
	
	return acc, UTxnsO
}

func (chain *Chain) Iterator() *Iterator {
	iter := &Iterator{chain.tip, chain.db}
	return iter
}

func (iter *Iterator) Next() *Block {
	var block *Block
	
	err := iter.db.View(
	func (tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		encBlock := b.Get(iter.currHash)
		block = BytesToBlock(encBlock)
		
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	iter.currHash = block.PrevHash
	
	return block
}

func dbExists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}
	
	return true
}

func NewChain(address string) *Chain {
	if dbExists() == false {
		fmt.Println("No existing blockchain found.")
		os.Exit(1)
	}
	
	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	
	err = db.Update(
	func (tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		tip = b.Get([]byte("l"))
		
		return nil	
	})
	if err != nil {
		log.Panic(err)
	}

	chain := Chain{tip, db}
	
	return &chain
}

func CreateChain(address string) *Chain {
	if dbExists() {
		fmt.Println("Blockchain already exists")
		os.Exit(1)
	}

	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	
	err = db.Update(
	func (tx *bolt.Tx) error {
		cbtx := NewCoinbaseTxn(address, genesisCoinbaseData)
		genesis := NewGenesisBlock(cbtx)

		b, err := tx.CreateBucket([]byte(bucket))
		if err != nil {
			log.Panic(err)
		}

		err = b.Put(genesis.CurrHash, genesis.BlockToBytes())
		if err != nil {
			log.Panic(err)
		}

		err = b.Put([]byte("l"), genesis.CurrHash)
		if err != nil {
			log.Panic(err)
		}

		tip = genesis.CurrHash
		
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	chain := Chain{tip, db}

	return &chain
}
