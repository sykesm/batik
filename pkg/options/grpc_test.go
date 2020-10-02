// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package options

import (
	"flag"
	"reflect"
	"testing"
	"time"

	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
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
	flags := (&GRPC{}).Flags()

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
			for _, f := range grpc.Flags() {
				err := f.Apply(flagSet)
				gt.Expect(err).NotTo(HaveOccurred())
			}

			err := flagSet.Parse(tt.args)
			gt.Expect(err).NotTo(HaveOccurred())
			gt.Expect(grpc).To(Equal(&tt.expected))
		})
	}
}

func TestGRPCFlagsDefaultText(t *testing.T) {
	flags := GRPCDefaults().Flags()
	assertWrappedFlagWithDefaultText(t, flags...)
}

func TestGRPCServerDefaults(t *testing.T) {
	gt := NewGomegaWithT(t)
	gs := GRPCServerDefaults()
	gt.Expect(gs).To(Equal(&GRPCServer{
		GRPC:          *GRPCDefaults(),
		ConnTimeout:   30 * time.Second,
		ListenAddress: ":9443",
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
	flags := (&GRPCServer{}).Flags()

	var names []string
	for _, f := range flags {
		names = append(names, f.Names()...)
	}

	gt.Expect(flags).To(HaveLen(4))
	gt.Expect(names).To(ConsistOf(
		"grpc-conn-timeout",
		"grpc-listen-address",
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
		"grpc listen address": {
			args: []string{
				"--grpc-listen-address", "127.0.0.1:9999",
			},
			expected: GRPCServer{
				ListenAddress: "127.0.0.1:9999",
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gt := NewGomegaWithT(t)

			grpcServer := &GRPCServer{}
			flagSet := flag.NewFlagSet("grpc-server-test", flag.ContinueOnError)
			for _, f := range grpcServer.Flags() {
				err := f.Apply(flagSet)
				gt.Expect(err).NotTo(HaveOccurred())
			}

			err := flagSet.Parse(tt.args)
			gt.Expect(err).NotTo(HaveOccurred())
			gt.Expect(grpcServer).To(Equal(&tt.expected))
		})
	}
}

func TestGRPCServerFlagsDefaultText(t *testing.T) {
	flags := GRPCServerDefaults().Flags()
	assertWrappedFlagWithDefaultText(t, flags...)
}

func TestBuildServerOptions(t *testing.T) {
	gt := NewGomegaWithT(t)

	// Because of how the gRPC option functions work, the only sane way to
	// test them is to create and start a server with the options and drive
	// a client against it to test the limits. That's pretty heavy.
	//
	// A simpler, fragile, and hacky mechanism is to at least check we're
	// using the correct option functions. These checks can't make assertions
	// on the values - just that the correction option set is being used.
	var setters []uintptr
	for _, o := range GRPCServerDefaults().BuildServerOptions() {
		p := reflect.ValueOf(o).Elem().FieldByName("f").Pointer()
		setters = append(setters, p)
	}

	gt.Expect(setters).To(ConsistOf(
		reflect.ValueOf(grpc.MaxRecvMsgSize(0)).Elem().FieldByName("f").Pointer(),
		reflect.ValueOf(grpc.MaxSendMsgSize(0)).Elem().FieldByName("f").Pointer(),
		reflect.ValueOf(grpc.ConnectionTimeout(0)).Elem().FieldByName("f").Pointer(),
	))
}
