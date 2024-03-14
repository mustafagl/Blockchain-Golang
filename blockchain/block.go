package blockchain

import (
	"crypto/sha256"
	"fmt"
	//"encoding/hex"
)

type Block struct {
	PrevHash []byte
	Data     []byte
	Hash     []byte
}

func NewBlock(prevHash []byte, transactions []*Transaction) *Block {
	var err error
	var data []byte
	data, err = SerializeTransactions(transactions)
	fmt.Println("DATA:", data)
	if err != nil {
		// Handle the error, e.g., log it or return a block with an error field
		// For simplicity, we'll just log it and return nil in this example
		fmt.Println("Error serializing transactions:", err)
		return nil
	}

	block := &Block{
		PrevHash: prevHash,
		Data: data,
	}
	block.Hash = block.CalculateHash()

	fmt.Println("HASH:", block.Hash)

	return block
}

func (b *Block) CalculateHash() []byte {
	data := append(b.PrevHash, b.Data...)
	hash := sha256.Sum256(data)
	return hash[:]
}