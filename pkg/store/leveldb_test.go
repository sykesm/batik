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
	"errors"
	"fmt"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/sykesm/batik/pkg/tested"
)

func BenchmarkLevelDBFileSystem(b *testing.B) {
	path, cleanup := tested.TempDir(b, "", "level")
	defer cleanup()

	benchmarkLevelDB(b, path)
}

func BenchmarkLevelDBInMemory(b *testing.B) {
	benchmarkLevelDB(b, "")
}

func benchmarkLevelDB(b *testing.B, path string) {
	gt := NewGomegaWithT(b)
	db, err := NewLevelDB(path)
	gt.Expect(err).NotTo(HaveOccurred())
	defer tested.Close(b, db)

	b.ResetTimer()
	b.ReportAllocs()
	defer b.StopTimer()

	for i := 0; i < b.N; i++ {
		var randomKey [128]byte
		var randomValue [600]byte

		_, err := rand.Read(randomKey[:])
		gt.Expect(err).NotTo(HaveOccurred())
		_, err = rand.Read(randomValue[:])
		gt.Expect(err).NotTo(HaveOccurred())

		err = db.Put(randomKey[:], randomValue[:])
		gt.Expect(err).NotTo(HaveOccurred())

		value, err := db.Get(randomKey[:])
		gt.Expect(err).NotTo(HaveOccurred())
		gt.Expect(value).To(Equal(randomValue[:]))
	}
}

func testExistence(t *testing.T, kv KV) {
	gt := NewGomegaWithT(t)
	_, err := kv.Get([]byte("not_exist"))
	gt.Expect(err).To(MatchError(leveldb.ErrNotFound))
	gt.Expect(errors.Is(err, leveldb.ErrNotFound)).To(BeTrue())

	err = kv.Put([]byte("exist"), []byte{})
	gt.Expect(err).NotTo(HaveOccurred())

	val, err := kv.Get([]byte("exist"))
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Expect(val).To(Equal([]byte{}))
}

func TestLevelDBExistence(t *testing.T) {
	gt := NewGomegaWithT(t)
	path, cleanup := tested.TempDir(t, "", "level")
	defer cleanup()

	db, err := NewLevelDB(path)
	gt.Expect(err).NotTo(HaveOccurred())
	defer tested.Close(t, db)

	testExistence(t, db)
}

func TestLevelDBMemStorage(t *testing.T) {
	gt := NewGomegaWithT(t)
	db, err := NewLevelDB("")
	gt.Expect(err).NotTo(HaveOccurred())
	defer tested.Close(t, db)

	testExistence(t, db)
}

func TestLevelDB(t *testing.T) {
	gt := NewGomegaWithT(t)
	path, cleanup := tested.TempDir(t, "", "level")
	defer cleanup()

	db, err := NewLevelDB(path)
	gt.Expect(err).NotTo(HaveOccurred())

	err = db.Put([]byte("exist"), []byte("value"))
	gt.Expect(err).NotTo(HaveOccurred())

	wb := db.NewWriteBatch()
	gt.Expect(wb.Put([]byte("key_batch1"), []byte("val_batch1"))).To(Succeed())
	gt.Expect(wb.Put([]byte("key_batch2"), []byte("val_batch2"))).To(Succeed())
	gt.Expect(wb.Put([]byte("key_batch3"), []byte("val_batch3"))).To(Succeed())
	gt.Expect(wb.Commit()).To(Succeed())
	gt.Expect(db.Close()).To(Succeed())

	db2, err := NewLevelDB(path)
	gt.Expect(err).NotTo(HaveOccurred())

	v, err := db2.Get([]byte("exist"))
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Expect(v).To(Equal([]byte("value")))

	// Check multiget
	mv, err := db2.MultiGet([]byte("key_batch1"), []byte("key_batch2"))
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Expect(mv).To(Equal([][]byte{[]byte("val_batch1"), []byte("val_batch2")}))

	_, err = db2.MultiGet([]byte("missing"), []byte("key_batch2"))
	gt.Expect(err).To(MatchError(leveldb.ErrNotFound))
	gt.Expect(errors.Is(err, leveldb.ErrNotFound)).To(BeTrue())

	// Check delete
	gt.Expect(db2.Delete([]byte("exist"))).To(Succeed())

	_, err = db2.Get([]byte("exist"))
	gt.Expect(err).To(MatchError(leveldb.ErrNotFound))
	gt.Expect(errors.Is(err, leveldb.ErrNotFound)).To(BeTrue())
}

