package blockchain

import (
	"crypto/sha256"
)

// MerkleNode represents a node in the Merkle tree
type MerkleNode struct {
	Left  *MerkleNode
	Right *MerkleNode
	Data  []byte
}

// MerkleTree represents the entire Merkle tree
type MerkleTree struct {
	RootNode *MerkleNode
}

// NewMerkleNode creates a new Merkle node with given left and right nodes
func NewMerkleNode(left, right *MerkleNode, data []byte) *MerkleNode {
	node := MerkleNode{}

	if left == nil && right == nil {
		hash := sha256.Sum256(data)
		node.Data = hash[:]
	} else {
		prevHashes := append(left.Data, right.Data...)
		hash := sha256.Sum256(prevHashes)
		node.Data = hash[:]
	}

	node.Left = left
	node.Right = right

	return &node
}

// NewMerkleTree creates a new Merkle tree from a slice of transactions
func NewMerkleTree(data [][]byte) *MerkleTree {
	if len(data) == 0 {
		return &MerkleTree{RootNode: &MerkleNode{Data: []byte{}}}
	}

	var nodes []MerkleNode

	for _, datum := range data {
		node := NewMerkleNode(nil, nil, datum)
		nodes = append(nodes, *node)
	}

	for len(nodes) > 1 {
		var newLevel []MerkleNode

		for i := 0; i < len(nodes); i += 2 {
			var left, right *MerkleNode
			left = &nodes[i]

			if i+1 < len(nodes) {
				right = &nodes[i+1]
			} else {
				right = left
			}

			newNode := NewMerkleNode(left, right, nil)
			newLevel = append(newLevel, *newNode)
		}

		nodes = newLevel
	}

	tree := MerkleTree{&nodes[0]}
	return &tree
}
