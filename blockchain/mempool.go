package blockchain

import (
	"fmt"
	"sync"
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
		txData := []byte(fmt.Sprintf("%s%s%f", tx.Sender, tx.Recipient, tx.Amount))

		isValid := VerifySignature(tx.PublicKey, txData, tx.Signature)

		address := Address(tx.PublicKey)
		if address != tx.Sender {
			fmt.Printf("Sender is wrong")
			break
		}
		//fmt.Println("Address:")
		//fmt.Println(address)

		if !isValid {
			fmt.Printf("Signature is wrong")
			break
		}
		txSize, err := transactionSize(tx)
		if err != nil || !isValid {
			continue
		}
		//fmt.Printf("Total Size: %d \n", totalSize)
		if totalSize+txSize > 1024*1024 {
			break
		}
		selectedTransactions = append(selectedTransactions, tx)
		totalSize += txSize
	}

	return selectedTransactions
}

func (m *Mempool) Print() {
	m.mu.Lock()
	defer m.mu.Unlock()

	fmt.Println("Transactions in Mempool:")
	for _, tx := range m.Transactions {
		fmt.Printf("Sender: %s, Recipient: %s, Amount: %.2f\n", tx.Sender, tx.Recipient, tx.Amount)
	}
}
