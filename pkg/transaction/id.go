// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package transaction

import (
	"crypto/hmac"

	"github.com/sykesm/batik/pkg/merkle"
	"github.com/sykesm/batik/pkg/pb/transaction"
	"github.com/sykesm/batik/pkg/protomsg"
)

// ID generates a transaction ID from a Transaction message. An error is
// returned if any element of the transaction cannot be marshaled into a
// protobuf message.
func ID(h merkle.Hasher, tx *transaction.Transaction) ([]byte, error) {
	// fieldGetters is used instead of proto reflection to get the list of fields
	// and their associated field numbers when generating merkle hashes used for
	// transaction ID generation.
	fieldGetters := []func(*transaction.Transaction) (fn uint32, list interface{}){
		func(tx *transaction.Transaction) (uint32, interface{}) { return 2, tx.Inputs },
		func(tx *transaction.Transaction) (uint32, interface{}) { return 3, tx.References },
		func(tx *transaction.Transaction) (uint32, interface{}) { return 4, tx.Outputs },
		func(tx *transaction.Transaction) (uint32, interface{}) { return 5, tx.Parameters },
		func(tx *transaction.Transaction) (uint32, interface{}) { return 6, tx.RequiredSigners },
	}

	var leaves [][]byte
	for _, getField := range fieldGetters {
		fn, list := getField(tx)
		m, err := marshalMessages(list)
		if err != nil {
			return nil, err
		}
		for i := range m {
			m[i] = append(salt(h, tx.Salt, fn, uint32(i)), m[i]...)
		}
		leaves = append(leaves, merkle.Root(h, m...))
	}
	return merkle.Root(h, leaves...), nil
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
