// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

// Package merkletree defines a basic implementation for a merkle tree
// data structure.
package merkletree

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"hash"
)

var defaultHashStrategy = sha256.New

type MerkleTree struct {
	Root   *Node
	Leaves []*Node

	hashStrategy func() hash.Hash
}

type Node struct {
	Left  *Node
	Right *Node
	Hash  []byte
}

type Proof struct {
	MerkleRoot []byte
	ProofSet   [][]byte
	ProofIndex uint64
	NumLeaves  uint64
}

// MerkleRoot returns the hash of the merkle tree.
func (m *MerkleTree) MerkleRoot() ([]byte, error) {
	// If the tree is empty, return the hash of an empty string
	if m.Root == nil {
		h := m.hashStrategy()
		if _, err := h.Write([]byte{}); err != nil {
			return nil, err
		}

		return h.Sum(nil), nil
	}
	return m.Root.Hash, nil
}

// NewTree builds a merkle tree based off the input data representing the
// serialized leaves. It defaults to sha256 for the internal hashing algorithm.
func NewTree(serializedLeaves [][]byte) (*MerkleTree, error) {
	return NewTreeWithHashStrategy(serializedLeaves, defaultHashStrategy)
}

// GetMerkleTreeWithHashStrategy builds a merkle tree based off the input data
// representing the serialized leaves and a specified hash algorithm.
func NewTreeWithHashStrategy(serializedLeaves [][]byte, hashStrategy func() hash.Hash) (*MerkleTree, error) {
	t := &MerkleTree{
		hashStrategy: hashStrategy,
	}

	numLeaves := len(serializedLeaves)

	if numLeaves == 0 {
		return t, nil
	}

	// Pad with empty hash if not full
	for !isPow2(numLeaves) {
		serializedLeaves = append(serializedLeaves, []byte{})
		numLeaves++
	}

	root, leaves, err := buildMerkleTree(serializedLeaves, t)
	if err != nil {
		return nil, err
	}

	t.Root = root
	t.Leaves = leaves

	return t, nil
}

func isPow2(n int) bool {
	return n&(n-1) == 0
}

func buildMerkleTree(serializedLeaves [][]byte, t *MerkleTree) (*Node, []*Node, error) {
	var leaves []*Node

	for _, l := range serializedLeaves {
		h := t.hashStrategy()
		if _, err := h.Write(l); err != nil {
			return nil, nil, err
		}

		leaves = append(leaves, &Node{
			Hash: leafHash(h.Sum(nil)),
		})
	}

	root, err := buildIntermediate(leaves, t)
	if err != nil {
		return nil, nil, err
	}

	return root, leaves, nil
}

func buildIntermediate(nodes []*Node, t *MerkleTree) (*Node, error) {
	numNodes := len(nodes)

	// Return as root node if tree has one input
	if numNodes == 1 {
		return &Node{
			Hash: nodes[0].Hash,
		}, nil
	}

	// If not the root node, numNodes should never be odd due to prior padding
	if numNodes%2 != 0 {
		return nil, errors.New("number of nodes is not even")
	}

	var combinedHashNodes []*Node

	for i := 0; i < numNodes; i += 2 {
		left := nodes[i]
		right := nodes[i+1]

		h := t.hashStrategy()
		newHash := append(left.Hash, right.Hash...)
		if _, err := h.Write(newHash); err != nil {
			return nil, err
		}

		combined := &Node{
			Left:  left,
			Right: right,
			Hash:  intermediateHash(h.Sum(nil)),
		}

		// When down to last 2 nodes, return the combined node
		if numNodes == 2 {
			return combined, nil
		}

		combinedHashNodes = append(combinedHashNodes, combined)
	}

	return buildIntermediate(combinedHashNodes, t)
}

// leafHash prepends a 0x00 byte to the hashed data to protect against second
// preimage attacks against merkle trees.
// See: https://en.wikipedia.org/wiki/Merkle_tree#Second_preimage_attack
func leafHash(hash []byte) []byte {
	return append([]byte{0x00}, hash...)
}

// intermediateHash prepends a 0x01 byte to the hashed data to protect against second
// preimage attacks against merkle trees.
// See: https://en.wikipedia.org/wiki/Merkle_tree#Second_preimage_attack
func intermediateHash(hash []byte) []byte {
	return append([]byte{0x01}, hash...)
}

// TODO
func (m *MerkleTree) MerkleProof(hash []byte) (*Proof, error) {
	var proofIndex uint64

	dataFoundInTree := false
	for i, l := range m.Leaves {
		if bytes.Equal(l.Hash, hash) {
			dataFoundInTree = true
			proofIndex = uint64(i)
		}
	}

	if !dataFoundInTree {
		return nil, errors.New("hash not found in tree")
	}

	return &Proof{
		MerkleRoot: m.Root.Hash,
		ProofSet:   [][]byte{},
		ProofIndex: proofIndex,
		NumLeaves:  uint64(len(m.Leaves)),
	}, nil
}

// TODO
func (p *Proof) VerifyProof(hashStrategy func() hash.Hash) bool {
	if p.MerkleRoot == nil {
		return false
	}
	if p.ProofIndex >= p.NumLeaves {
		return false
	}

	return true
}
