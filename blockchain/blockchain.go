package blockchain

import (
	"crypto/sha256"
	"fmt"
	"time"
)

type Blockchain struct {
	Blocks []*Block
}

func NewBlockchain() *Blockchain {
	return &Blockchain{
		Blocks: []*Block{GenesisBlock()},
	}
}

func (bc *Blockchain) AddBlock(block *Block) {
	//fmt.Printf("Blok ekleniyor: %+v\n", block)
	bc.Blocks = append(bc.Blocks, block)
}

func (bc *Blockchain) LastBlock() *Block {
	if len(bc.Blocks) == 0 {
		return nil
	}
	lastBlock := bc.Blocks[len(bc.Blocks)-1]
	//fmt.Printf("Son blok alındı: %+v\n", lastBlock)
	return lastBlock
}

func (bc *Blockchain) Print() {
	for i, block := range bc.Blocks {
		fmt.Printf("========== Block %d ==========\n", i)
		block.Print()
	}
}

func GenesisBlock() *Block {
	// Genesis block data
	data := []byte("Genesis Block")
	prevHash := []byte{}
	timestamp := time.Now().Unix()

	// Merkle Root (for simplicity, just the hash of the data)
	merkleRoot := sha256.Sum256(data)

	block := &Block{
		PrevHash:   prevHash,
		Data:       data,
		MerkleRoot: merkleRoot[:],
		Timestamp:  timestamp,
	}

	// Calculate the hash of the block
	block.Hash = block.CalculateHash()

	return block
}
