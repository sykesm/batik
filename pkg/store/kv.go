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
	"errors"
	"io"
)

// A NotFoundError indicates that a resource was not found.
type NotFoundError struct{ Err error }

func (n *NotFoundError) Error() string { return n.Err.Error() }
func (n *NotFoundError) Unwrap() error { return n.Err }

// IsNotFound returns a boolean indicating whether the error indicates
// that the resource was not found in the store.
func IsNotFound(err error) bool {
	var nfe *NotFoundError
	return errors.As(err, &nfe)
}

// An AlreadyExistsError indicates that a resource already exists.
type AlreadyExistsError struct{ Err error }

func (a *AlreadyExistsError) Error() string { return a.Err.Error() }
func (a *AlreadyExistsError) Unwrap() error { return a.Err }

// IsAlreadyExists returns a boolean indicating that the error indicates
// a resource already exists in the store.
func IsAlreadyExists(err error) bool {
	var aee *AlreadyExistsError
	return errors.As(err, &aee)
}

type Key []byte

type KV interface {
	io.Closer

	Get(key []byte) ([]byte, error)
	Put(key, value []byte) error
	Delete(key []byte) error

	NewWriteBatch() WriteBatch
}

type MultiGetter interface {
	MultiGet(keys ...[]byte) ([][]byte, error)
}

// WriteBatch batches a collection of put operations in memory before
// it's committed to disk.
//
// It's not guaranteed that all of the operations are kept in memory before
// the write batch is explicitly committed. It might be possible that the
// database decided commit the batch to disk earlier. For example, if a write
// batch is created, and 1000 put operations are batched, it might happen
// that while batching the 600th operation, the database decides to commit
// the first 599th operations first before proceeding.
type WriteBatch interface {
	Put(key, value []byte) error
	Delete(key []byte) error
	Commit() error

	Clear()
	Count() int
}

// Iterator iterates over a DB.
// The Iterator is not safe for concurrent use, but it is safe to use
// multiple iterators concurrently, with each in a dedicated goroutine.
// It is also safe to use an iterator concurrently with modifying its
// underlying DB. The resultant key/value pairs are guaranteed to be
// consistent.
//
// Also read Iterator documentation of the leveldb/iterator package.
type Iterator interface {
	Keys() ([]Key, error)
}
