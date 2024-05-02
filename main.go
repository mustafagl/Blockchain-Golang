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

func testMempool() *blockchain.Mempool {
	// Create a new mempool
	mempool := blockchain.NewMempool()

	tx1 := blockchain.NewTransaction("Alice", "Bob", 100)

	tx2 := blockchain.NewTransaction("Charlie", "Alice", 200)

	mempool.AddTransaction(tx1)
	mempool.AddTransaction(tx2)

	// Get all transactions from the mempool and print them
	transactions := mempool.GetTransactionsWithinLimit()
	for _, tx := range transactions {
		fmt.Printf("Transaction Sender: %s, Recipient: %s, Amount: %.2f\n",
			tx.Sender, tx.Recipient, tx.Amount)
	}

	return mempool
}
func main() {

	blockchain.CreateAndSaveNewWallet()

	mempool := testMempool()

	// Get all transactions from the mempool and print them
	transactions := mempool.GetTransactions()
	for _, tx := range transactions {
		fmt.Printf("Transaction Sender: %s, Recipient: %s, Amount: %.2f\n",
			tx.Sender, tx.Recipient, tx.Amount)
	}


	bc := blockchain.NewBlockchain()
	//fmt.Printf("%+v\n", tx1)
	tx1 := blockchain.NewTransaction("Alice", "Bob", 100)

	//tx2 := blockchain.NewTransaction("Bob", "Charlie", 50)

	//tx3 := blockchain.NewTransaction("Bob", "Alex", 50)

	coinbasetx := blockchain.NewTransaction("Reward System", "Miner0", 50)

	prevBlock := bc.Blocks[len(bc.Blocks)-1]

	block1 := blockchain.NewBlock(prevBlock.Hash,[]*blockchain.Transaction{coinbasetx})
	block2 := blockchain.NewBlock(prevBlock.Hash,[]*blockchain.Transaction{coinbasetx,tx1})

	bc.AddBlock(block1)
	bc.AddBlock(block2)

	//bc.AddBlock([]*blockchain.Transaction{coinbasetx, tx1, tx2, tx3})

	for _, block := range bc.Blocks {
		fmt.Printf("Prev. hash: %x\n", block.PrevHash)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Printf("Data: %x\n", block.Data)
		fmt.Printf("Nonce: %x\n", block.Nonce)

		transactions, err := DeserializeTransactions(block.Data)
		if err != nil {
			fmt.Println("Error deserializing transactions:", err)
			continue
		}

		fmt.Println("Transactions:")
		for _, tx := range transactions {
			fmt.Printf("\tSender: %s, Recipient: %s, Amount: %f, Signature:%x, PublicKey:%x\n", tx.Sender, tx.Recipient, tx.Amount, tx.Signature, tx.PublicKey)
		}
		fmt.Println()
	}
}