func TestLevelDBWriteBatch(t *testing.T) {
	path, cleanup := tested.TempDir(t, "", "level")
	defer cleanup()

	t.Run("CloseWithoutCommit", func(t *testing.T) {
		gt := NewGomegaWithT(t)
		db, err := NewLevelDB(path)
		gt.Expect(err).NotTo(HaveOccurred())

		wb := db.NewWriteBatch()
		for i := 0; i < 100000; i++ {
			gt.Expect(wb.Put([]byte(fmt.Sprintf("key_batch%d", i+1)), []byte(fmt.Sprintf("val_batch%d", i+1)))).To(Succeed())
		}
		gt.Expect(db.Close()).To(Succeed()) // Close without committing the WriteBatch

		db2, err := NewLevelDB(path)
		gt.Expect(err).NotTo(HaveOccurred())
		defer tested.Close(t, db2)

		_, err = db2.Get([]byte("key_batch100000"))
		gt.Expect(err).To(MatchError(leveldb.ErrNotFound))
		gt.Expect(errors.Is(err, leveldb.ErrNotFound)).To(BeTrue())
	})

	t.Run("CommitThenClose", func(t *testing.T) {
		gt := NewGomegaWithT(t)
		db, err := NewLevelDB(path)
		gt.Expect(err).NotTo(HaveOccurred())

		wb := db.NewWriteBatch()
		for i := 0; i < 100000; i++ {
			gt.Expect(wb.Put([]byte(fmt.Sprintf("key_batch%d", i+1)), []byte(fmt.Sprintf("val_batch%d", i+1)))).To(Succeed())
		}

		gt.Expect(wb.Commit()).To(Succeed())
		gt.Expect(db.Close()).To(Succeed())

		db2, err := NewLevelDB(path)
		gt.Expect(err).NotTo(HaveOccurred())
		defer tested.Close(t, db2)

		v, err := db2.Get([]byte("key_batch100000"))
		gt.Expect(err).NotTo(HaveOccurred())
		gt.Expect(v).To(Equal([]byte("val_batch100000")))
	})

	t.Run("DeleteAndPut", func(t *testing.T) {
		gt := NewGomegaWithT(t)
		path, cleanup := tested.TempDir(t, "", "level")
		defer cleanup()

		db, err := NewLevelDB(path)
		gt.Expect(err).NotTo(HaveOccurred())
		defer tested.Close(t, db)

		wb := db.NewWriteBatch()
		gt.Expect(wb.Count()).To(Equal(0))

		gt.Expect(wb.Delete([]byte("key_batch1"))).To(Succeed())
		gt.Expect(wb.Count()).To(Equal(1))

		gt.Expect(wb.Put([]byte("key_batch1"), []byte("val_batch1"))).To(Succeed())
		gt.Expect(wb.Put([]byte("key_batch2"), []byte("val_batch2"))).To(Succeed())
		gt.Expect(wb.Put([]byte("key_batch3"), []byte("val_batch3"))).To(Succeed())
		gt.Expect(wb.Count()).To(Equal(4))

		gt.Expect(wb.Commit()).To(Succeed())
	})

	t.Run("Clear", func(t *testing.T) {
		gt := NewGomegaWithT(t)
		path, cleanup := tested.TempDir(t, "", "level")
		defer cleanup()

		db, err := NewLevelDB(path)
		gt.Expect(err).NotTo(HaveOccurred())
		defer tested.Close(t, db)

		wb := db.NewWriteBatch()
		gt.Expect(wb.Put([]byte("key_batch1"), []byte("val_batch1"))).To(Succeed())
		gt.Expect(wb.Put([]byte("key_batch2"), []byte("val_batch2"))).To(Succeed())
		gt.Expect(wb.Put([]byte("key_batch3"), []byte("val_batch3"))).To(Succeed())
		gt.Expect(wb.Count()).To(Equal(3))

		wb.Clear()
		gt.Expect(wb.Count()).To(Equal(0))
	})
}

func TestLevelDBIterator(t *testing.T) {
	path, cleanup := tested.TempDir(t, "", "level")
	defer cleanup()

	t.Run("Keys", func(t *testing.T) {
		gt := NewGomegaWithT(t)
		db, err := NewLevelDB(path)
		gt.Expect(err).NotTo(HaveOccurred())
		defer tested.Close(t, db)

		err = db.Put([]byte("a"), []byte{})
		gt.Expect(err).NotTo(HaveOccurred())
		db.Put([]byte("b"), []byte{})
		gt.Expect(err).NotTo(HaveOccurred())
		db.Put([]byte("c"), []byte{})
		gt.Expect(err).NotTo(HaveOccurred())

		iter := db.NewIterator(nil, nil)
		keys, err := iter.Keys()
		gt.Expect(err).NotTo(HaveOccurred())

		gt.Expect(keys).To(ConsistOf(
			Key([]byte("a")),
			Key([]byte("b")),
			Key([]byte("c")),
		))
	})

	t.Run("KeysWithPrefix", func(t *testing.T) {
		gt := NewGomegaWithT(t)
		db, err := NewLevelDB(path)
		gt.Expect(err).NotTo(HaveOccurred())
		defer tested.Close(t, db)

		err = db.Put([]byte("a.a"), []byte{})
		gt.Expect(err).NotTo(HaveOccurred())
		db.Put([]byte("a.b"), []byte{})
		gt.Expect(err).NotTo(HaveOccurred())
		db.Put([]byte("b.c"), []byte{})
		gt.Expect(err).NotTo(HaveOccurred())

		iter := db.NewIterator([]byte("a."), nil)
		keys, err := iter.Keys()
		gt.Expect(err).NotTo(HaveOccurred())

		gt.Expect(keys).To(ConsistOf(
			Key([]byte("a.a")),
			Key([]byte("a.b")),
		))
	})
}
