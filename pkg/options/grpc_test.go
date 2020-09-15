// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package options

import (
	"flag"
	"testing"
	"time"

	. "github.com/onsi/gomega"
)

func TestGRPCDefaults(t *testing.T) {
	gt := NewGomegaWithT(t)
	grpc := GRPCDefaults()
	gt.Expect(grpc).To(Equal(&GRPC{
		MaxRecvMessageSize: 100 * 1024 * 1024,
		MaxSendMessageSize: 100 * 1024 * 1024,
	}))
}

func TestGRPCApplyDefaults(t *testing.T) {
	tests := map[string]struct {
		setup func(*GRPC)
	}{
		"empty":    {setup: func(g *GRPC) { *g = GRPC{} }},
		"max recv": {setup: func(g *GRPC) { g.MaxRecvMessageSize = 0 }},
		"max send": {setup: func(g *GRPC) { g.MaxSendMessageSize = 0 }},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gt := NewGomegaWithT(t)

			input := GRPCDefaults()
			tt.setup(input)

			input.ApplyDefaults()
			gt.Expect(input).To(Equal(GRPCDefaults()))
		})
	}
}

func TestGRPCFlagNames(t *testing.T) {
	gt := NewGomegaWithT(t)
	flags := (&GRPC{}).Flags("command name")

	var names []string
	for _, f := range flags {
		names = append(names, f.Names()...)
	}

	gt.Expect(flags).To(HaveLen(2))
	gt.Expect(names).To(ConsistOf(
		"grpc-max-recv-message-size",
		"grpc-max-send-message-size",
	))
}

func TestGRPCFlags(t *testing.T) {
	tests := map[string]struct {
		args     []string
		expected GRPC
	}{
		"no flags": {
			args:     []string{},
			expected: GRPC{},
		},
		"max recv": {
			args:     []string{"--grpc-max-recv-message-size", "9999"},
			expected: GRPC{MaxRecvMessageSize: 9999},
		},
		"max send": {
			args:     []string{"--grpc-max-send-message-size", "8888"},
			expected: GRPC{MaxSendMessageSize: 8888},
		},
		"max recv and send": {
			args: []string{
				"--grpc-max-recv-message-size", "1",
				"--grpc-max-send-message-size", "2",
			},
			expected: GRPC{MaxRecvMessageSize: 1, MaxSendMessageSize: 2},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gt := NewGomegaWithT(t)

			grpc := &GRPC{}
			flagSet := flag.NewFlagSet("grpc-test", flag.ContinueOnError)
			for _, f := range grpc.Flags("full command name") {
				err := f.Apply(flagSet)
				gt.Expect(err).NotTo(HaveOccurred())
			}

			err := flagSet.Parse(tt.args)
			gt.Expect(err).NotTo(HaveOccurred())
			gt.Expect(grpc).To(Equal(&tt.expected))
		})
	}
}

func TestGRPCServerDefaults(t *testing.T) {
	gt := NewGomegaWithT(t)
	gs := GRPCServerDefaults()
	gt.Expect(gs).To(Equal(&GRPCServer{
		GRPC:        *GRPCDefaults(),
		ConnTimeout: 30 * time.Second,
	}))
}

func TestGRPCServerApplyDefaults(t *testing.T) {
	tests := map[string]struct {
		setup func(*GRPCServer)
	}{
		"empty":        {setup: func(gs *GRPCServer) { *gs = GRPCServer{} }},
		"conn timeout": {setup: func(gs *GRPCServer) { gs.ConnTimeout = 0 }},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gt := NewGomegaWithT(t)

			input := GRPCServerDefaults()
			tt.setup(input)

			input.ApplyDefaults()
			gt.Expect(input).To(Equal(GRPCServerDefaults()))
		})
	}
}

func TestGRPCServerFlagNames(t *testing.T) {
	gt := NewGomegaWithT(t)
	flags := (&GRPCServer{}).Flags("command name")

	var names []string
	for _, f := range flags {
		names = append(names, f.Names()...)
	}

	gt.Expect(flags).To(HaveLen(3))
	gt.Expect(names).To(ConsistOf(
		"grpc-conn-timeout",
		"grpc-max-recv-message-size",
		"grpc-max-send-message-size",
	))
}

func TestGRPCServerFlags(t *testing.T) {
	tests := map[string]struct {
		args     []string
		expected GRPCServer
	}{
		"no flags": {
			args:     []string{},
			expected: GRPCServer{},
		},
		"conn timeout": {
			args:     []string{"--grpc-conn-timeout", "30s"},
			expected: GRPCServer{ConnTimeout: 30 * time.Second},
		},
		"max recv and send": {
			args: []string{
				"--grpc-max-recv-message-size", "1",
				"--grpc-max-send-message-size", "2",
				"--grpc-conn-timeout", "60s",
			},
			expected: GRPCServer{
				GRPC:        GRPC{MaxRecvMessageSize: 1, MaxSendMessageSize: 2},
				ConnTimeout: time.Minute,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gt := NewGomegaWithT(t)

			grpcServer := &GRPCServer{}
			flagSet := flag.NewFlagSet("grpc-server-test", flag.ContinueOnError)
			for _, f := range grpcServer.Flags("full command name") {
				err := f.Apply(flagSet)
				gt.Expect(err).NotTo(HaveOccurred())
			}

			err := flagSet.Parse(tt.args)
			gt.Expect(err).NotTo(HaveOccurred())
			gt.Expect(grpcServer).To(Equal(&tt.expected))
		})
	}
}
