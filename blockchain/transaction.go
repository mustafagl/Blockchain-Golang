package blockchain

import (
	"bytes"
	"crypto/rsa"
	"encoding/gob"
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
func DeserializeTransactions(data []byte) ([]*Transaction, error) {
	var transactions []*Transaction
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	err := decoder.Decode(&transactions)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func SerializeTransaction(tx *Transaction) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(tx)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
func DeserializeTransaction(data []byte) (*Transaction, error) {
	var tx Transaction
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	err := decoder.Decode(&tx)
	if err != nil {
		return nil, err
	}
	return &tx, nil
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

func isValidTransactions(data []byte) bool {
	transactions, err := DeserializeTransactions(data)
	if err != nil {
		fmt.Println("Transaction deserialize edilirken bir hata oluştu:", err)
		return false
	}

	for _, tx := range transactions {
		// İşlem miktarı negatif olmamalı

		// İşlem gönderen ve alıcı aynı olmamalı
		if tx.Sender == tx.Recipient {
			fmt.Println("Geçersiz işlem: Gönderen ve alıcı aynı")
			return false
		}

		// İşlem imzası doğru olmalı
		txData := []byte(fmt.Sprintf("%s%s%f", tx.Sender, tx.Recipient, tx.Amount))

		if tx.Sender != "Reward System" {
			if !VerifySignature(tx.PublicKey, txData, tx.Signature) {
				fmt.Println("Geçersiz işlem imzası")
				return false
			}
		}

	}

	return true
}
