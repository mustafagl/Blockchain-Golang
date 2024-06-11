package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"time"
)

const targetBits = 17

type ProofOfWork struct {
	block      *Block
	blockchain *Blockchain
	target     *big.Int
}

func NewProofOfWork(b *Block, bc *Blockchain) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))
	pow := &ProofOfWork{block: b, blockchain: bc, target: target}
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

func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0
	reset := make(chan bool)

	// Goroutine to check for block changes
	go func() {
		initialBlock := pow.blockchain.LastBlock()
		initialBlockHash := sha256.Sum256(initialBlock.Data)
		for {
			time.Sleep(5 * time.Second) // Check every 5 seconds
			currentBlock := pow.blockchain.LastBlock()
			currentBlockHash := sha256.Sum256(currentBlock.Data)
			if currentBlockHash != initialBlockHash {
				reset <- true
				break
			}
		}
	}()

	for nonce < math.MaxInt64 {
		select {
		case <-reset:
			fmt.Println("Blockchain's last block changed, restarting the proof of work")
			nonce = 0
		default:
			data := pow.prepareData(nonce)
			hash = sha256.Sum256(data)
			hashInt.SetBytes(hash[:])

			if hashInt.Cmp(pow.target) == -1 {
				fmt.Printf("Found hash: %s\n", hex.EncodeToString(hash[:]))
				return nonce, hash[:]
			} else {
				nonce++
			}
		}
	}

	return nonce, hash[:]
}

func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	return hashInt.Cmp(pow.target) == -1
}
