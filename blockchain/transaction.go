package blockchain

import (
	"bytes"
	"encoding/gob"
)

type Transaction struct {
	Sender    string
	Recipient string
	Amount    float64
	
}

func NewTransaction(sender, recipient string, amount float64) *Transaction {
	_, _ = GenerateKeyPair(2048)


	return &Transaction{
		Sender:    sender,
		Recipient: recipient,
		Amount:    amount,
	}
}

func SerializeTransactions(transactions []*Transaction) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(transactions)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}