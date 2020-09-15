// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package log

import (
	"math"

	"go.uber.org/zap/zapcore"
)

const (
	// disabledLevel represents a disabled log level. Logs at this level should
	// never be emitted.
	disabledLevel = zapcore.Level(math.MinInt8)

	// payloadLevel is used to log the extremely detailed message level debug
	// information.
	payloadLevel = zapcore.Level(zapcore.DebugLevel - 1)

	defaultLevel = zapcore.InfoLevel
)

// NameToLevel converts a level name to a zapcore.Level.  If the level name is
// unknown or not provided, the default level is returned.
func NameToLevel(level string) zapcore.Level {
	switch level {
	case "PAYLOAD", "payload":
		return payloadLevel
	case "DEBUG", "debug":
		return zapcore.DebugLevel
	case "INFO", "info":
		return zapcore.InfoLevel
	case "WARNING", "WARN", "warning", "warn":
		return zapcore.WarnLevel
	case "ERROR", "error":
		return zapcore.ErrorLevel
	case "DPANIC", "dpanic":
		return zapcore.DPanicLevel
	case "PANIC", "panic":
		return zapcore.PanicLevel
	case "FATAL", "fatal":
		return zapcore.FatalLevel

	default:
		return defaultLevel
	}
}
