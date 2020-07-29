// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package merkle

import (
	"bytes"
	"crypto"
	"crypto/rand"
	_ "crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"testing"

	. "github.com/onsi/gomega"
	"gopkg.in/yaml.v2"
)

func TestRoot(t *testing.T) {
	tests := []struct {
		desc     string
		leaves   []string
		expected string
	}{
		{"nil", nil, digest(crypto.SHA256, nil)},
		{"empty", []string{}, digest(crypto.SHA256, nil)},
		{"empty leaf", []string{""}, digest(crypto.SHA256, []byte{0})},
		{"L1", []string{"L1"}, digest(crypto.SHA256, append([]byte{0}, []byte("L1")...))},
		{"L2", []string{"L2"}, digest(crypto.SHA256, append([]byte{0}, []byte("L2")...))},
		{"L3", []string{"L3"}, digest(crypto.SHA256, append([]byte{0}, []byte("L3")...))},
		{"L1,L2", []string{"L1", "L2"}, "0458611336c5dfbf775a6ca6196b215413be1d4e129a3c837633276e458da501"},
		{"L1,L2,L3", []string{"L1", "L2", "L3"}, "fb790cff1cc41df6229c8b4e399b57a4263a9532e9a5dfdff190337682ee836f"},
		{"L1,L2,L3,L4", []string{"L1", "L2", "L3", "L4"}, "41d0c7082e1794f1133cb7cebeaedb2818a93d7f4d697c4db5d2c97a37c536aa"},
		{"L1,L2,L3,L4,L5", []string{"L1", "L2", "L3", "L4", "L5"}, "8d5fe8e8394e4a793a9cee344558017546f5005608ad52db4e388c13dec299f9"},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			gt := NewGomegaWithT(t)
			leaves := make([][]byte, len(tt.leaves))
			for i := 0; i < len(tt.leaves); i++ {
				leaves[i] = []byte(tt.leaves[i])
			}
			actual := Root(crypto.SHA256, leaves...)
			gt.Expect(toHex(actual)).To(Equal(tt.expected), "got %x, want: %s", actual, tt.expected)

			tree := NewTree(crypto.SHA256, leaves...)
			gt.Expect(toHex(tree.Root())).To(Equal(tt.expected))
		})
	}
}

func TestHashFunc(t *testing.T) {
	gt := NewGomegaWithT(t)
	h := Root(NewHashFunc(crypto.SHA224.New))
	gt.Expect(h).To(Equal(crypto.SHA224.New().Sum(nil)))
}

func TestTreeString(t *testing.T) {
	gt := NewGomegaWithT(t)
	data := randomData(8, 0)
	for i := 0; i < len(data); i++ {
		tree := NewTree(crypto.SHA256, data...)
		gt.Expect(tree.String()).To(Equal(toHex(tree.Root())))
	}
}

func TestTreeDump(t *testing.T) {
	var tests []struct {
		Input    []string `yaml:"input,omitempty"`
		Expected string   `yaml:"expected,omitempty"`
	}

	gt := NewGomegaWithT(t)
	testdata, err := ioutil.ReadFile("testdata/dumptest.yaml")
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Expect(yaml.Unmarshal(testdata, &tests)).To(Succeed())

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.Input), func(t *testing.T) {
			gt := NewGomegaWithT(t)
			var input [][]byte
			for _, i := range tt.Input {
				input = append(input, []byte(i))
			}

			tree := NewTree(crypto.SHA256, input...)
			buf := bytes.NewBuffer(nil)
			tree.Dump(buf)
			gt.Expect(buf.String()).To(Equal(tt.Expected))
		})
	}
}

func TestFormat(t *testing.T) {
	for i := 0; i < 9; i++ {
		t.Run(fmt.Sprintf("%d elements", i), func(t *testing.T) {
			gt := NewGomegaWithT(t)
			data := randomData(i, 0)

			tree := NewTree(crypto.SHA256, data...)
			gt.Expect(fmt.Sprintf("%s", tree)).To(Equal(tree.String()))
			gt.Expect(fmt.Sprintf("%q", tree)).To(Equal(tree.String()))
			gt.Expect(fmt.Sprintf("%v", tree)).To(Equal(tree.String()))

			buf := bytes.NewBuffer(nil)
			tree.Dump(buf)
			gt.Expect(fmt.Sprintf("%+v", tree)).To(Equal(buf.String()))
		})
	}
}

func digest(hash crypto.Hash, b []byte) string {
	h := hash.New()
	h.Write(b)
	return toHex(h.Sum(nil))
}

func toHex(b []byte) string {
	return hex.EncodeToString(b)
}

