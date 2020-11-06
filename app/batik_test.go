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
	"github.com/sykesm/batik/pkg/tested"
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
	gt.Expect(app.Flags[2].Names()[0]).To(Equal("data-dir"))
	gt.Expect(app.Flags[3].Names()[0]).To(Equal("log-format"))
	gt.Expect(app.Flags[4].Names()[0]).To(Equal("log-spec"))

	// Command implementations
	gt.Expect(app.Commands).To(HaveLen(2))
	gt.Expect(app.Commands[0].Name).To(Equal("db"))
	gt.Expect(app.Commands[1].Name).To(Equal("start"))

	// Subcommand implementations
	gt.Expect(app.Commands[0].Subcommands).To(HaveLen(3))
	gt.Expect(app.Commands[0].Subcommands[0].Name).To(Equal("get"))
	gt.Expect(app.Commands[0].Subcommands[0].Subcommands).To(HaveLen(2))
	gt.Expect(app.Commands[0].Subcommands[0].Subcommands[0].Name).To(Equal("state"))
	gt.Expect(app.Commands[0].Subcommands[0].Subcommands[0].Flags).To(HaveLen(1))
	gt.Expect(app.Commands[0].Subcommands[0].Subcommands[0].Flags[0].Names()[0]).To(Equal("consumed"))
	gt.Expect(app.Commands[0].Subcommands[0].Subcommands[1].Name).To(Equal("tx"))
	gt.Expect(app.Commands[0].Subcommands[1].Name).To(Equal("keys"))
	gt.Expect(app.Commands[0].Subcommands[1].Flags).To(HaveLen(1))
	gt.Expect(app.Commands[0].Subcommands[1].Flags[0].Names()[0]).To(Equal("prefix"))
	gt.Expect(app.Commands[0].Subcommands[2].Name).To(Equal("put"))
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
		{command: "logspec", stdout: "", stderr: ""},
		{command: "help", stdout: sa.CustomAppHelpTemplate, stderr: ""},
		{command: "unknown-command", stdout: "", stderr: "Unknown command: unknown-command\n"},
	}
	for _, tt := range tests {
		t.Run(tt.command, func(t *testing.T) {
			gt := NewGomegaWithT(t)

			path, cleanup := tested.TempDir(t, "", "level")
			defer cleanup()

			stdin := strings.NewReader(tt.command + "\n")
			stdout := bytes.NewBuffer(nil)
			stderr := bytes.NewBuffer(nil)

			app := Batik(nil, ioutil.NopCloser(stdin), stdout, stderr)
			err = app.Run([]string{"batik", "--data-dir=" + path})
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
		gt.Expect(sa.Commands).To(HaveLen(4))
		gt.Expect(sa.Commands[0].Name).To(Equal("db"))
		gt.Expect(sa.Commands[1].Name).To(Equal("exit"))
		gt.Expect(sa.Commands[2].Name).To(Equal("logspec"))
		gt.Expect(sa.Commands[3].Name).To(Equal("start"))

		gt.Expect(sa.Commands[0].Subcommands).To(HaveLen(3))
		gt.Expect(sa.Commands[0].Subcommands[0].Name).To(Equal("get"))
		gt.Expect(sa.Commands[0].Subcommands[0].Subcommands).To(HaveLen(2))
		gt.Expect(sa.Commands[0].Subcommands[0].Subcommands[0].Name).To(Equal("state"))
		gt.Expect(sa.Commands[0].Subcommands[0].Subcommands[0].Flags).To(HaveLen(1))
		gt.Expect(sa.Commands[0].Subcommands[0].Subcommands[0].Flags[0].Names()[0]).To(Equal("consumed"))
		gt.Expect(sa.Commands[0].Subcommands[0].Subcommands[1].Name).To(Equal("tx"))
		gt.Expect(sa.Commands[0].Subcommands[1].Name).To(Equal("keys"))
		gt.Expect(sa.Commands[0].Subcommands[1].Flags).To(HaveLen(1))
		gt.Expect(sa.Commands[0].Subcommands[1].Flags[0].Names()[0]).To(Equal("prefix"))
		gt.Expect(sa.Commands[0].Subcommands[2].Name).To(Equal("put"))
	})

	t.Run("HelpTemplate", func(t *testing.T) {
		gt := NewGomegaWithT(t)
		gt.Expect(strings.Split(strings.TrimSpace(sa.CustomAppHelpTemplate), "\n")).To(ConsistOf(
			"Commands:",
			"    db       perform operations against a kv store",
			"    exit     exit the shell",
			"    logspec  change the logspec of the logger leveler to any supported log level (eg. debug, info)",
			"    start    start the server",
		))
	})
}
