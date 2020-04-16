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
	"crypto/rand"
	"fmt"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"github.com/sykesm/batik/pkg/tested"
)

func BenchmarkLevelDB(b *testing.B) {
	b.StopTimer()

	path, cleanup := tested.TempDir(b, "", "level")
	defer cleanup()

	db, err := NewLevelDB(path)
	require.NoError(b, err)
	defer tested.Close(b, db)

	b.StartTimer()
	defer b.StopTimer()

	for i := 0; i < b.N; i++ {
		var randomKey [128]byte
		var randomValue [600]byte

		_, err := rand.Read(randomKey[:])
		require.NoError(b, err)
		_, err = rand.Read(randomValue[:])
		require.NoError(b, err)

		err = db.Put(randomKey[:], randomValue[:])
		require.NoError(b, err)

		value, err := db.Get(randomKey[:])
		require.NoError(b, err)

		require.EqualValues(b, randomValue[:], value)
	}
}

func testExistence(t *testing.T, kv KV) {
	_, err := kv.Get([]byte("not_exist"))
	require.Error(t, err)

	err = kv.Put([]byte("exist"), []byte{})
	require.NoError(t, err)

	val, err := kv.Get([]byte("exist"))
	require.NoError(t, err)
	require.Equal(t, []byte{}, val)
}

func TestLevelDBExistence(t *testing.T) {
	path, cleanup := tested.TempDir(t, "", "level")
	defer cleanup()

	db, err := NewLevelDB(path)
	require.NoError(t, err)
	defer tested.Close(t, db)

	testExistence(t, db)
}

func TestLevelDBMemStorage(t *testing.T) {
	db, err := NewLevelDB("")
	require.NoError(t, err)
	defer tested.Close(t, db)

	testExistence(t, db)
}

func TestLevelDB(t *testing.T) {
	path, cleanup := tested.TempDir(t, "", "level")
	defer cleanup()

	db, err := NewLevelDB(path)
	require.NoError(t, err)

	err = db.Put([]byte("exist"), []byte("value"))
	require.NoError(t, err)

	wb := db.NewWriteBatch()
	require.NoError(t, wb.Put([]byte("key_batch1"), []byte("val_batch1")))
	require.NoError(t, wb.Put([]byte("key_batch2"), []byte("val_batch2")))
	require.NoError(t, wb.Put([]byte("key_batch3"), []byte("val_batch3")))
	require.NoError(t, wb.Commit())

	require.NoError(t, db.Close())

	db2, err := NewLevelDB(path)
	require.NoError(t, err)

	v, err := db2.Get([]byte("exist"))
	require.NoError(t, err)
	require.Equal(t, []byte("value"), v)

	// Check multiget
	mv, err := db2.MultiGet([]byte("key_batch1"), []byte("key_batch2"))
	require.NoError(t, err)
	require.Equal(t, [][]byte{[]byte("val_batch1"), []byte("val_batch2")}, mv)

	// Check delete
	require.NoError(t, db2.Delete([]byte("exist")))

	_, err = db2.Get([]byte("exist"))
	require.Error(t, err)
}

func TestLevelDBWriteBatch(t *testing.T) {
	path, cleanup := tested.TempDir(t, "", "level")
	defer cleanup()

	t.Run("CloseWithoutCommit", func(t *testing.T) {
		db, err := NewLevelDB(path)
		require.NoError(t, err)

		wb := db.NewWriteBatch()
		for i := 0; i < 100000; i++ {
			require.NoError(t, wb.Put([]byte(fmt.Sprintf("key_batch%d", i+1)), []byte(fmt.Sprintf("val_batch%d", i+1))))
		}
		require.NoError(t, db.Close()) // Close without committing the WriteBatch

		db2, err := NewLevelDB(path)
		require.NoError(t, err)
		defer tested.Close(t, db2)

		_, err = db2.Get([]byte("key_batch100000"))
		require.EqualError(t, errors.Cause(err), ErrNotFound.Error())
	})

	t.Run("CommitThenClose", func(t *testing.T) {
		db, err := NewLevelDB(path)
		require.NoError(t, err)

		wb := db.NewWriteBatch()
		for i := 0; i < 100000; i++ {
			require.NoError(t, wb.Put([]byte(fmt.Sprintf("key_batch%d", i+1)), []byte(fmt.Sprintf("val_batch%d", i+1))))
		}

		require.NoError(t, wb.Commit())
		require.NoError(t, db.Close())

		db2, err := NewLevelDB(path)
		require.NoError(t, err)
		defer tested.Close(t, db2)

		v, err := db2.Get([]byte("key_batch100000"))
		require.NoError(t, err)
		require.EqualValues(t, []byte("val_batch100000"), v)
	})

	t.Run("DeleteAndPut", func(t *testing.T) {
		path, cleanup := tested.TempDir(t, "", "level")
		defer cleanup()

		db, err := NewLevelDB(path)
		require.NoError(t, err)
		defer tested.Close(t, db)

		wb := db.NewWriteBatch()
		require.Equal(t, 0, wb.Count())

		require.NoError(t, wb.Delete([]byte("key_batch1")))
		require.Equal(t, 1, wb.Count())

		require.NoError(t, wb.Put([]byte("key_batch1"), []byte("val_batch1")))
		require.NoError(t, wb.Put([]byte("key_batch2"), []byte("val_batch2")))
		require.NoError(t, wb.Put([]byte("key_batch3"), []byte("val_batch3")))
		require.Equal(t, 4, wb.Count())

		require.NoError(t, wb.Commit())
	})

	t.Run("Clear", func(t *testing.T) {
		path, cleanup := tested.TempDir(t, "", "level")
		defer cleanup()

		db, err := NewLevelDB(path)
		require.NoError(t, err)
		defer tested.Close(t, db)

		wb := db.NewWriteBatch()
		require.NoError(t, wb.Put([]byte("key_batch1"), []byte("val_batch1")))
		require.NoError(t, wb.Put([]byte("key_batch2"), []byte("val_batch2")))
		require.NoError(t, wb.Put([]byte("key_batch3"), []byte("val_batch3")))
		require.Equal(t, 3, wb.Count())

		wb.Clear()
		require.Equal(t, 0, wb.Count())
	})
}
