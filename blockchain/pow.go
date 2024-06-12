package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"math/rand"
	"time"
)

const targetBits = 21

type ProofOfWork struct {
	block  *Block
	node   *Node
	target *big.Int
}

func NewProofOfWork(b *Block, node *Node) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))
	pow := &ProofOfWork{block: b, node: node, target: target}
	return pow
}

func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := append(pow.block.PrevHash, pow.block.Data...)
	data = append(data, pow.block.MerkleRoot...)
	data = append(data, []byte(fmt.Sprintf("%d", pow.block.Timestamp))...)
	data = append(data, []byte(fmt.Sprintf("%x", targetBits))...)
	data = append(data, []byte(fmt.Sprintf("%d", nonce))...)
	return data
}

func (pow *ProofOfWork) Run() (int, []byte, bool) {
	var hashInt big.Int
	var hash [32]byte
	reset := make(chan bool)

	// Goroutine to check for block changes
	go func() {
		initialBlock := pow.node.Blockchain.LastBlock()
		initialBlockHash := sha256.Sum256(initialBlock.Data)
		for {

			currentBlock := pow.node.Blockchain.LastBlock()
			currentBlockHash := sha256.Sum256(currentBlock.Data)
			if currentBlockHash != initialBlockHash {
				reset <- true
				break
			}
		}
	}()

	for {
		select {
		case <-reset:
			fmt.Println("Blockchain's last block changed, stopping the proof of work")
			return 0, nil, true // Madencilik işlemini durdurmak için true değerini döndür
		default:
			rnd := rand.New(rand.NewSource(time.Now().UnixNano())) // Yeni bir rastgele sayı üretici
			nonce := rnd.Int()                                     // Rastgele bir nonce değeri üret
			data := pow.prepareData(nonce)
			hash = sha256.Sum256(data)
			hashInt.SetBytes(hash[:])

			if hashInt.Cmp(pow.target) == -1 {
				fmt.Printf("Found hash: %s\n", hex.EncodeToString(hash[:]))
				return nonce, hash[:], false // Madencilik işlemi başarıyla tamamlandı
			}
		}
	}
}

func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	return hashInt.Cmp(pow.target) == -1
}
