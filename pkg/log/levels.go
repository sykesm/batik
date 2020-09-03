// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package log

import (
	"fmt"
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
// unknown, an error is returned. If the level name is not provided, the default
// level is returned.
func NameToLevel(level string) (zapcore.Level, error) {
	switch level {
	case "PAYLOAD", "payload":
		return payloadLevel, nil
	case "DEBUG", "debug":
		return zapcore.DebugLevel, nil
	case "INFO", "info":
		return zapcore.InfoLevel, nil
	case "WARNING", "WARN", "warning", "warn":
		return zapcore.WarnLevel, nil
	case "ERROR", "error":
		return zapcore.ErrorLevel, nil
	case "DPANIC", "dpanic":
		return zapcore.DPanicLevel, nil
	case "PANIC", "panic":
		return zapcore.PanicLevel, nil
	case "FATAL", "fatal":
		return zapcore.FatalLevel, nil

	case "NOTICE", "notice":
		return zapcore.InfoLevel, nil // future
	case "CRITICAL", "critical":
		return zapcore.ErrorLevel, nil // future

	case "":
		return defaultLevel, nil

	default:
		return disabledLevel, fmt.Errorf("invalid log level: %s", level)
	}
}
