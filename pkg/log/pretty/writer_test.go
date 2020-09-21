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

	w := Writer{buf}

	testLine := "nonlogfmt line"
	_, err := w.Write([]byte(testLine))
	gt.Expect(err).To(MatchError("not a logfmt string"))

	testLine = `ts=1600356328.141956 level=info logger=batik caller=app/start.go:54 msg="Starting server"`
	expectedLine := `\x1b.*Sep 17 11:25:28.000000.*INFO.*|.*Starting server.*logger.*=.*batik.*caller.*=.*app/start.go:54.*`
	n, err := w.Write([]byte(testLine))
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Expect(n).To(Equal(145))
	gt.Expect(buf.String()).To(MatchRegexp(expectedLine))
}
