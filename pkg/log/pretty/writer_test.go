// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package pretty

import (
	"bytes"
	"fmt"
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
	expectedLine := `"\x1b[37mSep 17 11:25:28.000000\x1b[0m |\x1b[36mINFO\x1b[0m| \x1b[34mbatik\x1b[0m \x1b[0mapp/start.go:54\x1b[0m \x1b[97mStarting server\x1b[0m \n"`
	n, err := w.Write([]byte(testLine))
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Expect(n).To(Equal(113))
	fmt.Println(buf.String())
	gt.Expect(fmt.Sprintf("%q", buf.String())).To(Equal(expectedLine))
}
