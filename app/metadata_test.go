// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"bytes"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/urfave/cli/v2"
	"github.com/sykesm/batik/pkg/log"
)

func TestMetadata_Config(t *testing.T) {
	gt := NewGomegaWithT(t)

	ctx := &cli.Context{
		App: &cli.App{
			Metadata: make(map[string]interface{}),
		},
	}

	config := GetConfig(ctx)
	gt.Expect(config).To(Equal(Config{}))

	expectedConfig := Config{
		Server: Server{
			Address: "127.0.0.1:9000",
		},
	}

	ctx.App.Metadata[string(configKey)] = expectedConfig

	config = GetConfig(ctx)
	gt.Expect(config).To(Equal(expectedConfig))

	newConfig := Config{
		Server: Server{
			Address: "127.0.0.1:9001",
		},
	}

	SetConfig(ctx, newConfig)

	gt.Expect(ctx.App.Metadata[string(configKey)]).To(Equal(newConfig))
}

func TestMetadata_Logger(t *testing.T) {
	gt := NewGomegaWithT(t)

	ctx := &cli.Context{
		App: &cli.App{
			Metadata: make(map[string]interface{}),
		},
	}

	logger, err := GetLogger(ctx)
	gt.Expect(logger).NotTo(BeNil())
	gt.Expect(err).NotTo(HaveOccurred())

	var buf bytes.Buffer
	newLogger, err := log.NewLogger(log.Config{
		Writer: &buf,
	})
	gt.Expect(err).NotTo(HaveOccurred())

	SetLogger(ctx, newLogger)
	gt.Expect(ctx.App.Metadata[string(loggerKey)]).To(Equal(newLogger))

	logger, err = GetLogger(ctx)
	gt.Expect(err).NotTo(HaveOccurred())
	logger.Info("test")

	gt.Expect(buf.String()).To(MatchRegexp("TestMetadata_Logger.*test"))
}
