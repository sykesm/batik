// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package totalorder

import (
	"context"
	"encoding/binary"
	"hash"
	"sync"

	"github.com/pkg/errors"

	"github.com/sykesm/batik/pkg/store"
)

type Hasher interface {
	New() hash.Hash
}

var (
	// Global prefixes.
	keyMetadata  = [...]byte{0x1}
	keySequences = [...]byte{0x2}

	// MD keys
	keyMetadataLastCommitted = append(keyMetadata[:], 0x1)
	keyMetadataAccumulator   = append(keyMetadata[:], 0x2)

	// Statically defined keys
	// TODO, probably make this a more concise byte string?
	keyLastCommittedSeq = []byte("last-committed")
)

type Store struct {
	mutex        sync.Mutex
	hasher       Hasher
	kv           store.KV
	nextSequence uint64
	accumulator  []byte
	waitCs       map[uint64]chan struct{}
}

func NewStore(hasher Hasher, kv store.KV) *Store {
	// TODO, crash recovery

	return &Store{
		kv:     kv,
		hasher: hasher,
		waitCs: map[uint64]chan struct{}{},
	}
}

func (s *Store) Append(t TXIDAndHMAC) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	tahBytes := t.serialize()

	h := s.hasher.New()
	h.Write(s.accumulator)
	h.Write(tahBytes)
	nextAccumulator := h.Sum(nil)

	batch := s.kv.NewWriteBatch()
	batch.Put(txKey(s.nextSequence), tahBytes)
	batch.Put(keyMetadataLastCommitted, uint64ToBytes(s.nextSequence))
	batch.Put(keyMetadataAccumulator, nextAccumulator)
	if err := batch.Commit(); err != nil {
		return errors.WithMessage(err, "could not persist tx")
	}
	waitC, ok := s.waitCs[s.nextSequence]
	if ok {
		close(waitC)
	}
	delete(s.waitCs, s.nextSequence)
	s.nextSequence++
	s.accumulator = nextAccumulator
	return nil
}

func (s *Store) Get(ctx context.Context, seq uint64) (TXIDAndHMAC, error) {
	waitC := s.waitC(seq)
	select {
	case <-waitC:
		value, err := s.kv.Get(txKey(seq))
		if err != nil {
			return TXIDAndHMAC{}, errors.WithMessagef(err, "could not get key for seq %d", seq)
		}
		return txidAndHMACFromBytes(value), nil
	case <-ctx.Done():
		return TXIDAndHMAC{}, ctx.Err()
	}
}

func (s *Store) waitC(seq uint64) <-chan struct{} {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if seq < s.nextSequence {
		closedC := make(chan struct{})
		close(closedC)
		return closedC
	}

	if waitC, ok := s.waitCs[seq]; ok {
		return waitC
	}

	waitC := make(chan struct{})
	s.waitCs[seq] = waitC
	return waitC
}

func uint64ToBytes(val uint64) []byte {
	byteValue := make([]byte, 8)
	binary.BigEndian.PutUint64(byteValue, val)
	return byteValue
}

func bytesToUint64(byteValue []byte) uint64 {
	return binary.BigEndian.Uint64(byteValue)
}

func txKey(sequence uint64) []byte {
	byteValue := append(keySequences[:], 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0)
	binary.BigEndian.PutUint64(byteValue[1:], sequence)
	return byteValue
}
