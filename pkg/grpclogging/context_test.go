// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package grpclogging

import (
	"context"
	"testing"
	"time"

	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestFields(t *testing.T) {
	fields := func() []zapcore.Field {
		return []zapcore.Field{
			zap.String("string-key", "string-value"),
			zap.Duration("duration-key", time.Second),
			zap.Int("int-key", 42),
		}
	}

	tests := map[string]struct {
		ctx      context.Context
		expected types.GomegaMatcher
	}{
		"empty context": {
			context.Background(),
			BeNil(),
		},
		"populated context": {
			WithFields(context.Background(), fields()),
			ConsistOf(fields()),
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gt := NewGomegaWithT(t)

			f := Fields(tt.ctx)
			gt.Expect(f).To(tt.expected)
		})
	}
}
