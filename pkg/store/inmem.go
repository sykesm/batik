// Copyright (c) 2019 Perlin
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package store

import (
	"bytes"
	"sync"

	"github.com/huandu/skiplist"
	"github.com/pkg/errors"
)

type kvOpType uint8

const (
	kvOpTypeDel kvOpType = iota + 1
	kvOpTypePut
)

type kvOp struct {
	opType     kvOpType
	key, value []byte
}

var _ WriteBatch = (*inmemWriteBatch)(nil)

type inmemWriteBatch struct {
	ops []kvOp
	kv  *InmemKV
}

// It is safe to modify the contents of the argument after Put returns but not
// before.
func (b *inmemWriteBatch) Put(key, value []byte) error {
	keyCopy := make([]byte, len(key))
	copy(keyCopy, key)

	valueCopy := make([]byte, len(value))
	copy(valueCopy, value)

	b.ops = append(b.ops, kvOp{opType: kvOpTypePut, key: keyCopy, value: valueCopy})

	return nil
}

func (b *inmemWriteBatch) Delete(key []byte) error {
	keyCopy := make([]byte, len(key))
	copy(keyCopy, key)

	b.ops = append(b.ops, kvOp{opType: kvOpTypeDel, key: keyCopy, value: nil})

	return nil
}

func (b *inmemWriteBatch) Commit() error {
	return b.kv.commitWriteBatch(b)
}

func (b *inmemWriteBatch) Clear() {
	b.ops = nil
}

func (b *inmemWriteBatch) Count() int {
	return len(b.ops)
}

var _ KV = (*InmemKV)(nil)

type InmemKV struct {
	sync.RWMutex
	db *skiplist.SkipList
}

func (s *InmemKV) Close() error {
	s.Lock()
	defer s.Unlock()

	// Do nothing if already closed
	if s.db == nil {
		return nil
	}

	s.db.Init()
	s.db = nil

	return nil
}

func (s *InmemKV) get(key []byte) ([]byte, error) {
	v, found := s.db.GetValue(key)
	if !found {
		return nil, ErrNotFound
	}

	src := v.([]byte)
	dest := make([]byte, len(src))
	copy(dest, src)

	return dest, nil
}

// The returned slice is its own copy, it is safe to modify the contents
// of the returned slice.
// It is safe to modify the contents of the argument after Get returns.
func (s *InmemKV) Get(key []byte) ([]byte, error) {
	s.RLock()
	defer s.RUnlock()

	return s.get(key)
}

// The returned slice is its own copy, it is safe to modify the contents
// of the returned slice.
// It is safe to modify the contents of the argument after Get returns.
func (s *InmemKV) MultiGet(keys ...[]byte) ([][]byte, error) {
	s.RLock()
	defer s.RUnlock()

	bufs := make([][]byte, 0, len(keys))

	for _, key := range keys {
		buf, err := s.get(key)
		if err != nil {
			return nil, err
		}

		bufs = append(bufs, buf)
	}

	return bufs, nil
}

// It is safe to modify the contents of the arguments after Put returns but not
// before.
func (s *InmemKV) Put(key, value []byte) error {
	s.Lock()
	defer s.Unlock()

	keyCopy := make([]byte, len(key))
	copy(keyCopy, key)

	valueCopy := make([]byte, len(value))
	copy(valueCopy, value)

	_ = s.db.Set(keyCopy, valueCopy)

	return nil
}

var (
	writeBatchPool = sync.Pool{
		New: func() interface{} {
			return new(inmemWriteBatch)
		},
	}
)

func (s *InmemKV) NewWriteBatch() WriteBatch {
	wb := writeBatchPool.Get().(*inmemWriteBatch)
	wb.kv = s
	return wb
}

func (s *InmemKV) commitWriteBatch(wb *inmemWriteBatch) error {
	s.Lock()
	defer s.Unlock()

	for _, op := range wb.ops {
		switch op.opType {
		case kvOpTypePut:
			_ = s.db.Set(op.key, op.value)
		case kvOpTypeDel:
			_ = s.db.Remove(op.key)
		default:
			return errors.Errorf("inmem: unknown op type %d", op.opType)
		}
	}

	wb.ops = nil
	writeBatchPool.Put(wb)

	return nil
}

// It is safe to modify the contents of the arguments after Delete returns but
// not before.
func (s *InmemKV) Delete(key []byte) error {
	s.Lock()
	defer s.Unlock()

	_ = s.db.Remove(key)

	return nil
}

func NewInmem() *InmemKV {
	var comparator skiplist.GreaterThanFunc = func(lhs, rhs interface{}) bool {
		return bytes.Compare(lhs.([]byte), rhs.([]byte)) == 1
	}

	return &InmemKV{db: skiplist.New(comparator)}
}
