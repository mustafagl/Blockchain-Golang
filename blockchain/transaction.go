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

func NewTransaction(sender, recipient string, amount float64, signature []byte, pubKey *rsa.PublicKey) *Transaction {
	tx := &Transaction{
		Sender:    sender,
		Recipient: recipient,
		Amount:    amount,
		PublicKey: pubKey,
		Signature: signature,
	}	

	return tx
}

func CreateTransactionClient(sender, recipient string, amount float64, privKey *rsa.PrivateKey, pubKey *rsa.PublicKey) *Transaction {
	tx := &Transaction{
		Sender:    sender,
		Recipient: recipient,
		Amount:    amount,
		PublicKey: pubKey,
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

func transactionSize(tx *Transaction) (int, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(tx)
	if err != nil {
		return 0, err
	}
	return buffer.Len(), nil
}