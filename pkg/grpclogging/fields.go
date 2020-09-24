// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package grpclogging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type protoMarshaler struct {
	message proto.Message
}

func (m *protoMarshaler) MarshalJSON() ([]byte, error) {
	mo := protojson.MarshalOptions{
		AllowPartial:  true,
		UseProtoNames: true,
	}
	return mo.Marshal(m.message)
}

// ProtoMessage returns a zapcore.Field that attempts to log a protocol buffer
// message as JSON in the logs.
func ProtoMessage(key string, val interface{}) zapcore.Field {
	if pm, ok := val.(proto.Message); ok {
		return zap.Reflect(key, &protoMarshaler{message: pm})
	}
	return zap.Any(key, val)
}

// Error returns a zapcore.Field for a non-nil error. The errors is wrapped to
// hide any fmt.Formatter immplementation on the error.
func Error(err error) zapcore.Field {
	if err == nil {
		return zap.Skip()
	}

	// Wrap the error so it no longer implements fmt.Formatter. This will prevent
	// zap from adding the "verboseError" field to the log record that includes a
	// full stack trace.
	return zap.Error(struct{ error }{err})
}
