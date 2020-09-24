// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package grpclogging

import (
	"context"

	"go.uber.org/zap/zapcore"
)

type fieldKeyType struct{}

var fieldKey = &fieldKeyType{}

// Fields retrieves the zap fields populated by the gRPC interceptors from the
// provided context.
func Fields(ctx context.Context) []zapcore.Field {
	fields, ok := ctx.Value(fieldKey).([]zapcore.Field)
	if ok {
		return fields
	}
	return nil
}

// WithFields creates a new context decoarated with the provided
// zapcore.Fields.
func WithFields(ctx context.Context, fields []zapcore.Field) context.Context {
	return context.WithValue(ctx, fieldKey, fields)
}
