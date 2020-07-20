// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package merkletree

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"hash"
	"testing"

	. "github.com/onsi/gomega"
)

func TestNewTree(t *testing.T) {
	tests := []struct {
		testName     string
		inputs       [][]byte
		expectedHash string
	}{
		{
			testName: "0 levels",
			inputs: [][]byte{
				[]byte("inputA"),
			},
			/*
			* Leaves/Root (level 0):
			* A/root:	0x00+sha256("inputA") -> 003df26252b79cbbec76756cefc5bb3125df63180ce737c06f92bd4afbcc6d34f3
			 */
			expectedHash: "003df26252b79cbbec76756cefc5bb3125df63180ce737c06f92bd4afbcc6d34f3",
		},
		{
			testName: "1 level",
			inputs: [][]byte{
				[]byte("inputA"),
				[]byte("inputB"),
			},
			/*
			* Leaves (level 0):
			* A:				0x00+sha256("inputA") -> 003df26252b79cbbec76756cefc5bb3125df63180ce737c06f92bd4afbcc6d34f3
			* B:				0x00+sha256("inputB") -> 00b7686ae41a31ccc3db28939d8702fb6c2f31893852e3d8dfe6096d43cdc0f7c6
			*
			* Level 1/Root Level:
			* AB/root:	0x01+sha256(A+B) -> 014b2aab6bded4698071f17454b66940a19e67a9f3a70c3a91a811410b0c617c31
			 */
			expectedHash: "014b2aab6bded4698071f17454b66940a19e67a9f3a70c3a91a811410b0c617c31",
		},
		{
			testName: "2 levels - with padding",
			inputs: [][]byte{
				[]byte("inputA"),
				[]byte("inputB"),
				[]byte("inputC"),
			},
			/*
			* Leaves (level 0):
			* A:		0x00+sha256("inputA") -> 003df26252b79cbbec76756cefc5bb3125df63180ce737c06f92bd4afbcc6d34f3
			* B:		0x00+sha256("inputB") -> 00b7686ae41a31ccc3db28939d8702fb6c2f31893852e3d8dfe6096d43cdc0f7c6
			* C:		0x00+sha256("inputC") -> 000ec29b39efcb9b6c5ee1a0634246a66a2b1bf62a2f26611a319cb5c7583c018d
			* pad:	0x00+sha256() -> 00e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
			*
			* Level 1:
			* AB:		0x01+sha256(A+B) -> 014b2aab6bded4698071f17454b66940a19e67a9f3a70c3a91a811410b0c617c31
			* Cpad: 0x01+sha256(C+pad) -> 01c74bafa689760858a90edcf26cff27d017a7428a22847049b9967408d64e05ce
			*
			* Root Level (level 2):
			* root: 0x01+sha256(AB+Cpad) -> 01f068c720286abd5c811e93d4eb4e99e897b369be1596e4775b26c3ed646de02e
			 */
			expectedHash: "01f068c720286abd5c811e93d4eb4e99e897b369be1596e4775b26c3ed646de02e",
		},
		{
			testName: "2 levels - no padding",
			inputs: [][]byte{
				[]byte("inputA"),
				[]byte("inputB"),
				[]byte("inputC"),
				[]byte("inputD"),
			},
			/*
			* Leaves (level 0):
			* A:		0x00+sha256("inputA") -> 003df26252b79cbbec76756cefc5bb3125df63180ce737c06f92bd4afbcc6d34f3
			* B:		0x00+sha256("inputB") -> 00b7686ae41a31ccc3db28939d8702fb6c2f31893852e3d8dfe6096d43cdc0f7c6
			* C:		0x00+sha256("inputC") -> 000ec29b39efcb9b6c5ee1a0634246a66a2b1bf62a2f26611a319cb5c7583c018d
			* D:		0x00+sha256("inputD") -> 0003f95d1378ae9d1e00d59417ada7600dfc065410a9e3404fd8db44b1e799e3bd
			*
			* Level 1:
			* AB:		0x01+sha256(A+B) -> 014b2aab6bded4698071f17454b66940a19e67a9f3a70c3a91a811410b0c617c31
			* CD:		0x01+sha256(C+D) -> 012b7d6aaa2168ad4db8f431ddcdd0b7f1aadbda6b9c1f404d5755176e4339e9cf
			*
			* Root Level (level 2):
			* root: 0x01+sha256(AB+CD) -> 01133fb2026e2d0c7405c8b4edbbbbc2fd9374cb5aa083adf3f17feab06e81c7b0
			 */
			expectedHash: "01133fb2026e2d0c7405c8b4edbbbbc2fd9374cb5aa083adf3f17feab06e81c7b0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			gt := NewGomegaWithT(t)

			tree, err := NewTree(tt.inputs)

			gt.Expect(err).NotTo(HaveOccurred())

			expectedHash, err := hex.DecodeString(tt.expectedHash)
			gt.Expect(err).NotTo(HaveOccurred())
			gt.Expect(tree.MerkleRoot()).To(Equal(expectedHash))
		})
	}
}

func TestNewTree_Failure(t *testing.T) {
	gt := NewGomegaWithT(t)

	tree, err := NewTree([][]byte{})
	gt.Expect(err).To(MatchError("empty leaf hashes"))
	gt.Expect(tree).To(BeNil())
}

