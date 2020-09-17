// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package log

import (
	"io"
	"os"

	zaplogfmt "github.com/sykesm/zap-logfmt"
	"github.com/sykesm/batik/pkg/log/pretty"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	defaultFormat = "logfmt"
)

// Config is used to provide dependencies to a Logging instance.
type Config struct {
	// Name produces a named logger instance.
	//
	// If Name is not provided, the logger will be unnamed.
	Name string

	// Format is the log record format specifier for the Logging instance. If the
	// spec is the string "json", log records will be formatted as JSON. Any
	// other string will be processed as "logfmt". If color formatting is enabled
	// the format will be ignored.
	Format string

	// Leveler controls the log levels that are enabled for the logging system. The
	// leveler is a zap.AtomicLevel that can dynamically reassign the log level.
	Leveler zap.AtomicLevel

	// Writer is the sink for encoded and formatted log records.
	//
	// If a Writer is not provided, os.Stderr will be used as the log sink.
	Writer io.Writer

	// Color indicates whether color output post processing should be applied to logged
	// lines. Valid values are "yes" for forced color formatting, "no" for disabled color
	// formatting, and "auto" to automitcally process color formatting if the Writer is a tty.
	Color string
}

func NewLeveler(logspec string) zap.AtomicLevel {
	return zap.NewAtomicLevelAt(NameToLevel(logspec))
}

func NewLogger(config Config, options ...zap.Option) *zap.Logger {
	w := config.Writer

	switch config.Color {
	case "yes":
		w = pretty.NewWriter(w)
		config.Format = "logfmt"
	case "auto":
		if f, ok := w.(*os.File); ok && terminal.IsTerminal(int(f.Fd())) {
			w = pretty.NewWriter(w)
			config.Format = "logfmt"
		}
	}

	return zap.New(zapcore.NewCore(
		NewEncoder(config.Format),
		NewWriteSyncer(w),
		config.Leveler,
	), append(defaultZapOptions(), options...)...,
	).Named(config.Name)
}

func NewWriteSyncer(w io.Writer) zapcore.WriteSyncer {
	if w == nil {
		w = os.Stderr
	}

	var sw zapcore.WriteSyncer
	switch t := w.(type) {
	case *os.File:
		sw = zapcore.Lock(t)
	case zapcore.WriteSyncer:
		sw = t
	default:
		sw = zapcore.AddSync(w)
	}

	return sw
}

// NewEncoder returns a zapcore.Encoder based on the format string.
// Supported format strings include "json", and "logfmt". Any other string will
// default to "logfmt".
func NewEncoder(format string) zapcore.Encoder {
	switch format {
	case "json":
		return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	default:
		return zaplogfmt.NewEncoder(zap.NewProductionEncoderConfig())
	}
}

func defaultZapOptions() []zap.Option {
	return []zap.Option{
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	}
}
