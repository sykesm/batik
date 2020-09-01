// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package options

import (
	"flag"
	"strconv"
	"testing"
	"time"

	. "github.com/onsi/gomega"
)

func TestGRPCOptions(t *testing.T) {
	const MiB = 1024 * 1024

	tests := map[string]struct {
		args     []string
		expected GRPCOptions
	}{
		"no flags": {
			[]string{},
			GRPCOptions{MaxRecvMessageSize: 100 * MiB, MaxSendMessageSize: 100 * MiB},
		},
		"max recv": {
			[]string{"--max-recv-message-size", strconv.Itoa(90 * MiB)},
			GRPCOptions{MaxRecvMessageSize: 90 * MiB, MaxSendMessageSize: 100 * MiB},
		},
		"max send": {
			[]string{"--max-send-message-size", strconv.Itoa(50 * MiB)},
			GRPCOptions{MaxRecvMessageSize: 100 * MiB, MaxSendMessageSize: 50 * MiB},
		},
		"max recv and send": {
			[]string{
				"--max-recv-message-size", strconv.Itoa(1 * MiB),
				"--max-send-message-size", strconv.Itoa(2 * MiB),
			},
			GRPCOptions{MaxRecvMessageSize: 1 * MiB, MaxSendMessageSize: 2 * MiB},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gt := NewGomegaWithT(t)

			opts := DefaultGRPCOptions()
			flagSet := flag.NewFlagSet("grpc-options-test", flag.ContinueOnError)
			for _, f := range opts.Flags() {
				err := f.Apply(flagSet)
				gt.Expect(err).NotTo(HaveOccurred())
			}

			err := flagSet.Parse(tt.args)
			gt.Expect(err).NotTo(HaveOccurred())
			gt.Expect(opts).To(Equal(&tt.expected))
		})
	}
}

func TestGRPCServerOptions(t *testing.T) {
	const MiB = 1024 * 1024

	tests := map[string]struct {
		args     []string
		expected GRPCServerOptions
	}{
		"no flags": {
			[]string{},
			GRPCServerOptions{
				GRPCOptions:       *DefaultGRPCOptions(),
				ConnectionTimeout: 5 * time.Second,
			},
		},
		"conn timeout": {
			[]string{"--grpc-conn-timeout", "30s"},
			GRPCServerOptions{
				GRPCOptions:       *DefaultGRPCOptions(),
				ConnectionTimeout: 30 * time.Second,
			},
		},
		"max recv and send": {
			[]string{
				"--max-recv-message-size", strconv.Itoa(1 * MiB),
				"--max-send-message-size", strconv.Itoa(2 * MiB),
				"--grpc-conn-timeout", "1m",
			},
			GRPCServerOptions{
				GRPCOptions: GRPCOptions{
					MaxRecvMessageSize: 1 * MiB,
					MaxSendMessageSize: 2 * MiB,
				},
				ConnectionTimeout: time.Minute,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gt := NewGomegaWithT(t)

			opts := DefaultGRPCServerOptions()
			flagSet := flag.NewFlagSet("grpc-options-test", flag.ContinueOnError)
			for _, f := range opts.Flags() {
				err := f.Apply(flagSet)
				gt.Expect(err).NotTo(HaveOccurred())
			}

			err := flagSet.Parse(tt.args)
			gt.Expect(err).NotTo(HaveOccurred())
			gt.Expect(opts).To(Equal(&tt.expected))
		})
	}
}
