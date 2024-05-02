package blockchain

type Blockchain struct {
	Blocks []*Block
}

func NewBlockchain() *Blockchain {
	genesisBlock := NewBlock([]byte{}, []*Transaction{})
	return &Blockchain{Blocks: []*Block{genesisBlock}}
}

func (bc *Blockchain) AddBlock(nb *Block) {

	bc.Blocks = append(bc.Blocks, nb)
}
