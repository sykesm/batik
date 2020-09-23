// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package pretty

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestParseUnixTime(t *testing.T) {
	gt := NewGomegaWithT(t)

	_, err := ParseUnixTime("invalidtime")
	gt.Expect(err.Error()).To(Equal("strconv.ParseFloat: parsing \"invalidtime\": invalid syntax"))

	tm, err := ParseUnixTime("1599593917.548589")
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Expect(tm.UnixNano()).To(Equal(int64(1599593917 * 1e9)))
}
