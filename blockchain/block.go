package blockchain

import (
	"crypto/sha256"
	"fmt"
	"time"
	"unsafe"
)

type Block struct {
	PrevHash   []byte
	Data       []byte
	Hash       []byte
	Nonce      int
	MerkleRoot []byte
	Timestamp  int64
	MinerName  string // Madenci adı
}

func NewBlock(prevHash []byte, transactions []*Transaction, reward float64, miner string, bc *Blockchain) *Block {
	var err error
	var data []byte
	startTime := time.Now()

	// Add reward transaction
	rewardTx := &Transaction{
		Sender:    "Reward System",
		Recipient: miner,
		Amount:    reward,
	}
	transactions = append(transactions, rewardTx)

	data, err = SerializeTransactions(transactions)
	if err != nil {
		fmt.Println("Error serializing transactions:", err)
		return nil
	}

	// Create Merkle Tree
	var txData [][]byte
	for _, tx := range transactions {
		txBytes, err := SerializeTransaction(tx)
		if err != nil {
			fmt.Println("Error serializing transaction:", err)
			return nil
		}
		txData = append(txData, txBytes)
	}

	merkleTree := NewMerkleTree(txData)
	merkleRoot := merkleTree.RootNode.Data

	block := &Block{
		PrevHash:   prevHash,
		Data:       data,
		MerkleRoot: merkleRoot,
		Timestamp:  time.Now().Unix(),
		MinerName:  miner,
	}

	pow := NewProofOfWork(block, bc)
	nonce, hash := pow.Run()
	block.Nonce = nonce
	block.Hash = hash

	endTime := time.Now()

	// Süreyi hesapla
	duration := endTime.Sub(startTime)
	fmt.Printf("Duration pow: %v\n", duration)

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
	data = append(data, b.MerkleRoot...)
	hash := sha256.Sum256(data)
	return hash[:]
}

func (b *Block) Print() {
	fmt.Printf("Block - Hash: %x\n", b.Hash)
	fmt.Printf("Previous Hash: %x\n", b.PrevHash)
	fmt.Printf("Merkle Root: %x\n", b.MerkleRoot)
	fmt.Printf("Data (first 10 bytes): %x\n", b.Data[:10])
	fmt.Printf("Nonce: %d\n", b.Nonce)
	fmt.Printf("Miner: %s\n", b.MinerName)

	fmt.Println()
}
