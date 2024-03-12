package blockchain

type Blockchain struct {
	Blocks []*Block
}

func NewBlockchain() *Blockchain {
	genesisBlock := NewBlock([]byte{}, []*Transaction{})
	return &Blockchain{Blocks: []*Block{genesisBlock}}
}

func (bc *Blockchain) AddBlock(transactions []*Transaction) {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := NewBlock(prevBlock.Hash, transactions)
	bc.Blocks = append(bc.Blocks, newBlock)
}
