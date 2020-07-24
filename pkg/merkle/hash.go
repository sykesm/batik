// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package merkle

import (
	"hash"
	"math/bits"
)

// The following constants are used to mitigate a scont preimage attack.
// See https://en.wikipedia.org/wiki/Merkle_tree#Second_preimage_attack for
// details.
const (
	LeafPrefix byte = 0 // LeafPrefix is prepended to data prior to hashing.
	NodePrefix byte = 1 // NodePrefix is prepended to intermediate node data prior to hashing.
)

// A Hash is responsble for creating new instances of hash.Hash. All
// crypto.Hash implementaions from the standard library satisfy this interface.
type Hash interface {
	New() hash.Hash
}

// NewHashFunc is an adapter to allow the use of a function as a Hash.
type NewHashFunc func() hash.Hash

// New calls the NewHashFunc function n.
func (n NewHashFunc) New() hash.Hash { return n() }

// Root uses the provided hash function to calculate the merkle root of the
// specified elements.
func Root(h Hash, leaves ...[]byte) []byte {
	switch len(leaves) {
	case 0:
		return h.New().Sum(nil)
	default:
		return hashTree(h, leaves...)
	}
}

type node struct {
	hash []byte
}

// type tree struct {
// 	hash  Hash
// 	size  uint64
// 	nodes [][]node
// }

// func (t *tree) Root() []byte {
// 	if len(t.nodes) == 0 {
// 		return t.hash.New().Sum(nil)
// 	}
// 	return t.nodes[len(t.nodes)-1][0].hash
// }

// func (t *tree) String() string {
// 	return fmt.Sprintf("%x", t.Root())
// }

// func (t *tree) Format(s fmt.State, verb rune) {
// 	switch verb {
// 	case 'v':
// 		if s.Flag('+') {
// 			for level := len(t.nodes) - 1; level >= 0; level-- {
// 				for n := range t.nodes[level] {
// 					if n != 0 {
// 						fmt.Fprintf(s, " ")
// 					}
// 					fmt.Fprintf(s, "%x", t.nodes[level][n].hash)
// 				}
// 				fmt.Fprintf(s, "\n")
// 			}
// 			return
// 		}
// 		fallthrough

// 	case 's', 'q':
// 		io.WriteString(s, t.String())
// 	}
// }

func hashTree(h Hash, leaves ...[]byte) []byte {
	levels := bits.Len64(uint64(len(leaves))) + 1

	// Allocate tree
	nodes := make([][]node, levels)
	for n := range nodes {
		nodes[n] = make([]node, 1<<(levels-n-1))
	}

	// Leaves go in level 0
	for i := 0; i < len(leaves); i++ {
		nodes[0][i].hash = hashLeaf(h, leaves[i])
	}

	// Intermediate nodes
	for level := 1; level < levels; level++ {
		for i := range nodes[level] {
			nodes[level][i].hash = hashNode(h, nodes[level-1][i*2].hash, nodes[level-1][i*2+1].hash)
		}
	}

	return nodes[levels-1][0].hash
}

func hashLeaf(hash Hash, leaf []byte) []byte {
	h := hash.New()
	h.Write([]byte{LeafPrefix})
	h.Write(leaf)
	return h.Sum(nil)
}

func hashNode(hash Hash, left, right []byte) []byte {
	if right == nil {
		return left
	}
	h := hash.New()
	h.Write([]byte{NodePrefix})
	h.Write(left)
	h.Write(right)
	return h.Sum(nil)
}
