// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package protomsg

import (
	"crypto"
	_ "crypto/sha256"
	_ "crypto/sha512"
	"fmt"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/sykesm/batik/pkg/pb/transaction"
)

var refs = []*transaction.StateReference{
	nil,
	{},
	{Txid: []byte{0, 1, 2, 3, 4, 5, 6, 7, 8}, OutputIndex: 9999},
}

func TestHashMessage(t *testing.T) {
	tests := []struct {
		hash     crypto.Hash
		expected string
		sr       *transaction.StateReference
	}{
		{crypto.SHA256, "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855", refs[0]},
		{crypto.SHA384, "38b060a751ac96384cd9327eb1b1e36a21fdb71114be07434c0cc7bf63f6e1da274edebfe76f65fbd51ad2f14898b95b", refs[0]},

		{crypto.SHA256, "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855", refs[1]},
		{crypto.SHA384, "38b060a751ac96384cd9327eb1b1e36a21fdb71114be07434c0cc7bf63f6e1da274edebfe76f65fbd51ad2f14898b95b", refs[1]},

		{crypto.SHA256, "3b0d642194674cc10f85ea9624740042a4889d4c62d5197a4d2a82ba544ae62e", refs[2]},
		{crypto.SHA384, "cd5ee194fad1032a2c13ddd1a1a85610b2fbc7aeafae67d3a49daa3646d3e2e98aab31bfdab435858853786fde2f309b", refs[2]},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			gt := NewGomegaWithT(t)

			hasher := &Hasher{Hash: tt.hash}
			actual, err := hasher.HashMessage(tt.sr)
			gt.Expect(err).NotTo(HaveOccurred())
			gt.Expect(fmt.Sprintf("%x", actual)).To(Equal(tt.expected))
		})
	}
}

func TestHashMessages(t *testing.T) {
	tests := []struct {
		expected []string
		msgs     []*transaction.StateReference
	}{
		{
			[]string{
				"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
				"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
				"3b0d642194674cc10f85ea9624740042a4889d4c62d5197a4d2a82ba544ae62e",
			},
			refs,
		},
		{[]string{}, nil},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			gt := NewGomegaWithT(t)
			hasher := &Hasher{Hash: crypto.SHA256}

			actual, err := hasher.HashMessages(tt.msgs)
			gt.Expect(err).NotTo(HaveOccurred())
			gt.Expect(actual).To(HaveLen(len(tt.expected)))
			for i := range actual {
				gt.Expect(fmt.Sprintf("%x", actual[i])).To(Equal(tt.expected[i]))
			}
		})
	}
}

func TestHashMessagesNonMessage(t *testing.T) {
	gt := NewGomegaWithT(t)
	hasher := &Hasher{Hash: crypto.SHA256}

	_, err := hasher.HashMessages("bob")
	gt.Expect(err).To(MatchError("protomsg: index 0 of type string is not a proto.Message"))
}
