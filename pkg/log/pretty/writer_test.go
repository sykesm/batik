// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package pretty

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	. "github.com/onsi/gomega"
	"go.uber.org/zap"
)

func TestWrite(t *testing.T) {
	t.Run("NonLogfmt", func(t *testing.T) {
		gt := NewGomegaWithT(t)
		buf := &bytes.Buffer{}
		w := NewWriter(buf, zap.NewProductionEncoderConfig(), func(s string) (time.Time, error) { return time.Time{}, nil })
		_, err := w.Write([]byte("this is not a logfmt line"))
		gt.Expect(err).To(MatchError("not a logfmt string"))
	})

	t.Run("Logfmt", func(t *testing.T) {
		gt := NewGomegaWithT(t)

		ts := time.Now()
		buf := &bytes.Buffer{}
		w := NewWriter(buf, zap.NewProductionEncoderConfig(), func(s string) (time.Time, error) { return ts, nil })

		testLine := `ts=1600356328.141956 level=info logger=batik caller=app/caller.go:99 msg="the message"`
		expectedLine := fmt.Sprintf("\x1b[37m%s\x1b[0m |\x1b[36mINFO\x1b[0m| \x1b[34mbatik\x1b[0m \x1b[0mapp/caller.go:99\x1b[0m \x1b[97mthe message\x1b[0m \n", ts.Format(time.StampMicro))

		n, err := w.Write([]byte(testLine))
		gt.Expect(err).NotTo(HaveOccurred())
		gt.Expect(n).To(Equal(110))
		gt.Expect(buf.String()).To(Equal(expectedLine))
	})
}
