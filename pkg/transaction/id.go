// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package transaction

import (
	"crypto/hmac"

	"github.com/sykesm/batik/pkg/merkle"
	tb "github.com/sykesm/batik/pkg/pb/transaction"
	"github.com/sykesm/batik/pkg/protomsg"
	"google.golang.org/protobuf/encoding/protowire"
)

// IntermediateTx holds intermediate information for an encoded transaction.
type IntermediateTx struct {
	Tx      *tb.Transaction
	Encoded []byte
	ID      []byte
}

// Marshal encodes a Transaction message and also generates a transaction ID
// over the transaction. An error is returned if any element of the transaction cannot
// be marshaled into a protobuf message.
func Marshal(h merkle.Hasher, tx *tb.Transaction) (*IntermediateTx, error) {
	// fieldGetters is used instead of proto reflection to get the list of fields
	// and their associated field numbers when generating merkle hashes used for
	// transaction ID generation.
	fieldGetters := []func(*tb.Transaction) (fn uint32, list interface{}){
		func(tx *tb.Transaction) (uint32, interface{}) { return 2, tx.Inputs },
		func(tx *tb.Transaction) (uint32, interface{}) { return 3, tx.References },
		func(tx *tb.Transaction) (uint32, interface{}) { return 4, tx.Outputs },
		func(tx *tb.Transaction) (uint32, interface{}) { return 5, tx.Parameters },
		func(tx *tb.Transaction) (uint32, interface{}) { return 6, tx.RequiredSigners },
	}

	// The encoded transaction can be constructed from the encoded elements of
	// the transaction in order as they appear. Encoded elements are prepended
	// by the protowire encoded tag of the field number followed by the length
	// of the encoded message.
	var leaves [][]byte
	var encoded []byte
	if len(tx.Salt) != 0 {
		encoded = append(encoded, encodedElement(1, tx.Salt)...)
	}
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

	return &IntermediateTx{
		Tx:      tx,
		Encoded: encoded,
		ID:      merkle.Root(h, leaves...),
	}, nil
}

// encodedElement returns the encoded pieces of each message in a transaction.
// The element is prepended with the protowire encoded tag of the field number
// followed by the length of the encoded message.
//
// This logic loosely follows how protomsg.MarshalDeterministic encodes a proto.Message.
func encodedElement(fn uint32, m []byte) []byte {
	var encodedElement []byte

	encodedElement = append(encodedElement, byte(protowire.EncodeTag(protowire.Number(fn), protowire.BytesType)))
	encodedElement = append(encodedElement, byte(len(m)))
	encodedElement = append(encodedElement, m...)

	return encodedElement
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
