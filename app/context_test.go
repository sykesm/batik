// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"bytes"
	"context"
	"flag"
	"testing"

	. "github.com/onsi/gomega"
	zaplogfmt "github.com/sykesm/zap-logfmt"
	cli "github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/sykesm/batik/pkg/log"
	"github.com/sykesm/batik/pkg/namespace"
)

func TestContext_Logger(t *testing.T) {
	gt := NewGomegaWithT(t)

	ctx := cli.NewContext(cli.NewApp(), nil, nil)
	logger, err := GetLogger(ctx)
	gt.Expect(logger).NotTo(BeNil())
	gt.Expect(err).NotTo(HaveOccurred())

	buf := bytes.NewBuffer(nil)
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

func TestContext_Namespaces(t *testing.T) {
	gt := NewGomegaWithT(t)

	fs := flag.NewFlagSet("test", 0)
	fs.String("namespace", "", "")
	ctx := cli.NewContext(cli.NewApp(), fs, nil)

	nss := GetNamespaces(ctx)
	gt.Expect(nss).To(BeNil())

	_, err := GetCurrentNamespace(ctx)
	gt.Expect(err).To(MatchError("could not find namespaces from context"))

	configNSS := map[string]*namespace.Namespace{
		"ns1": {},
		"ns2": {},
	}

	SetNamespaces(ctx, configNSS)

	nss = GetNamespaces(ctx)
	gt.Expect(nss).To(Equal(configNSS))

	_, err = GetCurrentNamespace(ctx)
	gt.Expect(err).To(MatchError("target namespace is not set"))

	err = ctx.Set("namespace", "missing")
	gt.Expect(err).NotTo(HaveOccurred())

	_, err = GetCurrentNamespace(ctx)
	gt.Expect(err).To(MatchError("namespace \"missing\" is not defined"))

	err = ctx.Set("namespace", "ns1")
	gt.Expect(err).NotTo(HaveOccurred())

	ns, err := GetCurrentNamespace(ctx)
	gt.Expect(err).NotTo(HaveOccurred())
	gt.Expect(ns).To(Equal(configNSS["ns1"]))
}
