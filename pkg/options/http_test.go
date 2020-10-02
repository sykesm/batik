// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package options

import (
	"crypto/tls"
	"flag"
	"net/http"
	"testing"
	"time"

	. "github.com/onsi/gomega"
)

func TestHTTPServerDefaults(t *testing.T) {
	gt := NewGomegaWithT(t)
	httpServer := HTTPServerDefaults()
	gt.Expect(httpServer).To(Equal(&HTTPServer{
		ListenAddress:     ":8443",
		ReadHeaderTimeout: 30 * time.Second,
	}))
}

func TestHTTPServerApplyDefaults(t *testing.T) {
	tests := map[string]struct {
		setup func(*HTTPServer)
	}{
		"empty":               {setup: func(h *HTTPServer) { *h = HTTPServer{} }},
		"listen address":      {setup: func(h *HTTPServer) { h.ListenAddress = "" }},
		"read header timeout": {setup: func(h *HTTPServer) { h.ReadHeaderTimeout = 0 }},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gt := NewGomegaWithT(t)

			input := HTTPServerDefaults()
			tt.setup(input)

			input.ApplyDefaults()
			gt.Expect(input).To(Equal(HTTPServerDefaults()))
		})
	}
}

func TestHTTPServerFlagNames(t *testing.T) {
	gt := NewGomegaWithT(t)
	flags := (&HTTPServer{}).Flags()

	var names []string
	for _, f := range flags {
		names = append(names, f.Names()...)
	}

	gt.Expect(names).To(HaveLen(5))
	gt.Expect(names).To(ConsistOf(
		"http-listen-address",
		"http-read-timeout",
		"http-read-header-timeout",
		"http-write-timeout",
		"http-idle-timeout",
	))
}

func TestHTTPServerFlags(t *testing.T) {
	tests := map[string]struct {
		args     []string
		expected HTTPServer
	}{
		"no flags": {
			args:     []string{},
			expected: HTTPServer{},
		},
		"http-listen-address": {
			args:     []string{"--http-listen-address", "127.0.0.1:0"},
			expected: HTTPServer{ListenAddress: "127.0.0.1:0"},
		},
		"http timeouts": {
			args: []string{
				"--http-read-timeout", "1s",
				"--http-read-header-timeout", "1m",
				"--http-write-timeout", "30s",
				"--http-idle-timeout", "1h",
			},
			expected: HTTPServer{
				ReadTimeout:       time.Second,
				ReadHeaderTimeout: time.Minute,
				WriteTimeout:      30 * time.Second,
				IdleTimeout:       time.Hour,
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gt := NewGomegaWithT(t)

			httpServer := &HTTPServer{}
			flagSet := flag.NewFlagSet("http-server-test", flag.ContinueOnError)
			for _, f := range httpServer.Flags() {
				err := f.Apply(flagSet)
				gt.Expect(err).NotTo(HaveOccurred())
			}

			err := flagSet.Parse(tt.args)
			gt.Expect(err).NotTo(HaveOccurred())
			gt.Expect(httpServer).To(Equal(&tt.expected))
		})
	}
}

func TestHTTPServerFlagsDefaultText(t *testing.T) {
	flags := HTTPServerDefaults().Flags()
	assertWrappedFlagWithDefaultText(t, flags...)
}

func TestHTTPServerBuildServer(t *testing.T) {
	tests := map[string]struct {
		input     HTTPServer
		tlsConfig *tls.Config
		expected  *http.Server
	}{
		"empty": {
			input:     HTTPServer{},
			tlsConfig: nil,
			expected:  &http.Server{},
		},
		"with config": {
			input: HTTPServer{
				ListenAddress:     ":9999",
				ReadTimeout:       1,
				ReadHeaderTimeout: 2,
				WriteTimeout:      3,
				IdleTimeout:       4,
			},
			tlsConfig: nil,
			expected: &http.Server{
				Addr:              ":9999",
				ReadTimeout:       1,
				ReadHeaderTimeout: 2,
				WriteTimeout:      3,
				IdleTimeout:       4,
			},
		},
		"with TLS config": {
			input: HTTPServer{},
			tlsConfig: &tls.Config{
				ServerName: "test-server-name",
			},
			expected: &http.Server{
				TLSConfig: &tls.Config{
					ServerName: "test-server-name",
				},
			},
		},
		"with config and TLS": {
			input: HTTPServer{
				ListenAddress:     ":9999",
				ReadTimeout:       1,
				ReadHeaderTimeout: 2,
				WriteTimeout:      3,
				IdleTimeout:       4,
			},
			tlsConfig: &tls.Config{
				ServerName: "test-server-name",
			},
			expected: &http.Server{
				Addr:              ":9999",
				ReadTimeout:       1,
				ReadHeaderTimeout: 2,
				WriteTimeout:      3,
				IdleTimeout:       4,
				TLSConfig: &tls.Config{
					ServerName: "test-server-name",
				},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			gt := NewGomegaWithT(t)
			server := tt.input.BuildServer(tt.tlsConfig)
			gt.Expect(server).To(Equal(tt.expected))
		})
	}
}