// MTH is the algorithm described in RFC6962. We use it to ensure that our
// implmentation generates the same output.
func TestAgainstMTH(t *testing.T) {
	t.Run("NoNil", func(t *testing.T) {
		for i := 0; i <= 128; i++ {
			gt := NewGomegaWithT(t)
			data := randomData(i+1, 0)
			expected := MTH(data)
			actual := Root(crypto.SHA256, data...)
			gt.Expect(actual).To(Equal(expected), "got %x want %x for index %d", actual, expected, i)
		}
	})
	t.Run("WithNil", func(t *testing.T) {
		gt := NewGomegaWithT(t)
		data := randomData(8, 0)

		// Make one of the elements nil and ensure we get the same results
		for i := 0; i < len(data); i++ {
			d := data[i]
			data[i] = nil

			expected := MTH(data)
			actual := Root(crypto.SHA256, data...)
			gt.Expect(actual).To(Equal(expected), "got %x want %x for index %d", actual, expected, i)

			data[i] = d
		}
	})
	t.Run("WithEmpty", func(t *testing.T) {
		gt := NewGomegaWithT(t)
		data := randomData(8, 0)

		// Make one of the elements nil and ensure we get the same results
		for i := 0; i < len(data); i++ {
			d := data[i]
			data[i] = []byte{}

			expected := MTH(data)
			actual := Root(crypto.SHA256, data...)
			gt.Expect(actual).To(Equal(expected), "got %x want %x for index %d", actual, expected, i)

			data[i] = d
		}
	})
}

func randomData(n, l int) [][]byte {
	if l == 0 {
		l = int(randomBytes(1)[0]) // Up to 255 bytes
	}
	leaves := make([][]byte, n)
	for i := 0; i < n; i++ {
		leaves[i] = randomBytes(l)
	}
	return leaves
}

func randomBytes(n int) []byte {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return b
}

func BenchmarkMerkleRoot(b *testing.B) {
	leaves := randomData(2048, 256)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		hash := Root(crypto.SHA256, leaves...)
		if len(hash) == 0 {
			b.Fatalf("hash failed")
		}
	}
}

func BenchmarkMTH(b *testing.B) {
	leaves := randomData(2048, 256)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		hash := MTH(leaves)
		if len(hash) == 0 {
			b.Fatalf("hash failed")
		}
	}
}

// This is the algorithm defined by https://tools.ietf.org/html/rfc6962#section-2.1
//
// We compare our hash results to those produced by this algorithm.
func MTH(D [][]byte) []byte {
	h := crypto.SHA256.New()

	// Given an ordered list of n inputs, D[n] = {d(0), d(1), ..., d(n-1)},
	// the Merkle Tree Hash (MTH) is thus defined as follows:
	switch len(D) {
	case 0:
		// The hash of an empty list is the hash of an empty string:
		// MTH({}) = SHA-256().
		return h.Sum(nil)
	case 1:
		// The hash of a list with one entry (also known as a leaf hash) is:
		// MTH({d(0)}) = SHA-256(0x00 || d(0)).
		h.Write(append([]byte{0}, D[0]...))
		return h.Sum(nil)
	default:
		// For n > 1, let k be the largest power of two smaller than n
		// (ie. k < n <= 2k).
		//
		// The Merkle Tree Hash of an n-element list D[n] is then defined
		// recursively as
		//
		// MTH(D[n]) = SHA-256(0x01 || MTH(D[0:k]) || MTH(D[k:n])),
		//
		// where || is concatenation and D[k1:k2] denotes the list
		// {d(k1), d(k1+1), ..., d(k2-1)} of length (k2 - k1).
		n := len(D)
		k := largestPowerOfTwoLessThan(n)
		h.Write([]byte{1})
		h.Write(MTH(D[0:k]))
		h.Write(MTH(D[k:n]))
		return h.Sum(nil)
	}
}

func largestPowerOfTwoLessThan(n int) int {
	var k int
	for k = 1; k < n; {
		k *= 2
	}
	return k / 2
}

// This is the algorithm defined by https://tools.ietf.org/html/rfc6962#section-2.1.1
//
// We compare our audit paths to those produced by this algorithm.
func PATH(m int, D [][]byte) [][]byte {
	// Given an ordered list of n inputs to the tree, D[n] = {d(0), ...,
	// d(n-1)}, the Merkle audit path PATH(m, D[n]) for the (m+1)th input
	// d(m), 0 <= m < n, is defined as follows:
	switch {
	case m >= len(D):
		panic(fmt.Sprintf("%d is out of range", m))
	case m == 0 && len(D) == 1:
		// The path for the single leaf in a tree with a one-element input list
		// D[1] = {d(0)} is empty: PATH(0, {d(0)}) = {}
		return [][]byte{}
	default:
		// For n > 1, let k be the largest power of two smaller than n.  The
		// path for the (m+1)th element d(m) in a list of n > m elements is then
		// defined recursively as
		//
		// PATH(m, D[n]) = PATH(m, D[0:k]) : MTH(D[k:n]) for m < k; and
		//
		// PATH(m, D[n]) = PATH(m - k, D[k:n]) : MTH(D[0:k]) for m >= k,
		//
		// where : is concatenation of lists and D[k1:k2] denotes the length
		// (k2 - k1) list {d(k1), d(k1+1),..., d(k2-1)} as before.
		n := len(D)
		k := largestPowerOfTwoLessThan(n)
		if m < k {
			return append(PATH(m, D[0:k]), MTH(D[k:n]))
		}
		return append(PATH(m-k, D[k:n]), MTH(D[0:k]))
	}
}

//// This is the algorithm defined by https://tools.ietf.org/html/rfc6962#section-2.1.2
////
//// We compare our consistency proofs to those produced by this algorithm.
//func PROOF(m int, D [][]byte) [][]byte {
//	return SUBPROOF(m, D, true)
//}

//func SUBPROOF(m int, D [][]byte, b bool) [][]byte {
//	return nil
//}
