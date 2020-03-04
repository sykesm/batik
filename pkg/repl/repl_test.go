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
	"text/scanner"

	. "github.com/onsi/gomega"
	cli "github.com/urfave/cli/v2"
)

func TestScanner(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{`simple`, "simple", []string{"simple"}},
		{`word1 word2`, "word1 word2", []string{"word1", "word2"}},
		{`newline`, "new\nline", []string{"new", "\n", "line"}},
		{`newline`, "new\nline", []string{"new", "\n", "line"}},
		{`go comment`, "// comment", []string{"/", "/", "comment"}},
		{`id_with_underscores`, "id_with_underscore", []string{"id_with_underscore"}},
		{`id.with.dots`, "id.with.dots", []string{"id.with.dots"}},
		{`alphanum`, "abc 123 abc123", []string{"abc", "123", "abc123"}},
		{`decimal`, "1.2345", []string{"1.2345"}},
		{`"double quoted strings"`, `"double quoted strings"`, []string{`"double quoted strings"`}},
		{`'single quoted strings'`, `'single quoted strings'`, []string{"'", "single", "quoted", "strings", "'"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gt := NewGomegaWithT(t)

			reader := strings.NewReader(tt.input)
			s := newScanner(reader)
			var output []string
			for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
				output = append(output, s.TokenText())
			}
			gt.Expect(output).To(Equal(tt.expected))
			gt.Expect(s.ErrorCount).To(Equal(0))
		})
	}
}

func TestScanCommandLine(t *testing.T) {
	tests := []struct {
		input       string
		expected    []string
		expectedErr string
	}{
		{"command", []string{"command"}, ""},
		{"command subcommand", []string{"command", "subcommand"}, ""},
		{"command --flag1 arg1", []string{"command", "--flag1", "arg1"}, ""},
		{"command --flag1 comma,separated,args", []string{"command", "--flag1", "comma,separated,args"}, ""},
		{"command \"quoted argument string\"", []string{"command", "quoted argument string"}, ""},
		{"command \\\nacross many\\\nlines", []string{"command", "across", "many", "lines"}, ""},
		{`command \\ goo`, []string{"command", `\`, "goo"}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			gt := NewGomegaWithT(t)

			reader := strings.NewReader(tt.input + "\n")
			s := newScanner(reader)
			args, err := scanCommandLine(s)
			if tt.expectedErr == "" {
				gt.Expect(err).NotTo(HaveOccurred())
				gt.Expect(args).To(Equal(tt.expected))
			} else {
				gt.Expect(err).To(MatchError(tt.expectedErr))
			}
		})
	}
}

func TestRun(t *testing.T) {
	commands := []*cli.Command{
		{
			Name: "echo",
			Action: func(c *cli.Context) error {
				fmt.Fprintf(c.App.Writer, "%v", c.Args().Slice())
				return nil
			},
		},
		{
			Name: "fail",
			Action: func(c *cli.Context) error {
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
		{"", "", "", ""},
		{"echo", "echo arg1 arg2 arg3", "[arg1 arg2 arg3]", ""},
		{"fail", "fail", "", "command failed: bummer..."},
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
			app.CommandNotFound = func(c *cli.Context, cmd string) {
				fmt.Printf("Command not found: %s\n", cmd)
			}

			r := New(app, WithStdin(stdin), WithStdout(stdout), WithStderr(stderr))
			r.Run(context.Background())

			gt.Expect(stdout.String()).To(Equal(tt.expectedOut))
			gt.Expect(stderr.String()).To(Equal(tt.expectedErr))
		})
	}
}
