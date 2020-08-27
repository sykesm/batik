// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package repl

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"

	"github.com/peterh/liner"
	cli "github.com/urfave/cli/v2"
	"golang.org/x/crypto/ssh/terminal"
)

type exitError string

// ErrExit should be returned by commands that cleanly terminate the REPL.
const ErrExit = exitError("exit")

func (ee exitError) Error() string { return string(ee) }
func (ee exitError) ExitCode() int { return 0 }

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

// A prompter deals with issuing a command prmopt and reading the command line.
// This encapsulates the contract of liner.State.
type prompter interface {
	Prompt(prompt string) (string, error)
}

// A scanerPrompter is a prompter implementation that simply reads lines of text
// from the embedded scanner.
type scannerPrompter struct {
	scanner *bufio.Scanner
}

// Prompt is used to prompt and retrieve command line input.
func (s *scannerPrompter) Prompt(p string) (string, error) {
	if !s.scanner.Scan() && s.scanner.Err() == nil {
		return "", io.EOF
	}
	return s.scanner.Text(), nil
}

// readCommand reads and tokenizes a command and arguments.
func (r *REPL) readCommand(prompter prompter) ([]string, error) {
	scanner := &argScanner{}

	primaryPrompt := fmt.Sprintf("%s%% ", r.app.Name)
	secondaryPrompt := "> "

	prompt := primaryPrompt
	incomplete := false
	for {
		line, err := prompter.Prompt(prompt)
		if err == io.EOF && incomplete {
			break
		}
		if err == io.EOF {
			return nil, io.EOF
		}
		if err == liner.ErrPromptAborted {
			prompt = primaryPrompt
			scanner.Reset()
			continue
		}

		incomplete, err = scanner.ScanLine(line)
		if err != nil {
			return nil, err
		}
		if !incomplete {
			break
		}
		prompt = secondaryPrompt
	}

	return scanner.Args()
}

// isInteractiveTerminal returns true iff os.Stdin and os.Stdout are terminal
// devices.
func isInteractiveTerminal() bool {
	if !terminal.IsTerminal(int(os.Stdin.Fd())) {
		return false
	}
	if !terminal.IsTerminal(int(os.Stdout.Fd())) {
		return false
	}
	return true
}

// prompter returns an instance of a command prompt that can be used to read
// input from the command line. The liner package does not provide any
// mechanism to replace stdin or stdout streams so we use the scannerPrompter
// in environments without a terminal.
func (r *REPL) prompter() prompter {
	if r.stdin == os.Stdin && r.stdout == os.Stdout && isInteractiveTerminal() {
		rl := liner.NewLiner()
		rl.SetCtrlCAborts(true)
		rl.SetTabCompletionStyle(liner.TabPrints)

		// TODO(mjs): This is a really bare-bones starting point.
		rl.SetCompleter(func(line string) []string {
			l := strings.TrimLeftFunc(line, unicode.IsSpace)

			var candidates []string
			for _, command := range r.app.Commands {
				if strings.HasPrefix(command.Name, l) {
					candidates = append(candidates, command.Name)
				}
			}
			return candidates
		})
		return rl
	}

	return &scannerPrompter{
		scanner: bufio.NewScanner(r.stdin),
	}
}

// Run executes the REPL control loop. The method returns when any command returns
// ErrExit or when EOF is encountered on stdin.
func (r *REPL) Run(ctx context.Context) error {
	prompter := r.prompter()
	if closer, ok := prompter.(io.Closer); ok {
		defer closer.Close()
	}

	for {
		args, err := r.readCommand(prompter)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			fmt.Fprintf(r.app.ErrWriter, "%s\n", err.Error())
			continue
		}
		if len(args) == 0 {
			continue
		}

		args = append([]string{r.app.Name}, args...)
		err = r.app.RunContext(ctx, args)
		if err == ErrExit {
			return nil
		}
		if err != nil {
			fmt.Fprintf(r.app.ErrWriter, "command failed: %s\n", err.Error())
		}
	}
}
