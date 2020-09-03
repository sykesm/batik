// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"github.com/pkg/errors"
	cli "github.com/urfave/cli/v2"
)

func startCommand(interactive bool) *cli.Command {
	return &cli.Command{
		Name:        "start",
		Description: "start the grpc server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "address",
				Aliases: []string{"a"},
				Usage:   "Listen address for the grpc server",
			},
		},
		Action: func(ctx *cli.Context) error {
			config := GetConfig(ctx)
			logger, err := GetLogger(ctx)
			if err != nil {
				return cli.Exit(err, exitServerStartFailed)
			}

			if ctx.String("address") != "" {
				config.Server.Address = ctx.String("address")
			}

			server, err := NewServer(config, logger)
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
