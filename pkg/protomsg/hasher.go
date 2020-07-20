// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package protomsg

import (
	"crypto"

	"google.golang.org/protobuf/proto"
)

type Hasher struct {
	Hash crypto.Hash
}

func (h *Hasher) HashMessage(sr proto.Message) ([]byte, error) {
	encoded, err := MarshalDeterministic(sr)
	if err != nil {
		return nil, err
	}

	hash := h.Hash.New()
	hash.Write(encoded)
	return hash.Sum(nil), nil
}

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
