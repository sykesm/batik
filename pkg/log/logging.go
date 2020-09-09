// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package log

import (
	"io"
	"os"

	zaplogfmt "github.com/sykesm/zap-logfmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/sykesm/batik/pkg/log/encoder"
)

const (
	defaultFormat = "%{color}%{time:2006-01-02 15:04:05.000 MST} [%{module}] %{shortfunc} -> %{level:.4s} %{id:03x}%{color:reset} %{message}"
)

// Config is used to provide dependencies to a Logging instance.
type Config struct {
	// Name produces a named logger instance.
	//
	// If Name is not provided, the logger will be unnamed.
	Name string

	// Format is the log record format specifier for the Logging instance. If the
	// spec is the string "json", log records will be formatted as JSON. Any
	// other string will be provided to the FormatEncoder. Please see
	// encoder.ParseFormat for details on the supported verbs.
	//
	// If Format is not provided, a default format that provides basic information will
	// be used.
	Format string

	// LogSpec determines the log levels that are enabled for the logging system. The
	// spec must be in a format that can be processed by ActivateSpec.
	//
	// If LogSpec is not provided, loggers will be enabled at the INFO level.
	LogSpec string

	// Writer is the sink for encoded and formatted log records.
	//
	// If a Writer is not provided, os.Stderr will be used as the log sink.
	Writer io.Writer
}

func NewLogger(config Config, options ...zap.Option) (*zap.Logger, error) {
	e, err := NewEncoder(config.Format)
	if err != nil {
		return nil, err
	}

	w := NewWriteSyncer(config.Writer)

	l, err := NameToLevel(config.LogSpec)
	if err != nil {
		return nil, err
	}

	return zap.New(zapcore.NewCore(
		e,
		w,
		l,
	), append(defaultZapOptions(), options...)...,
	).Named(config.Name), nil
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
// be passed to the FormatEncoder.
func NewEncoder(format string) (zapcore.Encoder, error) {
	if format == "" {
		format = defaultFormat
	}

	var e zapcore.Encoder
	switch format {
	case "json":
		e = zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	case "logfmt":
		e = zaplogfmt.NewEncoder(zap.NewProductionEncoderConfig())
	default:
		// console
		formatters, err := encoder.ParseFormat(format)
		if err != nil {
			return nil, err
		}
		e = encoder.NewFormatEncoder(formatters...)
	}

	return e, nil
}

func defaultZapOptions() []zap.Option {
	return []zap.Option{
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	}
}
