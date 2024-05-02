package blockchain

import (
	"sync"
	"fmt"
)

// Mempool struct represents a mempool
type Mempool struct {
	Transactions []*Transaction
	mu           sync.Mutex
}

// NewMempool creates a new mempool
func NewMempool() *Mempool {
	return &Mempool{
		Transactions: []*Transaction{},
	}
}

// AddTransaction adds a transaction to the mempool
func (m *Mempool) AddTransaction(tx *Transaction) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Transactions = append(m.Transactions, tx)
}

// GetTransactions returns all transactions in the mempool
func (m *Mempool) GetTransactions() []*Transaction {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.Transactions
}

func (m *Mempool) GetTransactionsWithinLimit() []*Transaction {
	m.mu.Lock()
	defer m.mu.Unlock()

	var selectedTransactions []*Transaction
	var totalSize int

	for _, tx := range m.Transactions {
		txSize, err := transactionSize(tx)
		if err != nil {
			continue
		}		
		fmt.Printf("Total Size: %d \n", totalSize)
		if totalSize+txSize > 1024*1024 {
			break
		}
		selectedTransactions = append(selectedTransactions, tx)
		totalSize += txSize
	}

	return selectedTransactions
}

