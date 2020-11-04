// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package transaction

import (
	"crypto/hmac"
	"errors"

	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/proto"

	"github.com/sykesm/batik/pkg/merkle"
	txv1 "github.com/sykesm/batik/pkg/pb/tx/v1"
	"github.com/sykesm/batik/pkg/protomsg"
)

// Transaction holds intermediate information for an encoded transaction.
type Transaction struct {
	ID              ID                `json:"id"`
	Inputs          []*StateID        `json:"inputs"`
	References      []*StateID        `json:"references"`
	Outputs         []*State          `json:"outputs"`
	Parameters      []*Parameter      `json:"parameters"`
	RequiredSigners []*Party          `json:"required_signers"`
	Tx              *txv1.Transaction `json:"-"`
	Encoded         []byte            `json:"-"`
}

// New creates a Transaction from a protocol buffer message and also generates
// a transaction ID over the transaction. An error is returned if any element
// of the transaction cannot be marshaled into a protobuf message.
func New(h merkle.Hasher, tx *txv1.Transaction) (*Transaction, error) {
	// The transaction must be salted.
	if len(tx.Salt) < 32 {
		return nil, errors.New("transaction salt is missing or less than 32 bytes in length")
	}

	// fieldGetters is used instead of proto reflection to get the list of fields
	// and their associated field numbers when generating merkle hashes used for
	// transaction ID generation.
	fieldGetters := []func(*txv1.Transaction) (fn uint32, list interface{}){
		func(tx *txv1.Transaction) (uint32, interface{}) { return 2, tx.Inputs },
		func(tx *txv1.Transaction) (uint32, interface{}) { return 3, tx.References },
		func(tx *txv1.Transaction) (uint32, interface{}) { return 4, tx.Outputs },
		func(tx *txv1.Transaction) (uint32, interface{}) { return 5, tx.Parameters },
		func(tx *txv1.Transaction) (uint32, interface{}) { return 6, tx.RequiredSigners },
	}

	// The encoded transaction can be constructed from the encoded elements of
	// the transaction in order as they appear. Encoded elements are prepended
	// by the protowire encoded tag of the field number followed by the length
	// of the encoded message.
	var leaves [][]byte
	var encoded []byte
	encoded = append(encoded, encodedElement(1, tx.Salt)...)
	for _, getField := range fieldGetters {
		fn, list := getField(tx)
		m, err := marshalMessages(list)
		if err != nil {
			return nil, err
		}
		for i := range m {
			encoded = append(encoded, encodedElement(fn, m[i])...)
			m[i] = append(salt(h, tx.Salt, fn, uint32(i)), m[i]...)
		}
		leaves = append(leaves, merkle.Root(h, m...))
	}

	txid := merkle.Root(h, leaves...)
	return &Transaction{
		ID:              NewID(txid),
		Inputs:          ToStateIDs(tx.Inputs...),
		References:      ToStateIDs(tx.References...),
		Outputs:         ToStates(txid, tx.Outputs...),
		Parameters:      ToParameters(tx.Parameters...),
		RequiredSigners: ToParties(tx.RequiredSigners...),
		Tx:              tx,
		Encoded:         encoded,
	}, nil
}

// NewFromBytes creates a Transaction from the protocol buffer serialization
// of a txv1.Transaction.
func NewFromBytes(h merkle.Hasher, b []byte) (*Transaction, error) {
	var tx txv1.Transaction
	err := proto.Unmarshal(b, &tx)
	if err != nil {
		return nil, err
	}
	return New(h, &tx)
}

// Signed holds a transaction and (possibly unverified) signatures.
type Signed struct {
	*Transaction
	Signatures []*Signature
}

// encodedElement returns the encoded pieces of each message in a transaction.
// The element is prepended with the protowire encoded tag of the field number
// followed by the length of the encoded message.
//
// This logic loosely follows how protomsg.MarshalDeterministic encodes a
// proto.Message.
func encodedElement(fn uint32, m []byte) []byte {
	var encoded []byte
	encoded = protowire.AppendTag(encoded, protowire.Number(fn), protowire.BytesType)
	encoded = protowire.AppendBytes(encoded, m)
	return encoded
}

// A salt is generated for each leaf by caculating an HMAC over the protobuf
// field number and element index. The transaction salt is used as the key that
// seeds the HMAC.
func salt(h merkle.Hasher, seed []byte, fn, idx uint32) []byte {
	hash := hmac.New(h.New, seed)
	hash.Write([]byte{byte(fn >> 24), byte(fn >> 16), byte(fn >> 8), byte(fn)})     // big-endian field number
	hash.Write([]byte{byte(idx >> 24), byte(idx >> 16), byte(idx >> 8), byte(idx)}) // big-endian index
	return hash.Sum(nil)
}

// marshalMessages deterministically marshals a list of protobuf messages. The
// marshaled messages are returned in the same order as provided.
func marshalMessages(in ...interface{}) ([][]byte, error) {
	msgs, err := protomsg.ToMessageSlice(in...)
	if err != nil {
		return nil, err
	}

	var encoded [][]byte
	for _, m := range msgs {
		b, err := protomsg.MarshalDeterministic(m)
		if err != nil {
			return nil, err
		}
		encoded = append(encoded, b)
	}
	return encoded, nil
}
