// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package merkle

import (
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"math/bits"
)

// The following constants are used to mitigate a second preimage attack.
//
// See https://en.wikipedia.org/wiki/Merkle_tree#Second_preimage_attack for
// details.
const (
	LeafPrefix byte = 0 // LeafPrefix is prepended to data prior to hashing.
	NodePrefix byte = 1 // NodePrefix is prepended to intermediate node data prior to hashing.
)

// A Hasher is responsble for creating new instances of hash.Hash. All
// crypto.Hash implementaions from the standard library satisfy this
// interface.
type Hasher interface {
	New() hash.Hash
}

// NewHashFunc is an adapter to allow the use of a function as a Hasher.
type NewHashFunc func() hash.Hash

// New calls the NewHashFunc function n.
func (n NewHashFunc) New() hash.Hash { return n() }

// Root uses the provided hash function to calculate the merkle root of the
// specified elements.
func Root(h Hasher, leaves ...[]byte) []byte {
	return NewTree(h, leaves...).Root()
}

// A Tree represents a bnary Merkle Tree.
type Tree struct {
	hash  Hasher
	size  int
	nodes [][]node
}

type node struct {
	hash []byte
}

// NewTree constructs a binary Merkle tree from the provided data by using the
// provided Hasher. This implementation is based on RFC 6962 and uses leaf
// hashes directly as internal nodes when the tree is unbalanced.
//
// The tree does not support increment build processing so, unless diagnostic
// information is required, most consumers should use the Root function to
// calculate the root hash instead of creating a tree.
func NewTree(h Hasher, leaves ...[]byte) *Tree {
	levels := treeDepth(len(leaves))

	nodes := make([][]node, levels)
	for n := range nodes {
		// TODO(mjs): level 0 should be the same length as leaves
		nodes[n] = make([]node, 1<<(levels-n-1))
	}

	for i := 0; i < len(leaves); i++ {
		nodes[0][i].hash = hashLeaf(h, leaves[i])
	}

	for level := 1; level < levels; level++ {
		for i := range nodes[level] {
			// TODO(mjs): Need bounds check when level 0 is the same length as leaves
			nodes[level][i].hash = hashNode(h, nodes[level-1][i*2].hash, nodes[level-1][i*2+1].hash)
		}
	}

	return &Tree{hash: h, size: len(leaves), nodes: nodes}
}

func treeDepth(n int) int {
	switch {
	case n&(n-1) == 0:
		return bits.Len(uint(n))
	default:
		return bits.Len(uint(n)) + 1
	}
}

func hashLeaf(hash Hasher, leaf []byte) []byte {
	h := hash.New()
	h.Write([]byte{LeafPrefix})
	h.Write(leaf)
	return h.Sum(nil)
}

func hashNode(hash Hasher, left, right []byte) []byte {
	if right == nil {
		return left
	}
	h := hash.New()
	h.Write([]byte{NodePrefix})
	h.Write(left)
	h.Write(right)
	return h.Sum(nil)
}

// Root returns the Merkle root hash of the tree.
func (t *Tree) Root() []byte {
	if len(t.nodes) == 0 {
		return t.hash.New().Sum(nil)
	}
	return t.nodes[len(t.nodes)-1][0].hash
}

// String satisfies the fmt.Stringer interface and returns the hex encoded root
// hash of the tree.
func (t *Tree) String() string {
	return hex.EncodeToString(t.Root())
}

// Format satisfies the fmt.Formatter interace. The following format verbs are
// supported:
//
//     %s, %q, %v  print the hex encoded merkle hash root
//     %+v         dump the hashes of all elements in the tree
//
func (t *Tree) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			t.Dump(s)
			break
		}
		fallthrough

	case 's', 'q':
		io.WriteString(s, t.String())
	}
}

// Dump writes a text formatted representation of the tree to the provided
// writer. The tree is formatted horizontally in post-order traversal.
func (t *Tree) Dump(w io.Writer) {
	if t == nil || t.size == 0 {
		io.WriteString(w, "empty: "+hex.EncodeToString(t.Root())+"\n")
		return
	}

	io.WriteString(w, "root: "+hex.EncodeToString(t.Root())+"\n")
	if t.size > 1 {
		t.dumpLevel(w, len(t.nodes)-2, 0, "")
		t.dumpLevel(w, len(t.nodes)-2, 1, "")
	}
}

func (t *Tree) dumpLevel(w io.Writer, level, index int, prefix string) {
	glyphs := [][]string{
		{" ├─", " │ "}, // intermediate element
		{" └─", "   "}, // last element
	}
	p := glyphs[0]

	switch {
	case level == 0:
		if index == t.size-1 || index&1 == 1 {
			p = glyphs[1]
		}
		io.WriteString(w, prefix+p[0]+" leaf: "+hex.EncodeToString(t.nodes[level][index].hash)+"\n")

	case t.size <= index*(1<<level)+1<<(level-1):
		t.dumpLevel(w, level-1, 2*index, prefix)

	default:
		if index&1 == 1 || (index+1)*(1<<level) == t.size {
			p = glyphs[1]
		}

		io.WriteString(w, prefix+p[0]+" node: "+hex.EncodeToString(t.nodes[level][index].hash)+"\n")
		t.dumpLevel(w, level-1, 2*index, prefix+p[1])
		t.dumpLevel(w, level-1, 2*index+1, prefix+p[1])
	}
}
