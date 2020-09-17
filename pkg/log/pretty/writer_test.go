// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package pretty

import (
	"bytes"
	"testing"

	. "github.com/onsi/gomega"
)

func TestWrite(t *testing.T) {
	gt := NewGomegaWithT(t)

	buf := &bytes.Buffer{}

	w := NewWriter(buf)

	testLine := "nonlogfmt line"
	n, err := w.Write([]byte(testLine))
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Expect(n).To(Equal(len(testLine)))
	gt.Expect(buf.Bytes()).To(MatchRegexp(testLine))

	testLine = `ts=1600356328.141956 level=info logger=batik caller=app/start.go:54 msg="Starting server"`
	expectedLine := `\x1b.*Sep 17 11:09:29.*INFO.*|.*Starting server.*logger.*=.*batik.*caller.*=.*app/start.go:54.*`
	n, err = w.Write([]byte(testLine))
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Expect(n).To(Equal(138))
	gt.Expect(buf.Bytes()).To(MatchRegexp(expectedLine))
}
