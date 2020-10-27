// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package transaction

import (
	"crypto"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/sykesm/batik/pkg/merkle"
	txv1 "github.com/sykesm/batik/pkg/pb/tx/v1"
	"github.com/sykesm/batik/pkg/protomsg"
	"github.com/sykesm/batik/pkg/transaction/internal/testprotos/mutated"
)

func TestSalt(t *testing.T) {
	tests := []struct {
		salt []byte
		fn   uint32
		idx  uint32
	}{
		{nil, 0, 0},
		{[]byte("NaCl"), 0, 0},
		{nil, 1, 0},
		{nil, 0, 1},
		{nil, 1, 2},
		{[]byte("NaCl"), 0xffffffff, 2},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			gt := NewGomegaWithT(t)

			hash := hmac.New(sha256.New, tt.salt)
			b := make([]byte, 4)
			binary.BigEndian.PutUint32(b, tt.fn)
			hash.Write(b)
			binary.BigEndian.PutUint32(b, tt.idx)
			hash.Write(b)

			gt.Expect(salt(crypto.SHA256, tt.salt, tt.fn, tt.idx)).To(Equal(hash.Sum(nil)))
		})
	}
}

func TestMarshal(t *testing.T) {
	salt := fromHex(t, "000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f")
	noop := func(tx *txv1.Transaction) {}
	salted := func(tx *txv1.Transaction) { tx.Salt = salt }
	emptySalted := func(tx *txv1.Transaction) { tx.Reset(); tx.Salt = salt }
	shortSalt := func(tx *txv1.Transaction) { tx.Salt = tx.Salt[0:31] }
	noSigners := func(tx *txv1.Transaction) { tx.RequiredSigners = nil }
	nilElement := func(tx *txv1.Transaction) { tx.Inputs[0] = nil }
	reset := func(tx *txv1.Transaction) { tx.Reset() }
	unknownFields := func(tx *txv1.Transaction) { tx.Inputs[0].ProtoReflect().SetUnknown([]byte("garbage")) }

	var tests = map[string]struct {
		expected   ID
		errMatcher types.GomegaMatcher
		setup      func(*txv1.Transaction)
	}{
		"happy":         {fromHex(t, "77dc6e1729583cf7f1db9863b34a8951a3bb9369ab4cf0a86340ea92a8514cf5"), nil, noop},
		"changed salt":  {fromHex(t, "f1d081a486273dc66226a1fa30837e7583d9e0921163eac2300b350ff5ab4095"), nil, salted},
		"salted empty":  {fromHex(t, "38955e69c8db8963b3513c17631aebcf224c9c77017992dfe35a6dbba54b60a8"), nil, emptySalted},
		"empty vector":  {fromHex(t, "b2d9f8a592db1c411ee4ffe783c529af43172773107f5ed103535cd5e62ad1b4"), nil, noSigners},
		"nil element":   {fromHex(t, "d0ea4beee73ffa3597f00fe207e7c19f505247510da1dce3b7149d306d6d910a"), nil, nilElement},
		"empty":         {nil, MatchError("transaction salt is missing or less than 32 bytes in length"), reset},
		"short salt":    {nil, MatchError("transaction salt is missing or less than 32 bytes in length"), shortSalt},
		"unknown field": {nil, MatchError("protomsg: refusing to marshal unknown fields with length 7"), unknownFields},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gt := NewGomegaWithT(t)

			tx := newTestTransaction()
			tt.setup(tx)
			intTx, err := Marshal(crypto.SHA256, tx)
			if tt.errMatcher != nil {
				gt.Expect(err).To(tt.errMatcher)
				return
			}
			gt.Expect(err).NotTo(HaveOccurred())
			gt.Expect(intTx.Tx).To(Equal(tx))

			reflected, err := reflectTransactionID(crypto.SHA256, tx)
			gt.Expect(err).NotTo(HaveOccurred())

			gt.Expect(intTx.ID).To(Equal(tt.expected), "got %x want %x reflect %x", intTx.ID, tt.expected, reflected)

			expectedEncoded, err := protomsg.MarshalDeterministic(tx)
			gt.Expect(err).NotTo(HaveOccurred())
			gt.Expect(intTx.Encoded).To(Equal(expectedEncoded))
		})
	}
}

func TestIDMatchesReflected(t *testing.T) {
	gt := NewGomegaWithT(t)
	intTx, err := Marshal(crypto.SHA256, newTestTransaction())
	gt.Expect(err).NotTo(HaveOccurred())

	reflected, err := reflectTransactionID(crypto.SHA256, newTestTransaction())
	gt.Expect(err).NotTo(HaveOccurred())

	gt.Expect(intTx.ID).To(Equal(reflected))
}

