// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	. "github.com/onsi/gomega"
	cli "github.com/urfave/cli/v2"

	"github.com/sykesm/batik/pkg/options"
)

func TestBatikWiring(t *testing.T) {
	gt := NewGomegaWithT(t)

	stdin := bytes.NewBuffer(nil)
	stdout := bytes.NewBuffer(nil)
	stderr := bytes.NewBuffer(nil)
	app := Batik(nil, ioutil.NopCloser(stdin), stdout, stderr)
	gt.Expect(app.Copyright).To(MatchRegexp("© Copyright IBM Corporation [\\d]{4}. All rights reserved."))

	// Global flags
	gt.Expect(app.Flags).NotTo(BeEmpty())
	gt.Expect(app.Flags[0].Names()[0]).To(Equal("color"))
	gt.Expect(app.Flags[1].Names()[0]).To(Equal("config"))
	gt.Expect(app.Flags[2].Names()[0]).To(Equal("log-spec"))

	// Command implementations
	gt.Expect(app.Commands).To(HaveLen(1))
	gt.Expect(app.Commands[0].Name).To(Equal("start"))
}

func TestBatikCommandNotFound(t *testing.T) {
	gt := NewGomegaWithT(t)

	stdin := bytes.NewBuffer(nil)
	stdout := bytes.NewBuffer(nil)
	stderr := bytes.NewBuffer(nil)

	app := Batik(nil, ioutil.NopCloser(stdin), stdout, stderr)
	app.ExitErrHandler = func(ctx *cli.Context, err error) {}

	err := app.Run([]string{"batik", "bogus-command"})
	gt.Expect(err).To(HaveOccurred())
	gt.Expect(err.(cli.ExitCoder).ExitCode()).To(Equal(2))
	gt.Expect(stdout.String()).To(BeEmpty())
	gt.Expect(stderr.String()).To(Equal("batik: 'bogus-command' is not a batik command. See `batik --help`.\n"))
}

func TestBatikConfigNotFound(t *testing.T) {
	gt := NewGomegaWithT(t)

	stdin := bytes.NewBuffer(nil)
	stdout := bytes.NewBuffer(nil)
	stderr := bytes.NewBuffer(nil)

	app := Batik(nil, ioutil.NopCloser(stdin), stdout, stderr)
	app.ExitErrHandler = func(ctx *cli.Context, err error) {
		fmt.Fprintf(ctx.App.ErrWriter, "%+v\n", err)
	}

	err := app.Run([]string{"batik", "--config", "missing-file.txt"})
	gt.Expect(err).To(HaveOccurred())
	gt.Expect(err.(cli.ExitCoder).ExitCode()).To(Equal(3))
	gt.Expect(stdout.String()).To(BeEmpty())
	gt.Expect(stderr.String()).To(MatchRegexp("unable to read config:.*missing-file.txt"))
}

func TestBatikInteractive(t *testing.T) {
	gt := NewGomegaWithT(t)
	app := cli.NewApp()
	ctx := cli.NewContext(app, nil, nil)
	sa, err := shellApp(ctx, options.BatikDefaults())
	gt.Expect(err).NotTo(HaveOccurred())

	tests := []struct {
		command string
		stdout  string
		stderr  string
	}{
		{command: "exit", stdout: "", stderr: ""},
		{command: "help", stdout: sa.CustomAppHelpTemplate, stderr: ""},
		{command: "unknown-command", stdout: "", stderr: "Unknown command: unknown-command\n"},
	}
	for _, tt := range tests {
		t.Run(tt.command, func(t *testing.T) {
			gt := NewGomegaWithT(t)

			stdin := strings.NewReader(tt.command + "\n")
			stdout := bytes.NewBuffer(nil)
			stderr := bytes.NewBuffer(nil)

			app := Batik(nil, ioutil.NopCloser(stdin), stdout, stderr)
			err := app.Run([]string{"batik"})
			gt.Expect(err).NotTo(HaveOccurred())
			gt.Expect(stdout.String()).To(Equal(tt.stdout))
			gt.Expect(stderr.String()).To(Equal(tt.stderr))
		})
	}
}

func TestBatikInteractiveWiring(t *testing.T) {
	gt := NewGomegaWithT(t)
	app := cli.NewApp()
	ctx := cli.NewContext(app, nil, nil)
	sa, err := shellApp(ctx, options.BatikDefaults())
	gt.Expect(err).NotTo(HaveOccurred())

	t.Run("AvailableCommands", func(t *testing.T) {
		gt := NewGomegaWithT(t)
		gt.Expect(sa.Commands).To(HaveLen(2))
		gt.Expect(sa.Commands[0].Name).To(Equal("exit"))
		gt.Expect(sa.Commands[1].Name).To(Equal("start"))
	})

	t.Run("HelpTemplate", func(t *testing.T) {
		gt := NewGomegaWithT(t)
		gt.Expect(strings.Split(strings.TrimSpace(sa.CustomAppHelpTemplate), "\n")).To(ConsistOf(
			"Commands:",
			"    exit   exit the shell",
			"    start  start the server",
		))
	})
}
