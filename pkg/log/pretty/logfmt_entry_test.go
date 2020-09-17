// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package pretty

import (
	"bytes"
	"testing"
	"text/tabwriter"

	. "github.com/onsi/gomega"
)

func TestUnmarshalLogfmt(t *testing.T) {
	gt := NewGomegaWithT(t)

	testLine := []byte("nonlogfmt line")
	logfmtEntry := UnmarshalLogfmt(testLine)
	gt.Expect(logfmtEntry).To(BeNil())

	testLine = []byte(`ts=1600356328.141956 level=info logger=batik caller=app/start.go:54 msg="Starting server"`)
	logfmtEntry = UnmarshalLogfmt(testLine)
	gt.Expect(logfmtEntry).NotTo(BeNil())
}

func TestLogfmtEntry_Prettify(t *testing.T) {
	gt := NewGomegaWithT(t)

	buf := bytes.NewBuffer(nil)
	out := tabwriter.NewWriter(buf, 0, 1, 0, '\t', 0)
	logfmtEntry := &LogfmtEntry{
		buf:  buf,
		out:  out,
		Opts: DefaultOptions,
	}
	pretty := logfmtEntry.Prettify()
	gt.Expect(pretty).To(MatchRegexp(`\x1b.*Jan  1 00:00:00.*|.*|.*<no msg>.*`))

	testLine := []byte(`ts=1600356328.141956 level=info logger=batik caller=app/start.go:54 msg="Starting server"`)
	logfmtEntry = UnmarshalLogfmt(testLine)
	gt.Expect(logfmtEntry).NotTo(BeNil())

	expectedLine := `\x1b.*Sep 17 11:09:29.*INFO.*|.*Starting server.*logger.*=.*batik.*caller.*=.*app/start.go:54.*`
	pretty = logfmtEntry.Prettify()
	gt.Expect(pretty).To(MatchRegexp(expectedLine))
}
