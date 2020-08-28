// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"fmt"

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
			if ctx.String("address") != "" {
				config.Server.Address = ctx.String("address")
			}

			server, err := NewServer(config, ctx.App.Writer, ctx.App.ErrWriter)
			if err != nil {
				return cli.Exit(fmt.Sprintf("failed to create server: %s", err), exitServerCreateFailed)
			}
			ctx.App.Metadata["server"] = server

			start := func() error {
				if err := server.Start(); err != nil {
					return cli.Exit(err.Error(), exitServerStartFailed)
				}

				fmt.Fprintln(ctx.App.ErrWriter, "Server stopped")
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
