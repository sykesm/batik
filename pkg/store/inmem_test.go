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
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/sykesm/batik/pkg/btest"
)

func BenchmarkInmem(b *testing.B) {
	b.StopTimer()

	kv := NewInmem()
	defer btest.Close(b, kv)

	b.StartTimer()
	defer b.StopTimer()

	for i := 0; i < b.N; i++ {
		var randomKey [128]byte
		var randomValue [600]byte

		_, err := rand.Read(randomKey[:])
		require.NoError(b, err)
		_, err = rand.Read(randomValue[:])
		require.NoError(b, err)

		err = kv.Put(randomKey[:], randomValue[:])
		require.NoError(b, err)

		value, err := kv.Get(randomKey[:])
		require.NoError(b, err)

		require.EqualValues(b, randomValue[:], value)
	}
}

func TestInmemExistence(t *testing.T) {
	kv := NewInmem()
	defer btest.Close(t, kv)

	testExistence(t, kv)
}

func TestInmemKV(t *testing.T) {
	kv := NewInmem()
	defer btest.Close(t, kv)

	err := kv.Put([]byte("exist"), []byte("value"))
	require.NoError(t, err)

	v, err := kv.Get([]byte("exist"))
	require.NoError(t, err)
	require.Equal(t, []byte("value"), v)

	require.NoError(t, kv.Delete([]byte("exist")))
	_, err = kv.Get([]byte("exist"))
	require.Error(t, err)

	wb := kv.NewWriteBatch()
	require.NoError(t, wb.Put([]byte("key_batch1"), []byte("val_batch1")))
	require.NoError(t, wb.Put([]byte("key_batch2"), []byte("val_batch2")))
	require.NoError(t, wb.Put([]byte("key_batch3"), []byte("val_batch3")))
	require.NoError(t, wb.Commit())

	mv, err := kv.MultiGet([]byte("key_batch1"), []byte("key_batch2"), []byte("key_batch3"))
	require.NoError(t, err)
	require.Equal(t, [][]byte{[]byte("val_batch1"), []byte("val_batch2"), []byte("val_batch3")}, mv)
}

func TestInmemIdempotentClose(t *testing.T) {
	kv := NewInmem()
	require.NoError(t, kv.Close())
	require.NoError(t, kv.Close())
}

func TestInmemWriteBatch(t *testing.T) {
	t.Run("DeleteAndPut", func(t *testing.T) {
		kv := NewInmem()
		defer btest.Close(t, kv)

		wb := kv.NewWriteBatch()
		require.Equal(t, 0, wb.Count())

		require.NoError(t, wb.Delete([]byte("key_batch1")))
		require.Equal(t, 1, wb.Count())

		require.NoError(t, wb.Put([]byte("key_batch1"), []byte("val_batch1")))
		require.NoError(t, wb.Put([]byte("key_batch2"), []byte("val_batch2")))
		require.NoError(t, wb.Put([]byte("key_batch3"), []byte("val_batch3")))
		require.Equal(t, 4, wb.Count())

		require.NoError(t, wb.Commit())

		mv, err := kv.MultiGet([]byte("key_batch1"), []byte("key_batch2"), []byte("key_batch3"))
		require.NoError(t, err)
		require.Equal(t, [][]byte{[]byte("val_batch1"), []byte("val_batch2"), []byte("val_batch3")}, mv)
	})

	t.Run("Clear", func(t *testing.T) {
		kv := NewInmem()
		defer btest.Close(t, kv)

		wb := kv.NewWriteBatch()
		require.NoError(t, wb.Put([]byte("key_batch1"), []byte("val_batch1")))
		require.NoError(t, wb.Put([]byte("key_batch2"), []byte("val_batch2")))
		require.NoError(t, wb.Put([]byte("key_batch3"), []byte("val_batch3")))
		require.Equal(t, 3, wb.Count())

		wb.Clear()
		require.Equal(t, 0, wb.Count())
	})
}
