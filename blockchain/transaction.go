package blockchain

import (
	"bytes"
	"encoding/gob"
	"crypto/rsa"
	"fmt"
)

type Transaction struct {
	Sender    string
	Recipient string
	Amount    float64
	PublicKey *rsa.PublicKey
	Signature []byte
}

func NewTransaction(sender, recipient string, amount float64) *Transaction {
	privKey, pubKey := GenerateKeyPair(2048)
	tx := &Transaction{
		Sender:    sender,
		Recipient: recipient,
		Amount:    amount,
	}	

	txData := []byte(fmt.Sprintf("%s%s%f", tx.Sender, tx.Recipient, tx.Amount))

	tx.Signature = SignTransaction(privKey, txData)
	tx.PublicKey = pubKey

	return tx
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