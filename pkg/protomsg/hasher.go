// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package protomsg

import (
	"hash"

	"google.golang.org/protobuf/proto"
)

// A HashNewer is responsble for creating new instances of hash.Hash. All
// crypto.Hash implementaions from the standard library satisfy this interface.
type HashNewer interface {
	New() hash.Hash
}

// HashNewerFunc is an adapter to allow the use of a function as a HashNewer.
type HashNewerFunc func() hash.Hash

// New calls the handlder function h.
func (h HashNewerFunc) New() hash.Hash { return h() }

// A Hasher is responsible for deterministic hashing of a proto.Message.
type Hasher struct {
	Hash HashNewer
}

// HashMessage encodes a proto.Message and returns the hash of the result.
func (h *Hasher) HashMessage(sr proto.Message) ([]byte, error) {
	encoded, err := MarshalDeterministic(sr)
	if err != nil {
		return nil, err
	}

	hash := h.Hash.New()
	hash.Write(encoded)
	return hash.Sum(nil), nil
}

// HashMessages encodes and hashes a slice of objects that implement
// proto.Message.
func (h *Hasher) HashMessages(in interface{}) ([][]byte, error) {
	msgs, err := toMessageSlice(in)
	if err != nil {
		return nil, err
	}

	hashes := make([][]byte, len(msgs), len(msgs))
	for i := range msgs {
		hash, err := h.HashMessage(msgs[i])
		if err != nil {
			return nil, err
		}
		hashes[i] = hash
	}

	return hashes, nil
}
