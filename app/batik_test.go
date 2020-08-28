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
	gt.Expect(app.Flags[0].Names()[0]).To(Equal("config"))

	// Command implementations
	gt.Expect(app.Commands).NotTo(BeEmpty())
	gt.Expect(app.Commands[0].Name).To(Equal("start"))
	gt.Expect(app.Commands[1].Name).To(Equal("status"))
}

func TestBatikCommandNotFound(t *testing.T) {
	gt := NewGomegaWithT(t)

	stdin := bytes.NewBuffer(nil)
	stdout := bytes.NewBuffer(nil)
	stderr := bytes.NewBuffer(nil)

	app := Batik(nil, ioutil.NopCloser(stdin), stdout, stderr)
	app.ExitErrHandler = func(c *cli.Context, err error) {}

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
	app.ExitErrHandler = func(c *cli.Context, err error) {
		fmt.Fprintf(c.App.ErrWriter, "%+v\n", err)
	}

	err := app.Run([]string{"batik", "--config", "missing-file.txt"})
	gt.Expect(err).To(HaveOccurred())
	gt.Expect(err.(cli.ExitCoder).ExitCode()).To(Equal(3))
	gt.Expect(stdout.String()).To(BeEmpty())
	gt.Expect(stderr.String()).To(MatchRegexp("failed loading batik config:.*missing-file.txt"))
}

func TestBatikInteractive(t *testing.T) {
	gt := NewGomegaWithT(t)
	app := cli.NewApp()
	ctx := cli.NewContext(app, nil, nil)
	sa, err := shellApp(ctx)
	gt.Expect(err).NotTo(HaveOccurred())

	tests := []struct {
		command string
		stdout  string
		stderr  string
	}{
		{command: "", stdout: "", stderr: ""},
		{command: "exit", stdout: "", stderr: ""},
		{command: "help", stdout: sa.CustomAppHelpTemplate, stderr: ""},
		{command: "unknown-command", stdout: "", stderr: "Unknown command: unknown-command\n"},
	}
	for _, tt := range tests {
		gt := NewGomegaWithT(t)

		stdin := strings.NewReader(tt.command + "\n")
		stdout := bytes.NewBuffer(nil)
		stderr := bytes.NewBuffer(nil)

		app := Batik(nil, ioutil.NopCloser(stdin), stdout, stderr)
		err := app.Run([]string{"batik"})
		gt.Expect(err).NotTo(HaveOccurred())
		gt.Expect(stdout.String()).To(Equal(tt.stdout))
		gt.Expect(stderr.String()).To(Equal(tt.stderr))
	}
}
