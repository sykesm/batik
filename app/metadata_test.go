// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"testing"

	. "github.com/onsi/gomega"
	"github.com/urfave/cli/v2"
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
