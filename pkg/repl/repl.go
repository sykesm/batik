// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package repl

import (
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"text/scanner"
	"unicode"

	cli "github.com/urfave/cli/v2"
)

// A REPL implements an interactive Read, Evaluate, Print Loop for a cli.App.
type REPL struct {
	app    *cli.App
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
}

// An Option configures a REPL instance during construction.
type Option func(*REPL)

// WithStdin provides the io.Reader that command lines will be written from.
func WithStdin(stdin io.Reader) Option { return func(r *REPL) { r.stdin = stdin } }

// WithStdout provides the io.Writer that command output will be written to.
func WithStdout(stdout io.Writer) Option { return func(r *REPL) { r.stdout = stdout } }

// WithStderr provides the io.Writer that command errors will be written to.
func WithStderr(stderr io.Writer) Option { return func(r *REPL) { r.stderr = stderr } }

// New creates an instance of a Read, Evaluate, Print Loop control structure.
// In the absense of any options, input will be read from os.Stding, output
// will be written to os.Stdout, and errors will be written to os.Stderr.
func New(app *cli.App, opts ...Option) *REPL {
	r := &REPL{
		app:    app,
		stdin:  os.Stdin,
		stdout: os.Stdout,
		stderr: os.Stderr,
	}

	for _, o := range opts {
		o(r)
	}

	app.Writer = r.stdout
	app.ErrWriter = r.stderr

	return r
}

func newScanner(r io.Reader, stderr io.Writer) *scanner.Scanner {
	s := &scanner.Scanner{}
	s.Init(r)
	s.Mode ^= scanner.ScanChars      // don't scan go character literals ('\n' is a char literal)
	s.Mode ^= scanner.ScanComments   // don't scan go comments
	s.Mode ^= scanner.ScanRawStrings // don't scan go raw strings (`this is a raw string`)
	// s.Mode ^= scanner.ScanStrings    // don't scan go go strings ("string" is a string)
	s.Whitespace ^= 1 << '\n' // don't skip new lines
	s.IsIdentRune = func(ch rune, i int) bool {
		switch ch {
		case '_', '.', '-', ',':
			return true
		default:
			return unicode.IsLetter(ch) || unicode.IsDigit(ch)
		}
	}
	s.Error = func(s *scanner.Scanner, msg string) {
		fmt.Fprintf(stderr, "%s: %s\n", msg, s.TokenText())
	}

	return s
}

// TODO: The scanning functions should probabably all hang off an object that
// embeds the text/scanner that we interact with.

// scanCommandLine attempts to read and tokenize a command and arguments.
// Tokens are separted by whitespace and a unescaped newline characters
// indicate the end of the command.
//
// BUG(mjs): This does not properly handle basic shell semantics such as
// single quoted and multi-line quoted strings.
func scanCommandLine(s *scanner.Scanner) ([]string, error) {
	var args []string
	for {
		switch s.Scan() {
		case '\\':
			if ch := s.Next(); ch != '\n' {
				args = append(args, string(ch))
			}
		case '\n':
			return args, nil
		case scanner.String:
			str, err := strconv.Unquote(s.TokenText())
			if err != nil {
				return nil, err
			}
			args = append(args, str)
		case scanner.EOF:
			return nil, io.EOF
		default:
			args = append(args, s.TokenText())
		}
	}
}

// Run starts executing the REPL control loop.
func (r *REPL) Run(ctx context.Context) {
	s := newScanner(r.stdin, r.stderr)
	for {
		args, err := scanCommandLine(s)
		if err == io.EOF {
			return
		}
		if err != nil {
			fmt.Fprintf(r.app.ErrWriter, "scan command line: %s", err.Error())
			continue
		}
		if len(args) == 0 {
			continue
		}

		args = append([]string{r.app.Name}, args...)
		err = r.app.RunContext(ctx, args)
		if err != nil {
			fmt.Fprintf(r.app.ErrWriter, "command failed: %s", err.Error())
		}
	}
}
