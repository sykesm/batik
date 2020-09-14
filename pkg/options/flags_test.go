// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package options

import (
	"testing"
	"time"

	. "github.com/onsi/gomega"

	cli "github.com/urfave/cli/v2"
)

func TestFlagApply(t *testing.T) {
	gt := NewGomegaWithT(t)

	var (
		duration time.Duration
		str      string
		ui       uint
	)

	app := cli.NewApp()
	app.Name = "flagtest"
	app.Action = func(ctx *cli.Context) error { return nil }
	app.Flags = []cli.Flag{
		NewDurationFlag(&cli.DurationFlag{Name: "duration", Value: duration, Destination: &duration}),
		NewStringFlag(&cli.StringFlag{Name: "string", Value: str, Destination: &str}),
		NewUintFlag(&cli.UintFlag{Name: "uint", Value: ui, Destination: &ui}),
	}

	// Simulate reading the config file by updating the flag destinations
	duration = time.Minute
	str = "updated-string"
	ui = 1234

	err := app.Run([]string{"flagtest"})
	gt.Expect(err).NotTo(HaveOccurred())

	gt.Expect(duration).To(Equal(time.Minute))
	gt.Expect(str).To(Equal("updated-string"))
	gt.Expect(ui).To(Equal(uint(1234)))

	err = app.Run([]string{
		"flagtest",
		"--duration", "1s",
		"--string", "flag-string",
		"--uint", "9876",
	})
	gt.Expect(err).NotTo(HaveOccurred())

	gt.Expect(duration).To(Equal(time.Second))
	gt.Expect(str).To(Equal("flag-string"))
	gt.Expect(ui).To(Equal(uint(9876)))
}
