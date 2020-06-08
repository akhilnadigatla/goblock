package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

const subsidy = 10

type TxInput struct {
	Txid   []byte
	Vout   int
	ScrSig string
}

type TxOutput struct {
	Value     int
	ScrPubKey string
}

type Txn struct {
	ID   []byte
	Vin  []TxInput
	Vout []TxOutput
}

func (txn *Txn) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	encoder := gob.NewEncoder(&encoded)
	err := encoder.Encode(txn)
	if err != nil {
		log.Panic(err)
	}
	hash = sha256.Sum256(encoded.Bytes())
	txn.ID = hash[:]
}

func (txn Txn) IsCoinbase() bool {
	return len(txn.Vin) == 1 && len(txn.Vin[0].Txid) == 0 && txn.Vin[0].Vout == -1
}

func (in *TxInput) CanUnlockOutputWith(data string) bool {
	return in.ScrSig == data
}

func (out *TxOutput) CanBeUnlockedWith(data string) bool {
	return out.ScrPubKey == data
}

func NewCoinbaseTxn(to, data string) *Txn {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}
	
	txIn := TxInput{[]byte{}, -1, data}
	txOut := TxOutput{subsidy, to}
	txn := Txn{nil, []TxInput{txIn}, []TxOutput{txOut}}	
	txn.SetID()

	return &txn
}

func NewUTOTxn(from, to string, amount int, chain *Chain) *Txn {
	var inputs []TxInput
	var outputs []TxOutput

	acc, vOutputs := chain.FindSTxnsO(from, amount)

	if acc < amount {
		log.Panic("Error: Insufficient funds.")
	}

	for txid, outs := range vOutputs {
		txID, err := hex.DecodeString(txid)
		if err != nil {
			log.Panic(err)
		}
		
		for _, out := range outs {
			input := TxInput{txID, out, from}
			inputs = append(inputs, input)
		}
	}

	outputs = append(outputs, TxOutput{amount, to})
	if acc > amount {
		outputs = append(outputs, TxOutput{acc - amount, from})
	}

	txn := Txn{nil, inputs, outputs}
	txn.SetID()

	return &txn	
}
