// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package encoder

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestReset(t *testing.T) {
	gt := NewGomegaWithT(t)
	gt.Expect(ResetColor()).To(Equal("\x1b[0m"))
}

func TestNormalColors(t *testing.T) {
	gt := NewGomegaWithT(t)
	gt.Expect(ColorBlack.Normal()).To(Equal("\x1b[30m"))
	gt.Expect(ColorRed.Normal()).To(Equal("\x1b[31m"))
	gt.Expect(ColorGreen.Normal()).To(Equal("\x1b[32m"))
	gt.Expect(ColorYellow.Normal()).To(Equal("\x1b[33m"))
	gt.Expect(ColorBlue.Normal()).To(Equal("\x1b[34m"))
	gt.Expect(ColorMagenta.Normal()).To(Equal("\x1b[35m"))
	gt.Expect(ColorCyan.Normal()).To(Equal("\x1b[36m"))
	gt.Expect(ColorWhite.Normal()).To(Equal("\x1b[37m"))
}

func TestBoldColors(t *testing.T) {
	gt := NewGomegaWithT(t)
	gt.Expect(ColorBlack.Bold()).To(Equal("\x1b[30;1m"))
	gt.Expect(ColorRed.Bold()).To(Equal("\x1b[31;1m"))
	gt.Expect(ColorGreen.Bold()).To(Equal("\x1b[32;1m"))
	gt.Expect(ColorYellow.Bold()).To(Equal("\x1b[33;1m"))
	gt.Expect(ColorBlue.Bold()).To(Equal("\x1b[34;1m"))
	gt.Expect(ColorMagenta.Bold()).To(Equal("\x1b[35;1m"))
	gt.Expect(ColorCyan.Bold()).To(Equal("\x1b[36;1m"))
	gt.Expect(ColorWhite.Bold()).To(Equal("\x1b[37;1m"))
}
