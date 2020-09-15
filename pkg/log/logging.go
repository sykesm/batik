// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package log

import (
	"io"
	"os"

	zaplogfmt "github.com/sykesm/zap-logfmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
	// other string will be processed as "logfmt".
	Format string

	// Leveler controls the log levels that are enabled for the logging system. The
	// leveler is a zap.AtomicLevel that can dynamically reassign the log level.
	Leveler zap.AtomicLevel

	// Writer is the sink for encoded and formatted log records.
	//
	// If a Writer is not provided, os.Stderr will be used as the log sink.
	Writer io.Writer
}

func NewLeveler(logspec string) zap.AtomicLevel {
	return zap.NewAtomicLevelAt(NameToLevel(logspec))
}

func NewLogger(config Config, options ...zap.Option) *zap.Logger {
	return zap.New(zapcore.NewCore(
		NewEncoder(config.Format),
		NewWriteSyncer(config.Writer),
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
