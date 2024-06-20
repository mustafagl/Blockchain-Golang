package blockchain

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
)

type API struct {
	Node *Node
}

func NewAPI(node *Node) *API {
	fmt.Println("API Başlatıldı.")

	return &API{Node: node}
}

func (api *API) getTransactionsFromBlock(w http.ResponseWriter, r *http.Request) {
	// Parse request parameters
	params := r.URL.Query()
	blockHash := params.Get("blockHash")
	fmt.Println(blockHash)

	if blockHash == "" {
		http.Error(w, "Block hash is required", http.StatusBadRequest)
		return
	}

	byteHash, err := hex.DecodeString(blockHash)
	if err != nil {
		fmt.Println("Error decoding hex string:", err)
		return
	}

	var block *Block
	block = api.Node.Blockchain.GetBlockFromHash(byteHash)
	if block == nil {
		http.Error(w, "Invalid block hash", http.StatusBadRequest)
		return
	}

	fmt.Println("Data")
	fmt.Println(block.Data)
	fmt.Printf("Data (first 10 bytes): %x\n", block.Data[:])

	if string(block.Data) == "Genesis Block" {
		fmt.Println("Genesis Block detected, no transactions to deserialize.")
		// Kullanıcıya bilgi olarak bildir
		w.WriteHeader(http.StatusOK)                                                            // İsteğin başarılı olduğunu belirt
		w.Write([]byte("This block is a genesis block and does not contain any transactions.")) // Kullanıcıya mesajı gönder
		return
	}

	transactions, err := DeserializeTransactions(block.Data[:])
	if err != nil {
		fmt.Println("Error deserializing transactions:", err)
		http.Error(w, "Error processing block data", http.StatusInternalServerError)
		return
	}

	// Return the transactions as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
	fmt.Println(transactions)
	// Return the transaction as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}

func (api *API) getBlockchain(w http.ResponseWriter, r *http.Request) {
	blockchain := *api.Node.Blockchain
	blockchain.Print()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blockchain)
}

func (api *API) getLastBlock(w http.ResponseWriter, r *http.Request) {
	lastBlock := api.Node.Blockchain.LastBlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lastBlock)
}

func (api *API) getMempool(w http.ResponseWriter, r *http.Request) {
	mempool := api.Node.Mempool
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mempool)
}

func (api *API) Run(address string) {
	fmt.Println("API Run.")
	http.HandleFunc("/blockchain", api.getBlockchain)
	http.HandleFunc("/blockchain/last", api.getLastBlock)
	http.HandleFunc("/mempool", api.getMempool)
	http.HandleFunc("/transactions", api.getTransactionsFromBlock)

	fmt.Printf("API server started at %s\n", address)
	err := http.ListenAndServe(address, nil)
	if err != nil {
		fmt.Printf("Failed to start API server: %s\n", err)
	}
}
