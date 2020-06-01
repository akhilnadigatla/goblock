package main

import (
	"fmt"
	"log"
	"github.com/boltdb/bolt"
)

const dbFile = "chain.db"
const bucket = "blocks"

// Chain composed of a list of blocks
type chain struct {
	tip []byte
	db  *bolt.DB
}

// Structurre to iterate over blocks
type iterator struct {
	currHash []byte
	db 	 *bolt.DB
}

// Add a new block to the end of the chain
func (chain *chain) addBlock(data string) {
	var lastHash []byte
	
	err := chain.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		lastHash = b.Get([]byte("l"))
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	newBlock := newBlock(data, lastHash)

	err := chain.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		err := b.Put(newBlock.currHash, newBlock.blockToBytes())
		if err != nil {
			log.Panic(err)
		}
		err = b.Put([]byte("l"), newBlock.currHash)
		if err != nil {
			log.Panic(err)
		}
		chain.tip = newBlock.currHash
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

func (chain *chain) iterator() *iterator {
	iter := &iterator{chain.tip, chain.db}
	
	return iter
}

func (iter *iterator) next() *block {
	var block *block

	err := iter.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		encBlock := b.Get(iter.currHash)
		block = bytesToBlock(encBlock)
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
	
	iter.currHash = block.prevHash
	
	return block
}

// Create a new blockchain
func newChain() *chain {
	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			fmt.Println("No existing blockhain found. Creating a new one...")
			first := newFirstBlock()
			b, err := tx.CreateBucket([]byte(bucket))
			if err != nil {
				log.Panic(err)
			}
			err = b.Put(first.currHash, first.blockToBytes())
			if err != nil {
				log.Panic(err)
			}
			err = b.Put([]byte("l"), first.currHash)
			if err != nil {
				log.Panic(err)
			}
			tip = first.currHash
		} else {
			tip = b.Get([]byte("l"))
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	chain := chain{tip, db}
	return &chain 
}
