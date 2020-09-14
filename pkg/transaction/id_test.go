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
	"github.com/sykesm/batik/pkg/pb/transaction"
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
	noop := func(tx *transaction.Transaction) {}
	salted := func(tx *transaction.Transaction) { tx.Salt = []byte("sodium-chloride") }
	reset := func(tx *transaction.Transaction) { tx.Reset() }
	emptySalted := func(tx *transaction.Transaction) { tx.Reset(); tx.Salt = []byte("sodium-chloride") }
	noSigners := func(tx *transaction.Transaction) { tx.RequiredSigners = nil }
	nilElement := func(tx *transaction.Transaction) { tx.Inputs[0] = nil }
	unknownFields := func(tx *transaction.Transaction) { tx.Inputs[0].ProtoReflect().SetUnknown([]byte("garbage")) }

	var tests = map[string]struct {
		expected   []byte
		errMatcher types.GomegaMatcher
		setup      func(*transaction.Transaction)
	}{
		"happy":         {fromHex(t, "53e33ae87fb6cf2e4aaaabcdae3a93d578d9b7366e905dfff0446356774f726f"), nil, noop},
		"changed salt":  {fromHex(t, "a79fefe6edf500b30ab220e29de6fa22f9b1df876ce9d40a2dbee69c158ac491"), nil, salted},
		"empty":         {fromHex(t, "38955e69c8db8963b3513c17631aebcf224c9c77017992dfe35a6dbba54b60a8"), nil, reset},
		"salted empty":  {fromHex(t, "38955e69c8db8963b3513c17631aebcf224c9c77017992dfe35a6dbba54b60a8"), nil, emptySalted},
		"empty vector":  {fromHex(t, "d7469b8edbbccc748143923aae7073118ba4995a3b249b956a2b6366824c6c12"), nil, noSigners},
		"nil element":   {fromHex(t, "8b14c856d21568ba559d699bb15736bb41f83f12982c89dd33bd8d4149f8dc80"), nil, nilElement},
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
func reflectTransactionID(h merkle.Hasher, tx proto.Message) ([]byte, error) {
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

func newTestTransaction() *transaction.Transaction {
	return &transaction.Transaction{
		Inputs: []*transaction.StateReference{
			{Txid: []byte("input-transaction-id-0"), OutputIndex: 1},
			{Txid: []byte("input-transaction-id-1"), OutputIndex: 0},
		},
		References: []*transaction.StateReference{
			{Txid: []byte("ref-transaction-id-0"), OutputIndex: 1},
			{Txid: []byte("ref-transaction-id-1"), OutputIndex: 0},
		},
		Outputs: []*transaction.State{
			{
				Info: &transaction.StateInfo{
					Owners: []*transaction.Party{{Credential: []byte("owner-1")}, {Credential: []byte("owner-2")}},
					Kind:   "state-kind-0",
				},
				State: []byte("state-0"),
			},
			{
				Info: &transaction.StateInfo{
					Owners: []*transaction.Party{{Credential: []byte("owner-1")}, {Credential: []byte("owner-2")}},
					Kind:   "state-kind-1",
				},
				State: []byte("state-1"),
			},
		},
		Parameters: []*transaction.Parameter{
			{Name: "name-0", Value: []byte("value-0")},
			{Name: "name-1", Value: []byte("value-1")},
		},
		RequiredSigners: []*transaction.Party{
			{Credential: []byte("observer-1")},
			{Credential: []byte("observer-2")},
		},
		Salt: []byte("NaCl"),
	}
}

func fromHex(t *testing.T, s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		t.Fatalf("failed to decode %q as hex string", s)
	}
	return b
}
