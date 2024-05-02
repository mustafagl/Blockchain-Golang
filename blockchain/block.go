package blockchain

import (
	"crypto/sha256"
	"fmt"
	//"encoding/hex"
	 "unsafe"

)

type Block struct {
	PrevHash []byte
	Data     []byte
	Hash     []byte
	Nonce    int

}

func NewBlock(prevHash []byte, transactions []*Transaction) *Block {
	var err error
	var data []byte

	
	data, err = SerializeTransactions(transactions)
	//fmt.Println("DATA:", data)
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


	blockSize := int(unsafe.Sizeof(block)) +
		len(block.PrevHash) +
		len(block.Data) +
		len(block.Hash)

	fmt.Println("BLOK BOYUT:", blockSize)

	if blockSize > 1024*1024 {
		fmt.Println("Block struct size with data is larger than 1 MB")
	} else {
		fmt.Println("Block struct size with data is within 1 MB")
	}


	return block
}

func (b *Block) CalculateHash() []byte {
	data := append(b.PrevHash, b.Data...)
	hash := sha256.Sum256(data)
	return hash[:]
}