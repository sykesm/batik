// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package store

import (
	"errors"
	"strconv"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/syndtr/goleveldb/leveldb"
)

func TestIsNotFound(t *testing.T) {
	tests := []struct {
		err        error
		isNotFound bool
	}{
		{err: nil, isNotFound: false},
		{err: errors.New("an error"), isNotFound: false},
		{err: leveldb.ErrNotFound, isNotFound: false},
		{err: &NotFoundError{}, isNotFound: true},
		{err: &NotFoundError{Err: errors.New("an error")}, isNotFound: true},
		{err: &NotFoundError{Err: leveldb.ErrNotFound}, isNotFound: true},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			gt := NewGomegaWithT(t)
			gt.Expect(IsNotFound(tt.err)).To(Equal(tt.isNotFound))
		})
	}
}

func TestIsAlreadyExists(t *testing.T) {
	tests := []struct {
		err             error
		isAlreadyExists bool
	}{
		{err: nil, isAlreadyExists: false},
		{err: errors.New("an error"), isAlreadyExists: false},
		{err: leveldb.ErrNotFound, isAlreadyExists: false},
		{err: &AlreadyExistsError{}, isAlreadyExists: true},
		{err: &AlreadyExistsError{Err: errors.New("an error")}, isAlreadyExists: true},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			gt := NewGomegaWithT(t)
			gt.Expect(IsAlreadyExists(tt.err)).To(Equal(tt.isAlreadyExists))
		})
	}
}