func TestNewTreeWithHashingStrategy(t *testing.T) {
	tests := []struct {
		testName        string
		hashingStrategy func() hash.Hash
		inputs          [][]byte
		expectedHash    string
	}{
		{
			testName:        "MD5",
			hashingStrategy: md5.New,
			inputs: [][]byte{
				[]byte("inputA"),
				[]byte("inputB"),
				[]byte("inputC"),
			},
			/*
			* Leaves (level 0):
			* A:		0x00+md5("inputA") -> 00c481761c4701a77c08bcfcf074a506eb
			* B:		0x00+md5("inputB") -> 00e634df8e289d7f0b748e16b71d502373
			* C:		0x00+md5("inputC") -> 004e269439d9ffe1e516f0dd2e8d7f6bbc
			* pad:	0x00+md5() -> 00d41d8cd98f00b204e9800998ecf8427e
			*
			* Level 1:
			* AB:		0x01+md5(A+B) -> 0111a7dcd31e9e5fc1a221b060ad22c84b
			* Cpad: 0x01+md5(C+pad) -> 019f4e040136fb6be14a4ade6ca57746c4
			*
			* Root Level (level 2):
			* root: 0x01+md5(AB+Cpad) -> 010f28e622d7bd149088ef921fcae88c2a
			 */
			expectedHash: "010f28e622d7bd149088ef921fcae88c2a",
		},
		{
			testName:        "SHA1",
			hashingStrategy: sha1.New,
			inputs: [][]byte{
				[]byte("inputA"),
				[]byte("inputB"),
				[]byte("inputC"),
			},
			/*
			* Leaves (level 0):
			* A:		0x00+sha1("inputA") -> 0060c07de44bf3680efbf00c0d92a2fc9ba7138c1a
			* B:		0x00+sha1("inputB") -> 00460dbf3f2b631aaaacded7dd0f279985c46b2d42
			* C:		0x00+sha1("inputC") -> 00dbdb58c0193d5fdb9c06e9d487247560538c43d9
			* pad:	0x00+sha1() -> 00da39a3ee5e6b4b0d3255bfef95601890afd80709
			*
			* Level 1:
			* AB:		0x01+sha1(A+B) -> 0137c4ce778516b16b9b95405f044a8f4ff5b3d8bf
			* Cpad: 0x01+sha1(C+pad) -> 01e696c54e163029d818b33c6a5c94d5f0478388bd
			*
			* Root Level (level 2):
			* root: 0x01+sha1(AB+Cpad) -> 019a2130adb8408cd7ed01c08561d55c0309b3f8d8
			 */
			expectedHash: "019a2130adb8408cd7ed01c08561d55c0309b3f8d8",
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			gt := NewGomegaWithT(t)

			tree, err := NewTreeWithHashStrategy(tt.inputs, tt.hashingStrategy)

			gt.Expect(err).NotTo(HaveOccurred())

			expectedHash, err := hex.DecodeString(tt.expectedHash)
			gt.Expect(err).NotTo(HaveOccurred())
			gt.Expect(tree.MerkleRoot()).To(Equal(expectedHash))
		})
	}
}

func TestSecondPreImageAttack(t *testing.T) {
	gt := NewGomegaWithT(t)

	/*
	* Leaves (level 0):
	* A:		0x00+sha256("inputA") -> 003df26252b79cbbec76756cefc5bb3125df63180ce737c06f92bd4afbcc6d34f3
	* B:		0x00+sha256("inputB") -> 00b7686ae41a31ccc3db28939d8702fb6c2f31893852e3d8dfe6096d43cdc0f7c6
	* C:		0x00+sha256("inputC") -> 000ec29b39efcb9b6c5ee1a0634246a66a2b1bf62a2f26611a319cb5c7583c018d
	* pad:	0x00+sha256() -> 00e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
	*
	* Level 1:
	* AB:		0x01+sha256(A+B) -> 014b2aab6bded4698071f17454b66940a19e67a9f3a70c3a91a811410b0c617c31
	* Cpad: 0x01+sha256(C+pad) -> 01c74bafa689760858a90edcf26cff27d017a7428a22847049b9967408d64e05ce
	*
	* Root Level (level 2):
	* root: 0x01+sha256(AB+Cpad) -> 01f068c720286abd5c811e93d4eb4e99e897b369be1596e4775b26c3ed646de02e
	 */

	tree1, err := NewTree([][]byte{
		[]byte("inputA"),
		[]byte("inputB"),
		[]byte("inputC"),
	})
	gt.Expect(err).NotTo(HaveOccurred())

	ab, err := hex.DecodeString("014b2aab6bded4698071f17454b66940a19e67a9f3a70c3a91a811410b0c617c31")
	gt.Expect(err).NotTo(HaveOccurred())

	cpad, err := hex.DecodeString("01c74bafa689760858a90edcf26cff27d017a7428a22847049b9967408d64e05ce")
	gt.Expect(err).NotTo(HaveOccurred())

	tree2, err := NewTree([][]byte{
		ab, cpad,
	})
	gt.Expect(err).NotTo(HaveOccurred())

	// A tree formed from the combined hashes of the original leaves
	// should not produce the same tree
	gt.Expect(tree1.MerkleRoot()).NotTo(Equal(tree2.MerkleRoot()))
}
