package main

import (
	"fmt"
	"log"
	"github.com/boltdb/bolt"
)

const dbFile = "chain.db"
const bucket = "blocks"

type Chain struct {
	tip []byte
	db  *bolt.DB
}

type Iterator struct {
	currHash []byte
	db	 *bolt.DB
}

func (chain *Chain) AddBlock(data string) {
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
	
	newBlock := NewBlock(data, lastHash)

	err = chain.db.Update(
	func (tx *bolt.Tx) error {
		b:= tx.Bucket([]byte(bucket))
		
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

func NewChain() *Chain {
	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(
	func (tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		
		if b == nil {
			fmt.Println("No existing blockchain found. Creating a new one...")
			genesis := NewGenesisBlock()

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
		} else {
			tip = b.Get([]byte("l"))
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	chain := Chain{tip, db}

	return &chain
}
