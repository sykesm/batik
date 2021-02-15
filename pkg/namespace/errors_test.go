// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package namespace

import (
	"testing"

	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
)

func TestErrHalt(t *testing.T) {
	he := newHaltError(errors.New("whoops!"), "message number %d", 1)

	gt := NewGomegaWithT(t)
	gt.Expect(errors.Is(he, ErrHalt)).To(BeTrue())
	gt.Expect(he).To(MatchError(ErrHalt))
	gt.Expect(he.Error()).To(Equal("message number 1: halt processing: whoops!"))
	unwrapped := he.(interface{ Unwrap() error }).Unwrap()
	gt.Expect(unwrapped).To(MatchError("halt processing: whoops!"))
	unwrapped = he.(interface{ Unwrap() error }).Unwrap()
	gt.Expect(unwrapped).To(MatchError(ErrHalt))
}
