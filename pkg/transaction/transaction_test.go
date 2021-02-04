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
	"strconv"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/sykesm/batik/pkg/merkle"
	txv1 "github.com/sykesm/batik/pkg/pb/tx/v1"
	"github.com/sykesm/batik/pkg/protomsg"
	. "github.com/sykesm/batik/pkg/tested/matcher"
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

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
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

func TestNew(t *testing.T) {
	salt := fromHex(t, "000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f")
	long := fromHex(t, "000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f")
	noop := func(tx *txv1.Transaction) {}
	salted := func(tx *txv1.Transaction) { tx.Salt = salt }
	emptySalted := func(tx *txv1.Transaction) { tx.Reset(); tx.Salt = salt }
	shortSalt := func(tx *txv1.Transaction) { tx.Salt = tx.Salt[0:31] }
	noSigners := func(tx *txv1.Transaction) { tx.RequiredSigners = nil }
	nilElement := func(tx *txv1.Transaction) { tx.Inputs[0] = nil }
	reset := func(tx *txv1.Transaction) { tx.Reset() }
	unknownFields := func(tx *txv1.Transaction) { tx.Inputs[0].ProtoReflect().SetUnknown([]byte("garbage")) }
	longKey := func(tx *txv1.Transaction) { tx.Outputs[0].Info.Owners[0].PublicKey = long }

	tests := map[string]struct {
		expected   ID
		errMatcher types.GomegaMatcher
		setup      func(*txv1.Transaction)
	}{
		"happy":         {fromHex(t, "74ab83202b777ab9f27931fd76827cb848048e3abd70d2718cf1b60ed740bd89"), nil, noop},
		"changed salt":  {fromHex(t, "2b28f98d9ee6806fea1942ae130a6c75b8f9bf5a6ace24695f4a25e021a9af53"), nil, salted},
		"salted empty":  {fromHex(t, "38955e69c8db8963b3513c17631aebcf224c9c77017992dfe35a6dbba54b60a8"), nil, emptySalted},
		"empty vector":  {fromHex(t, "6c8847e9e9cd65a88e17116599d47ed0c521f4cc3fd8696bda2e41b1bb10733a"), nil, noSigners},
		"nil element":   {fromHex(t, "c85b907ec17ab566b36147964b65122e8468ea7944a0c40015552f17bb21a1f5"), nil, nilElement},
		"long key":      {fromHex(t, "bbecd8c2a804b25788f70cbe13aaca8c5d63b777a947d81cce15e853304e2ee0"), nil, longKey},
		"empty":         {nil, MatchError("transaction salt is missing or less than 32 bytes in length"), reset},
		"short salt":    {nil, MatchError("transaction salt is missing or less than 32 bytes in length"), shortSalt},
		"unknown field": {nil, MatchError("protomsg: refusing to marshal unknown fields with length 7"), unknownFields},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gt := NewGomegaWithT(t)

			tx := newTestTransaction()
			tt.setup(tx)
			intTx, err := New(crypto.SHA256, tx)
			if tt.errMatcher != nil {
				gt.Expect(err).To(tt.errMatcher)
				return
			}
			gt.Expect(err).NotTo(HaveOccurred())
			gt.Expect(intTx.Tx).To(Equal(tx))

			reflected, err := reflectTransactionID(crypto.SHA256, tx)
			gt.Expect(err).NotTo(HaveOccurred())

			gt.Expect(intTx.ID).To(Equal(tt.expected), "got %s want %s reflect %s", intTx.ID, tt.expected, reflected)

			expectedEncoded, err := protomsg.MarshalDeterministic(tx)
			gt.Expect(err).NotTo(HaveOccurred())
			gt.Expect(intTx.Encoded).To(Equal(expectedEncoded))

			// Ensure encoded transaction can be unmarshalled back to a Transaction
			unmarshalled := &txv1.Transaction{}
			err = proto.Unmarshal(intTx.Encoded, unmarshalled)
			gt.Expect(err).NotTo(HaveOccurred())
			gt.Expect(unmarshalled).To(ProtoEqual(tx))
		})
	}
}

func TestNewFromBytes(t *testing.T) {
	gt := NewGomegaWithT(t)
	intTx, err := New(crypto.SHA256, newTestTransaction())
	gt.Expect(err).NotTo(HaveOccurred())

	tx, err := NewFromBytes(crypto.SHA256, intTx.Encoded)
	gt.Expect(err).NotTo(HaveOccurred())

	gt.Expect(tx.ID).To(Equal(intTx.ID))
	gt.Expect(tx.Tx).To(ProtoEqual(intTx.Tx))
}

func TestNewFromBytesError(t *testing.T) {
	gt := NewGomegaWithT(t)
	_, err := NewFromBytes(crypto.SHA256, []byte("deaadbeef"))
	gt.Expect(err).To(MatchError(proto.Error))
}

func TestIDMatchesReflected(t *testing.T) {
	gt := NewGomegaWithT(t)
	intTx, err := New(crypto.SHA256, newTestTransaction())
	gt.Expect(err).NotTo(HaveOccurred())

	reflected, err := reflectTransactionID(crypto.SHA256, newTestTransaction())
	gt.Expect(err).NotTo(HaveOccurred())

	gt.Expect(intTx.ID).To(Equal(reflected))
}

func TestReflectDetectsChanges(t *testing.T) {
	gt := NewGomegaWithT(t)

	intTx, err := New(crypto.SHA256, newTestTransaction())
	gt.Expect(err).NotTo(HaveOccurred())

	encoded, err := protomsg.MarshalDeterministic(newTestTransaction())
	gt.Expect(err).NotTo(HaveOccurred())

	tests := map[string]struct {
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
					Owners: []*txv1.Party{{PublicKey: []byte("owner-1")}, {PublicKey: []byte("owner-2")}},
					Kind:   "state-kind-0",
				},
				State: []byte("state-0"),
			},
			{
				Info: &txv1.StateInfo{
					Owners: []*txv1.Party{{PublicKey: []byte("owner-1")}, {PublicKey: []byte("owner-2")}},
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
			{PublicKey: []byte("observer-1")},
			{PublicKey: []byte("observer-2")},
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
