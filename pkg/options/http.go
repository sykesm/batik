// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package options

import (
	"crypto/tls"
	"net/http"
	"time"

	cli "github.com/urfave/cli/v2"
)

// HTTPServer exposes configuration for HTTP servers.
type HTTPServer struct {
	// ListenAddress determines the addresses that the server will listen on. The
	// address should be in a form that is compatible with net.Listen from the Go
	// standard library.
	ListenAddress string `yaml:"listen_address,omitempty"`
	// ReadTimeout defines the maximum duration for reading a request. This
	// timeout includes the header, body, and any connecton processing such as a
	// TLS handshake.  This timeout is not set by default.
	ReadTimeout time.Duration `yaml:"read_timeout,omitempty"`
	// ReadHeaderTimeout defines the maximum duration for reading a requset. This
	// timeout includes the header and any connection processing but does not
	// include time time it takes to read the request body.
	// The default value for this timeout is 30s.
	ReadHeaderTimeout time.Duration `yaml:"read_header_timeout,omitempty"`
	// WriteTimeout defines the maximum duration for processing and responding to
	// a request. This timeout inclues all processing from time time the last
	// request header was read. This timeout is not set by default.
	WriteTimeout time.Duration `yaml:"write_timeout,omitempty"`
	// IdleTimeout defines the maximum amount of time to wait for a new request when
	// keep-alives are enabled. This timeout is not set by default.
	IdleTimeout time.Duration `yaml:"idle_timeout,omitempty"`
}

// HTTPServerDefaults returns the default configuration values for HTTP
// servers.
func HTTPServerDefaults() *HTTPServer {
	return &HTTPServer{
		ListenAddress:     ":8443",
		ReadHeaderTimeout: 30 * time.Second,
	}
}

// ApplyDefaults applies default values for missing configuration fields.
func (h *HTTPServer) ApplyDefaults() {
	defaults := HTTPServerDefaults()
	if h.ListenAddress == "" {
		h.ListenAddress = defaults.ListenAddress
	}
	if h.ReadHeaderTimeout == 0 {
		h.ReadHeaderTimeout = defaults.ReadHeaderTimeout
	}
}

// Flags exposes configuration fields as flags. The current value of the
// receiver is used as the default value of the flag so ApplyDefaults should
// be called before requesting flags.
func (h *HTTPServer) Flags() []cli.Flag {
	def := HTTPServerDefaults()
	flags := []cli.Flag{
		NewStringFlag(&cli.StringFlag{
			Name:        "http-listen-address",
			Value:       h.ListenAddress,
			Destination: &h.ListenAddress,
			Usage:       "FIXME: http-listen-address",
			DefaultText: def.ListenAddress,
		}),
		NewDurationFlag(&cli.DurationFlag{
			Name:        "http-read-timeout",
			Value:       h.ReadTimeout,
			Destination: &h.ReadTimeout,
			Usage:       "FIXME: http-read-timeout",
			DefaultText: def.ReadTimeout.String(),
		}),
		NewDurationFlag(&cli.DurationFlag{
			Name:        "http-read-header-timeout",
			Value:       h.ReadHeaderTimeout,
			Destination: &h.ReadHeaderTimeout,
			Usage:       "FIXME: http-read-header-timeout",
			DefaultText: def.ReadHeaderTimeout.String(),
		}),
		NewDurationFlag(&cli.DurationFlag{
			Name:        "http-write-timeout",
			Value:       h.WriteTimeout,
			Destination: &h.WriteTimeout,
			Usage:       "FIXME: http-write-timeout",
			DefaultText: def.WriteTimeout.String(),
		}),
		NewDurationFlag(&cli.DurationFlag{
			Name:        "http-idle-timeout",
			Value:       h.IdleTimeout,
			Destination: &h.IdleTimeout,
			Usage:       "FIXME: http-idle-timeout",
			DefaultText: def.IdleTimeout.String(),
		}),
	}
	return flags
}

// BuildServer creates a basic *http.Server instance that is populated from the
// associated configuration. Additional customization can be performed by the
// caller.
func (h HTTPServer) BuildServer(tc *tls.Config) *http.Server {
	return &http.Server{
		Addr:              h.ListenAddress,
		TLSConfig:         tc,
		ReadTimeout:       h.ReadTimeout,
		ReadHeaderTimeout: h.ReadHeaderTimeout,
		WriteTimeout:      h.WriteTimeout,
		IdleTimeout:       h.IdleTimeout,
	}
}
