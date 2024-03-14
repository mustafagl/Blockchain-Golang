package main

import (
	"fmt"
	"go-blockchain/blockchain"
	"bytes"
	"encoding/gob"
)

func DeserializeTransactions(data []byte) ([]*blockchain.Transaction, error) {
	var transactions []*blockchain.Transaction
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	err := decoder.Decode(&transactions)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func main() {
	bc := blockchain.NewBlockchain()

	tx1 := blockchain.NewTransaction("Alice", "Bob", 100)
	fmt.Printf("%+v\n", tx1)

	tx2 := blockchain.NewTransaction("Bob", "Charlie", 50)

	bc.AddBlock([]*blockchain.Transaction{tx1, tx2})

	for _, block := range bc.Blocks {
		fmt.Printf("Prev. hash: %x\n", block.PrevHash)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Printf("Data: %x\n", block.Data)

		transactions, err := DeserializeTransactions(block.Data)
		if err != nil {
			fmt.Println("Error deserializing transactions:", err)
			continue
		}

		fmt.Println("Transactions:")
		for _, tx := range transactions {
			fmt.Printf("\tSender: %s, Recipient: %s, Amount: %f\n", tx.Sender, tx.Recipient, tx.Amount)
		}
		fmt.Println()
	}
}
