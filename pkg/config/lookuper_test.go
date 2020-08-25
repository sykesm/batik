// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package config

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestMapLookuper(t *testing.T) {
	gt := NewGomegaWithT(t)

	ml := MapLookuper(nil)
	v, ok := ml.Lookup("key1")
	gt.Expect(ok).To(BeFalse())
	gt.Expect(v).To(BeEmpty())

	ml = MapLookuper(map[string]string{"key": "value"})
	v, ok = ml.Lookup("key")
	gt.Expect(ok).To(BeTrue())
	gt.Expect(v).To(Equal("value"))

	v, ok = ml.Lookup("missing")
	gt.Expect(ok).To(BeFalse())
	gt.Expect(v).To(BeEmpty())
}