func TestReflectDetectsChanges(t *testing.T) {
	gt := NewGomegaWithT(t)

	intTx, err := Marshal(crypto.SHA256, newTestTransaction())
	gt.Expect(err).NotTo(HaveOccurred())

	encoded, err := protomsg.MarshalDeterministic(newTestTransaction())
	gt.Expect(err).NotTo(HaveOccurred())

	var tests = map[string]struct {
		message    proto.Message
		errMatcher types.GomegaMatcher
	}{
		"no salt":       {&mutated.NoSaltTransaction{}, MatchError("transaction field number 1 must be a byte slice to use as a salt")},
		"removed field": {&mutated.RemovedFieldTransaction{}, nil},
		"extra field":   {&mutated.ExtraFieldTransaction{}, nil},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gt := NewGomegaWithT(t)

			err := proto.Unmarshal(encoded, tt.message)
			gt.Expect(err).NotTo(HaveOccurred())

			id, err := reflectTransactionID(crypto.SHA256, tt.message)
			if tt.errMatcher != nil {
				gt.Expect(err).To(tt.errMatcher)
				return
			}
			gt.Expect(err).NotTo(HaveOccurred())
			gt.Expect(id).NotTo(Equal(intTx.ID))
		})
	}
}

// The reflectTransactionID generates a transaction ID using metadata acquired
// by the protoreflect API. We want to avoid reflection in the production code
// its runtime cost. By using reflection in test, we can detect breaking
// changes to the structure of the transaction message that may require
// modifications to the production code.
func reflectTransactionID(h merkle.Hasher, tx proto.Message) (ID, error) {
	m := tx.ProtoReflect()

	// Retrieve the message fields and sort them by their field number.
	fields := m.Descriptor().Fields()
	var fds []protoreflect.FieldDescriptor
	for i := 0; i < fields.Len(); i++ {
		fds = append(fds, fields.Get(i))
	}
	sort.Slice(fds, func(i, j int) bool { return fds[i].Number() < fds[j].Number() })

	// Assert that field number 1 is a byte slice and use it as the seed.
	if fds[0].Number() != 1 {
		return nil, errors.New("transaction field number 1 must be a byte slice to use as a salt")
	}
	seed := m.Get(fds[0]).Interface().([]byte)

	// Iterate over the fields (except the salt). Each element should be a list.
	// For each element of the list, mashaled to the nonce for the field and
	// index. Finally, calculate the merkle root hash of the elements.
	var txLeaves [][]byte
	for _, fd := range fds[1:] {
		var listLeaves [][]byte
		for i, l := 0, m.Get(fd).List(); i < l.Len(); i++ {
			s := salt(h, seed, uint32(fd.Number()), uint32(i))
			msg := l.Get(i).Message().Interface()
			encoded, err := protomsg.MarshalDeterministic(msg)
			if err != nil {
				return nil, fmt.Errorf("failed to marhsal message %#v: %w", msg, err)
			}
			listLeaves = append(listLeaves, append(s, encoded...))
		}
		txLeaves = append(txLeaves, merkle.Root(h, listLeaves...))
	}
	return merkle.Root(h, txLeaves...), nil
}

func newTestTransaction() *txv1.Transaction {
	return &txv1.Transaction{
		Salt: []byte("NaCl - abcdefghijklmnopqrstuvwxyz"),
		Inputs: []*txv1.StateReference{
			{Txid: []byte("input-transaction-id-0"), OutputIndex: 1},
			{Txid: []byte("input-transaction-id-1"), OutputIndex: 0},
		},
		References: []*txv1.StateReference{
			{Txid: []byte("ref-transaction-id-0"), OutputIndex: 1},
			{Txid: []byte("ref-transaction-id-1"), OutputIndex: 0},
		},
		Outputs: []*txv1.State{
			{
				Info: &txv1.StateInfo{
					Owners: []*txv1.Party{{Credential: []byte("owner-1")}, {Credential: []byte("owner-2")}},
					Kind:   "state-kind-0",
				},
				State: []byte("state-0"),
			},
			{
				Info: &txv1.StateInfo{
					Owners: []*txv1.Party{{Credential: []byte("owner-1")}, {Credential: []byte("owner-2")}},
					Kind:   "state-kind-1",
				},
				State: []byte("state-1"),
			},
		},
		Parameters: []*txv1.Parameter{
			{Name: "name-0", Value: []byte("value-0")},
			{Name: "name-1", Value: []byte("value-1")},
		},
		RequiredSigners: []*txv1.Party{
			{Credential: []byte("observer-1")},
			{Credential: []byte("observer-2")},
		},
	}
}

func fromHex(t *testing.T, s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		t.Fatalf("failed to decode %q as hex string", s)
	}
	return b
}
