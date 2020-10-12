// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package log

import (
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	defaultFormat = "logfmt"
)

func NewLeveler(logspec string) zap.AtomicLevel {
	return zap.NewAtomicLevelAt(NameToLevel(logspec))
}

func NewLogger(e zapcore.Encoder, w zapcore.WriteSyncer, l zapcore.LevelEnabler, options ...zap.Option) *zap.Logger {
	return zap.New(zapcore.NewCore(e, w, l), append(defaultZapOptions(), options...)...)
}

func NewWriteSyncer(w io.Writer) zapcore.WriteSyncer {
	if w == nil {
		w = os.Stderr
	}

	var sw zapcore.WriteSyncer
	switch w := w.(type) {
	case *os.File:
		fi, err := w.Stat()
		if err == nil && !fi.Mode().IsRegular() {
			// Add no-op sync method for special files
			sw = zapcore.AddSync(struct{ io.Writer }{w})
		}
		sw = zapcore.Lock(sw)
	case zapcore.WriteSyncer:
		sw = w
	default:
		sw = zapcore.AddSync(w)
	}

	return sw
}

func defaultZapOptions() []zap.Option {
	return []zap.Option{
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	}
}
