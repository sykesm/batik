// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package repl

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	. "github.com/onsi/gomega"
	cli "github.com/urfave/cli/v2"
)

func TestRun(t *testing.T) {
	commands := []*cli.Command{
		{
			Name: "echo",
			Action: func(ctx *cli.Context) error {
				fmt.Fprintf(ctx.App.Writer, "%v", ctx.Args().Slice())
				return nil
			},
		},
		{
			Name: "fail",
			Action: func(ctx *cli.Context) error {
				return errors.New("bummer...")
			},
		},
	}

	tests := []struct {
		name        string
		input       string
		expectedOut string
		expectedErr string
	}{
		{"empty", "", "", ""},
		{"echo", "echo arg1 arg2 arg3", "[arg1 arg2 arg3]", ""},
		{"echo-unterminated", `"arg1'`, "", "scanner: double quoted string not terminated\n"},
		{"fail", "fail", "", "command failed: bummer...\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gt := NewGomegaWithT(t)

			stdin := strings.NewReader(tt.input + "\n")
			stdout := &bytes.Buffer{}
			stderr := &bytes.Buffer{}

			app := cli.NewApp()
			app.Name = "repltest"
			app.HideVersion = true
			app.Commands = commands
			app.CommandNotFound = func(ctx *cli.Context, cmd string) {
				fmt.Printf("Command not found: %s\n", cmd)
			}

			r := New(app, WithStdin(stdin), WithStdout(stdout), WithStderr(stderr))
			r.Run(context.Background())

			gt.Expect(stdout.String()).To(Equal(tt.expectedOut))
			gt.Expect(stderr.String()).To(Equal(tt.expectedErr))
		})
	}
}

func TestErrExit(t *testing.T) {
	gt := NewGomegaWithT(t)
	gt.Expect(ErrExit.Error()).To(Equal("exit"))
	gt.Expect(ErrExit.ExitCode()).To(Equal(0))
}
