package blockchain

import (
	"encoding/gob"
	"fmt"
	"net"
	"sync"
)

func init() {
	gob.Register(Block{})
	gob.Register(Transaction{})
}

type Node struct {
	Address    string
	Blockchain *Blockchain
	Mempool    *Mempool
	Peers      []string
	mutex      sync.Mutex
	Name       string // New field
}

func NewNode(address string) *Node {
	return &Node{
		Address:    address,
		Blockchain: NewBlockchain(),
		Mempool:    NewMempool(),
	}
}

func (node *Node) AddPeer(peer string) {
	node.mutex.Lock()
	defer node.mutex.Unlock()
	node.Peers = append(node.Peers, peer)
}

func (node *Node) HandleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Yeni bağlantı işleniyor") // Debug ifadesi

	decoder := gob.NewDecoder(conn)
	for {
		var request map[string]interface{}
		err := decoder.Decode(&request)
		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println("İstemci bağlantıyı kapattı")
				break
			}
			fmt.Println("İstek çözümlenemedi:", err)
			continue
		}

		fmt.Println("Gelen istek:", request) // Debug ifadesi

		switch request["type"] {
		case "block":
			fmt.Println("Yeni blok alındı")
			var block Block
			err := decoder.Decode(&block)
			if err != nil {
				fmt.Println("Blok çözümlenemedi:", err)
				continue
			}
			if isValidBlock(block) {
				node.Blockchain.AddBlock(&block)
				node.Blockchain.Print()
			}
		case "transaction":
			var tx Transaction
			err := decoder.Decode(&tx)
			if err != nil {
				fmt.Println("İşlem çözümlenemedi:", err)
				continue
			}
			node.Mempool.AddTransaction(&tx)
			node.Mempool.Print()
		case "sync":
			encoder := gob.NewEncoder(conn)
			err := encoder.Encode(node.Blockchain)
			if err != nil {
				fmt.Println("Blockchain gönderilemedi:", err)
			}
		default:
			fmt.Println("Bilinmeyen istek türü:", request["type"])
		}
	}
}

func isValidBlock(block Block) bool {
	// Implement your validation logic here
	// For example, checking the block's hash, transactions, etc.
	return true
}

func (node *Node) Start() {
	ln, err := net.Listen("tcp", node.Address)
	if err != nil {
		fmt.Println("Failed to start server:", err)
		return
	}
	defer ln.Close()

	node.SyncBlockchain()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Failed to accept connection:", err)
			continue
		}
		go node.HandleConnection(conn)
	}
}

func (node *Node) BroadcastBlock(block *Block) {
	for _, peer := range node.Peers {
		go func(peer string) {
			fmt.Println("Peer'e bağlanılıyor:", peer) // Debug ifadesi
			conn, err := net.Dial("tcp", peer)
			if err != nil {
				fmt.Println("Peer'e bağlanılamadı:", err)
				return
			}
			defer conn.Close()

			encoder := gob.NewEncoder(conn)
			// Önce mesaj türünü kodlayın
			err = encoder.Encode(map[string]interface{}{
				"type": "block",
			})
			if err != nil {
				fmt.Println("Blok türü gönderilemedi:", err)
				return
			}

			// Blok verisini kodlayın
			err = encoder.Encode(block)
			if err != nil {
				fmt.Println("Blok verisi gönderilemedi:", err)
				return
			}

			fmt.Println("Blok peer'e gönderildi:", peer) // Debug ifadesi
		}(peer)
	}
}

func (node *Node) SyncBlockchain() {
	for _, peer := range node.Peers {
		conn, err := net.Dial("tcp", peer)
		if err != nil {
			fmt.Println("Failed to connect to peer:", err)
			continue
		}
		defer conn.Close()

		encoder := gob.NewEncoder(conn)
		err = encoder.Encode(map[string]interface{}{
			"type": "sync",
		})
		if err != nil {
			fmt.Println("Failed to send sync request:", err)
			continue
		}

		var response Blockchain
		decoder := gob.NewDecoder(conn)
		err = decoder.Decode(&response)
		if err != nil {
			fmt.Println("Failed to receive blockchain:", err)
			continue
		}

		// Replace the local blockchain with the received blockchain if it's longer
		if len(response.Blocks) > len(node.Blockchain.Blocks) {
			node.Blockchain = &response
		}
	}

	fmt.Println("Blockchain synchronized")
}

func (node *Node) GenerateAndBroadcastBlock(reward float64, minerAddress string) {
	for {
		prevBlock := node.Blockchain.Blocks[len(node.Blockchain.Blocks)-1]
		transactions := node.Mempool.GetTransactionsWithinLimit()
		if len(transactions) == 0 {
			continue
		}
		newBlock := NewBlock(prevBlock.Hash, transactions, reward, node.Name, node.Blockchain)
		newBlock.MinerName = node.Name // Add node name to the block
		node.Blockchain.AddBlock(newBlock)
		node.BroadcastBlock(newBlock)
		node.Blockchain.Print()
	}
}

func (node *Node) BroadcastTransaction(tx *Transaction) {
	for _, peer := range node.Peers {
		go func(peer string) {
			conn, err := net.Dial("tcp", peer)
			if err != nil {
				fmt.Println("Failed to connect to peer:", err)
				return
			}
			defer conn.Close()

			encoder := gob.NewEncoder(conn)
			err = encoder.Encode(map[string]interface{}{
				"type":        "transaction",
				"transaction": tx,
			})
			if err != nil {
				fmt.Println("Failed to send transaction:", err)
			}
		}(peer)
	}
}
