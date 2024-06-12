package main

import (
	"flag"
	"fmt"
	"go-blockchain/blockchain"
)

func testMempool() *blockchain.Mempool {
	mempool := blockchain.NewMempool()

	numTransactions := 20

	for i := 0; i < numTransactions; i++ {
		privKey, pubKey := blockchain.GenerateKeyPair(2048)

		recipient := fmt.Sprintf("Recipient%d", i)
		amount := float64(i)

		tx := blockchain.CreateTransactionClient(blockchain.Address(pubKey), recipient, amount, privKey, pubKey)

		mempool.AddTransaction(tx)
	}

	return mempool
}

func main() {
	// Define command-line arguments
	nodeAddress := flag.String("node", "localhost:3000", "Address of the node")
	nodeName := flag.String("name", "Node", "Name of the node")
	peerAddress := flag.String("peer", "", "Address of the peer node")

	// Parse command-line arguments
	flag.Parse()

	// Create miner's keypair
	minerPrivKey, minerPubKey := blockchain.GenerateKeyPair(2048)
	minerAddress := blockchain.Address(minerPubKey)

	// Set initial reward for mining a block
	blockReward := 50.0

	node := blockchain.NewNode(*nodeAddress)
	node.Name = *nodeName

	if *peerAddress != "" {
		node.AddPeer(*peerAddress)
		fmt.Println("Added peer:", *peerAddress)
	}

	mempool := testMempool()
	node.Mempool = mempool

	node.PrivateKey = minerPrivKey
	node.PublicKey = minerPubKey

	go node.Start()
	fmt.Println(node.Name, "started at", *nodeAddress)

	go node.GenerateAndBroadcastBlock(blockReward, minerAddress)

	fmt.Println("Block reward:", blockReward, "Miner address:", minerAddress)
	select {} // Run forever
}
