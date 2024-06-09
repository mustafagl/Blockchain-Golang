package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"go-blockchain/blockchain"
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
	//byteSlice := make([]byte, 10)
	//Signature := blockchain.SignTransaction(privKey, byteSlice)

	mempool := blockchain.NewMempool()

	numTransactions := 20

	for i := 0; i < numTransactions; i++ {
		privKey, pubKey := blockchain.GenerateKeyPair(2048)

		// Create random sender and recipient names
		//sender := fmt.Sprintf("Sender%d", i)
		recipient := fmt.Sprintf("Recipient%d", i)
		// Create a random amount between 1 and 100
		amount := float64(i)
		// Create a new transaction
		tx := blockchain.CreateTransactionClient(blockchain.Address(pubKey), recipient, amount, privKey, pubKey)
		//tx := blockchain.NewTransaction(sender, recipient, amount, Signature, pubKey)

		// Add the transaction to the mempool
		mempool.AddTransaction(tx)
	}

	// Get all transactions from the mempool and print them

	return mempool
}
func main() {
	privKey, pubKey := blockchain.GenerateKeyPair(2048)
	byteSlice := make([]byte, 10)
	Signature := blockchain.SignTransaction(privKey, byteSlice)

	blockchain.CreateAndSaveNewWallet()

	mempool := testMempool()

	// Get all transactions from the mempool and print them
	transactions := mempool.GetTransactionsWithinLimit()
	for _, tx := range transactions {
		fmt.Printf("Transaction Sender: %s, Recipient: %s, Amount: %.2f\n",
			tx.Sender, tx.Recipient, tx.Amount)
	}

	bc := blockchain.NewBlockchain()
	//fmt.Printf("%+v\n", tx1)
	tx1 := blockchain.NewTransaction("Alice", "Bob", 100, Signature, pubKey)

	//tx2 := blockchain.NewTransaction("Bob", "Charlie", 50)

	//tx3 := blockchain.NewTransaction("Bob", "Alex", 50)

	coinbasetx := blockchain.NewTransaction("Reward System", "Miner0", 50, Signature, pubKey)

	prevBlock := bc.Blocks[len(bc.Blocks)-1]

	block1 := blockchain.NewBlock(prevBlock.Hash, transactions)
	block2 := blockchain.NewBlock(prevBlock.Hash, []*blockchain.Transaction{coinbasetx, tx1})

	bc.AddBlock(block1)
	bc.AddBlock(block2)

	//bc.AddBlock([]*blockchain.Transaction{coinbasetx, tx1, tx2, tx3})

}
