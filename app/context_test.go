// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"bytes"
	"context"
	"testing"

	. "github.com/onsi/gomega"
	zaplogfmt "github.com/sykesm/zap-logfmt"
	cli "github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/sykesm/batik/pkg/log"
)

func TestContext_Logger(t *testing.T) {
	gt := NewGomegaWithT(t)

	ctx := cli.NewContext(cli.NewApp(), nil, nil)

	logger, err := GetLogger(ctx)
	gt.Expect(logger).NotTo(BeNil())
	gt.Expect(err).NotTo(HaveOccurred())

	ctx.Context = context.WithValue(ctx.Context, loggerKey, "invalidlogger")
	logger, err = GetLogger(ctx)
	gt.Expect(err).To(MatchError("logger not of type *zap.Logger"))

	buf := &bytes.Buffer{}
	encoder := zaplogfmt.NewEncoder(zap.NewProductionEncoderConfig())
	writer := log.NewWriteSyncer(buf)
	leveler := log.NewLeveler("info")
	newLogger := log.NewLogger(encoder, writer, leveler)

	SetLogger(ctx, newLogger)
	gt.Expect(ctx.Context.Value(loggerKey)).To(Equal(newLogger))

	logger, err = GetLogger(ctx)
	gt.Expect(err).NotTo(HaveOccurred())
	logger.Info("test")

	gt.Expect(buf.String()).To(MatchRegexp("msg=test"))
}

func TestContext_Leveler(t *testing.T) {
	gt := NewGomegaWithT(t)

	ctx := cli.NewContext(cli.NewApp(), nil, nil)

	leveler, err := GetLeveler(ctx)
	gt.Expect(err).To(MatchError("leveler does not exist"))

	ctx.Context = context.WithValue(ctx.Context, levelerKey, "invalidleveler")
	leveler, err = GetLeveler(ctx)
	gt.Expect(err).To(MatchError("leveler not of type zapcore.LevelEnabler"))

	buf := &bytes.Buffer{}
	encoder := zaplogfmt.NewEncoder(zap.NewProductionEncoderConfig())
	writer := log.NewWriteSyncer(buf)
	newLeveler := log.NewLeveler("info")
	logger := log.NewLogger(encoder, writer, newLeveler)

	SetLeveler(ctx, newLeveler)
	gt.Expect(ctx.Context.Value(levelerKey)).To(Equal(newLeveler))

	leveler, err = GetLeveler(ctx)
	gt.Expect(err).NotTo(HaveOccurred())
	logger.Info("test-info")
	logger.Debug("test-debug")

	gt.Expect(buf.String()).To(MatchRegexp("msg=test-info"))
	gt.Expect(buf.String()).NotTo(MatchRegexp("msg=test-debug"))

	l, ok := leveler.(zap.AtomicLevel)
	gt.Expect(ok).To(BeTrue())

	l.SetLevel(zapcore.DebugLevel)
	logger.Info("test-info")
	logger.Debug("test-debug")

	gt.Expect(buf.String()).To(MatchRegexp("msg=test-info"))
	gt.Expect(buf.String()).To(MatchRegexp("msg=test-debug"))
}
