// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package color

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestSprintColors(t *testing.T) {
	gt := NewGomegaWithT(t)
	gt.Expect(None.Sprint("test")).To(Equal("\x1b[0mtest\x1b[0m"))
	gt.Expect(Bold.Sprint("test")).To(Equal("\x1b[1mtest\x1b[0m"))
	gt.Expect(Faint.Sprint("test")).To(Equal("\x1b[2mtest\x1b[0m"))
	gt.Expect(Italic.Sprint("test")).To(Equal("\x1b[3mtest\x1b[0m"))
	gt.Expect(Underline.Sprint("test")).To(Equal("\x1b[4mtest\x1b[0m"))
	gt.Expect(BlinkSlow.Sprint("test")).To(Equal("\x1b[5mtest\x1b[0m"))
	gt.Expect(BlinkRapid.Sprint("test")).To(Equal("\x1b[6mtest\x1b[0m"))
	gt.Expect(ReverseVideo.Sprint("test")).To(Equal("\x1b[7mtest\x1b[0m"))
	gt.Expect(Concealed.Sprint("test")).To(Equal("\x1b[8mtest\x1b[0m"))
	gt.Expect(CrossedOut.Sprint("test")).To(Equal("\x1b[9mtest\x1b[0m"))

	gt.Expect(FgBlack.Sprint("test")).To(Equal("\x1b[30mtest\x1b[0m"))
	gt.Expect(FgRed.Sprint("test")).To(Equal("\x1b[31mtest\x1b[0m"))
	gt.Expect(FgGreen.Sprint("test")).To(Equal("\x1b[32mtest\x1b[0m"))
	gt.Expect(FgYellow.Sprint("test")).To(Equal("\x1b[33mtest\x1b[0m"))
	gt.Expect(FgBlue.Sprint("test")).To(Equal("\x1b[34mtest\x1b[0m"))
	gt.Expect(FgMagenta.Sprint("test")).To(Equal("\x1b[35mtest\x1b[0m"))
	gt.Expect(FgCyan.Sprint("test")).To(Equal("\x1b[36mtest\x1b[0m"))
	gt.Expect(FgWhite.Sprint("test")).To(Equal("\x1b[37mtest\x1b[0m"))

	gt.Expect(FgHiBlack.Sprint("test")).To(Equal("\x1b[90mtest\x1b[0m"))
	gt.Expect(FgHiRed.Sprint("test")).To(Equal("\x1b[91mtest\x1b[0m"))
	gt.Expect(FgHiGreen.Sprint("test")).To(Equal("\x1b[92mtest\x1b[0m"))
	gt.Expect(FgHiYellow.Sprint("test")).To(Equal("\x1b[93mtest\x1b[0m"))
	gt.Expect(FgHiBlue.Sprint("test")).To(Equal("\x1b[94mtest\x1b[0m"))
	gt.Expect(FgHiMagenta.Sprint("test")).To(Equal("\x1b[95mtest\x1b[0m"))
	gt.Expect(FgHiCyan.Sprint("test")).To(Equal("\x1b[96mtest\x1b[0m"))
	gt.Expect(FgHiWhite.Sprint("test")).To(Equal("\x1b[97mtest\x1b[0m"))

	gt.Expect(BgBlack.Sprint("test")).To(Equal("\x1b[40mtest\x1b[0m"))
	gt.Expect(BgRed.Sprint("test")).To(Equal("\x1b[41mtest\x1b[0m"))
	gt.Expect(BgGreen.Sprint("test")).To(Equal("\x1b[42mtest\x1b[0m"))
	gt.Expect(BgYellow.Sprint("test")).To(Equal("\x1b[43mtest\x1b[0m"))
	gt.Expect(BgBlue.Sprint("test")).To(Equal("\x1b[44mtest\x1b[0m"))
	gt.Expect(BgMagenta.Sprint("test")).To(Equal("\x1b[45mtest\x1b[0m"))
	gt.Expect(BgCyan.Sprint("test")).To(Equal("\x1b[46mtest\x1b[0m"))
	gt.Expect(BgWhite.Sprint("test")).To(Equal("\x1b[47mtest\x1b[0m"))

	gt.Expect(BgHiBlack.Sprint("test")).To(Equal("\x1b[100mtest\x1b[0m"))
	gt.Expect(BgHiRed.Sprint("test")).To(Equal("\x1b[101mtest\x1b[0m"))
	gt.Expect(BgHiGreen.Sprint("test")).To(Equal("\x1b[102mtest\x1b[0m"))
	gt.Expect(BgHiYellow.Sprint("test")).To(Equal("\x1b[103mtest\x1b[0m"))
	gt.Expect(BgHiBlue.Sprint("test")).To(Equal("\x1b[104mtest\x1b[0m"))
	gt.Expect(BgHiMagenta.Sprint("test")).To(Equal("\x1b[105mtest\x1b[0m"))
	gt.Expect(BgHiCyan.Sprint("test")).To(Equal("\x1b[106mtest\x1b[0m"))
	gt.Expect(BgHiWhite.Sprint("test")).To(Equal("\x1b[107mtest\x1b[0m"))
}
