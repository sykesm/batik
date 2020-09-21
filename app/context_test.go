// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"bytes"
	"testing"

	. "github.com/onsi/gomega"
	zaplogfmt "github.com/sykesm/zap-logfmt"
	cli "github.com/urfave/cli/v2"
	"go.uber.org/zap"

	"github.com/sykesm/batik/pkg/log"
)

func TestContext_Logger(t *testing.T) {
	gt := NewGomegaWithT(t)

	ctx := cli.NewContext(cli.NewApp(), nil, nil)

	logger, err := GetLogger(ctx)
	gt.Expect(logger).NotTo(BeNil())
	gt.Expect(err).NotTo(HaveOccurred())

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
