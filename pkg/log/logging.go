// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package log

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"golang.org/x/crypto/ssh/terminal"
)

// FilteredWriter writes filtered logs to the designated level.
type FilteredWriter struct {
	w     zerolog.LevelWriter
	level zerolog.Level
}

func (w *FilteredWriter) Write(p []byte) (n int, err error) {
	return w.w.Write(p)
}

func (w *FilteredWriter) WriteLevel(level zerolog.Level, p []byte) (n int, err error) {
	if level >= w.level {
		return w.w.WriteLevel(level, p)
	}
	return len(p), nil
}

func NewDefaultErrLogger() *zerolog.Logger {
	l, _ := NewLogger("warn", "stderr")
	return l
}

func NewDefaultLogger() *zerolog.Logger {
	l, _ := NewLogger("info", "stdout")
	return l
}

// NewLogger creates a new *zerolog.Logger instance that writes logs at the indicated level to the logPath.
// logPath can be any of the following:
// * string representing a file path to write logs to
// * the string literal "stdout" or "stderr" indicating to write to the respective buffer
// * any io.Writer
func NewLogger(level string, logPath interface{}) (*zerolog.Logger, error) {
	logLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		return nil, err
	}

	w, err := logOutput(logLevel, logPath)
	if err != nil {
		return nil, err
	}

	logger := zerolog.New(w).With().Timestamp().Logger()

	return &logger, nil
}

// logOutput determines where we should send logs (if anywhere) and the log level.
func logOutput(level zerolog.Level, logPath interface{}) (io.Writer, error) {
	logOutput := ioutil.Discard

	if level == zerolog.Disabled {
		return logOutput, nil
	}

	switch logPath.(type) {
	case io.Writer:
		logOutput = logPath.(io.Writer)
	case string:
		if logPath == "stdout" || logPath == "stderr" {
			logOutput = parseStream(logPath.(string))

			// Pretty log if tty
			if terminal.IsTerminal(int(os.Stdout.Fd())) {
				logOutput = zerolog.ConsoleWriter{Out: logOutput, TimeFormat: time.RFC3339}
			}
		} else {
			var err error
			logOutput, err = os.OpenFile(logPath.(string), syscall.O_CREAT|syscall.O_RDWR|syscall.O_APPEND, 0666)
			if err != nil {
				return nil, err
			}
		}
	default:
		return nil, fmt.Errorf("unrecognized logPath: %s of type %T", logPath, logPath)
	}

	levelWriter := zerolog.MultiLevelWriter(logOutput)
	// A filtered writer allows us to set the log level for a specific logger instance instead of using the global
	// log level.
	filteredWriter := &FilteredWriter{levelWriter, level}
	w := zerolog.MultiLevelWriter(filteredWriter)

	return w, nil
}

func parseStream(stream string) *os.File {
	switch stream {
	case "stdout":
		return os.Stdout
	case "stderr":
		return os.Stderr
	}

	return nil
}
