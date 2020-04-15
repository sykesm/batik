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

	db := NewInmem()
	defer btest.Close(b, db)

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

func TestExistence(t *testing.T) {
	db := NewInmem()
	defer btest.Close(t, db)

	_, err := db.Get([]byte("not_exist"))
	require.Error(t, err)

	err = db.Put([]byte("exist"), []byte{})
	require.NoError(t, err)

	val, err := db.Get([]byte("exist"))
	require.NoError(t, err)
	require.Equal(t, []byte{}, val)
}
