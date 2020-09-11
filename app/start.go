// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"github.com/pkg/errors"
	cli "github.com/urfave/cli/v2"

	"github.com/sykesm/batik/app/options"
)

func startCommand(config *options.Config, interactive bool) *cli.Command {
	return &cli.Command{
		Name:        "start",
		Description: "start the grpc server",
		Flags: append(
			config.Server.Flags("start"),
			config.Ledger.Flags("start")...,
		),
		Action: func(ctx *cli.Context) error {
			logger, err := GetLogger(ctx)
			if err != nil {
				return cli.Exit(err, exitServerStartFailed)
			}

			serverConfig := Config{
				Server: Server{Address: config.Server.ListenAddress},
				DBPath: config.Ledger.DataDir,
			}
			server, err := NewServer(serverConfig, logger)
			if err != nil {
				return cli.Exit(errors.Wrap(err, "failed to create server"), exitServerCreateFailed)
			}

			SetServer(ctx, server)

			start := func() error {
				if err := server.Start(); err != nil {
					return cli.Exit(err, exitServerStartFailed)
				}

				logger.Info("Server stopped")
				return nil
			}

			// TODO: The server needs work to enable us to start, wait for ready, and then
			// return.
			if interactive {
				go start()
				return nil
			}
			return start()
		},
	}
}
