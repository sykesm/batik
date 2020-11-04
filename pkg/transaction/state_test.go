// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package transaction

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestPartyStringer(t *testing.T) {
	gt := NewGomegaWithT(t)
	gt.Expect(Party{}.String()).To(Equal(""))
	gt.Expect(Party{PublicKey: []byte{1, 2, 3}}.String()).To(Equal("010203"))
}
