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
	"github.com/pkg/errors"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/storage"
	"github.com/syndtr/goleveldb/leveldb/util"
)

var _ Iterator = (*leveldbIterator)(nil)

type leveldbIterator struct {
	iterator.Iterator
}

// Keys returns all keys that the iterator can iterate over.
// The underlying returned value can safely be modified and the
// iterator is released on return.
func (i *leveldbIterator) Keys() ([]Key, error) {
	defer i.Iterator.Release()

	var keys []Key
	for i.Iterator.Next() {
		keys = append(keys, append([]byte(nil), i.Iterator.Key()...))
	}

	return keys, i.Iterator.Error()
}

var _ WriteBatch = (*leveldbWriteBatch)(nil)

type leveldbWriteBatch struct {
	batch *leveldb.Batch
	kv    *LevelDBKV
}

func (b *leveldbWriteBatch) Put(key, value []byte) error {
	b.batch.Put(key, value)
	return nil
}

func (b *leveldbWriteBatch) Delete(key []byte) error {
	b.batch.Delete(key)
	return nil
}

func (b *leveldbWriteBatch) Commit() error {
	return b.kv.commitWriteBatch(b)
}

func (b *leveldbWriteBatch) Clear() {
	b.batch.Reset()
}

func (b *leveldbWriteBatch) Count() int {
	return b.batch.Len()
}

var _ KV = (*LevelDBKV)(nil)

type LevelDBKV struct {
	dir string
	db  *leveldb.DB
}

func (l *LevelDBKV) Close() error {
	return l.db.Close()
}

func (l *LevelDBKV) Get(key []byte) ([]byte, error) {
	v, err := l.db.Get(key, nil)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func (l *LevelDBKV) MultiGet(keys ...[]byte) ([][]byte, error) {
	var bufs = make([][]byte, len(keys))

	for i := range keys {
		b, err := l.Get(keys[i])
		if err != nil {
			return nil, err
		}

		bufs[i] = b
	}

	return bufs, nil
}

func (l *LevelDBKV) Put(key, value []byte) error {
	return l.db.Put(key, value, nil)
}

func (l *LevelDBKV) NewWriteBatch() WriteBatch {
	return &leveldbWriteBatch{
		batch: &leveldb.Batch{},
		kv:    l,
	}
}

// NewIterator returns an iterator that can be used to fetch Keys over a range
// from the DB. Prefix allows slicing the iterator to only contains keys in the given
// range. An empty prefix iterates over all keys in the DB.
//
// Values returned by the Iterator are cloned and can be safely modified.
func (l *LevelDBKV) NewIterator(prefix []byte, ro *opt.ReadOptions) Iterator {
	var bytesPrefix *util.Range
	if len(prefix) > 0 {
		bytesPrefix = util.BytesPrefix(prefix)
	}

	return &leveldbIterator{
		Iterator: l.db.NewIterator(bytesPrefix, ro),
	}
}

func (l *LevelDBKV) commitWriteBatch(wb *leveldbWriteBatch) error {
	return l.db.Write(wb.batch, nil)
}

func (l *LevelDBKV) Delete(key []byte) error {
	return l.db.Delete(key, nil)
}

func NewLevelDB(dir string) (*LevelDBKV, error) { // nolint:golint
	opts := &opt.Options{
		Filter:       filter.NewBloomFilter(10),
		NoWriteMerge: true,
	}

	var (
		db  *leveldb.DB
		err error
	)

	if len(dir) == 0 {
		db, err = leveldb.Open(storage.NewMemStorage(), opts)
		if err != nil {
			return nil, errors.Wrap(err, "failed to init leveldb")
		}
	} else {
		db, err = leveldb.OpenFile(dir, opts)
		if err != nil {
			return nil, errors.Wrap(err, "failed to init leveldb")
		}
	}

	return &LevelDBKV{
		dir: dir,
		db:  db,
	}, nil
}
