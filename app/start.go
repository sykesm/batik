// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"fmt"

	cli "github.com/urfave/cli/v2"
)

func startCommand() *cli.Command {
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
			config := ctx.App.Metadata["config"].(Config)
			if ctx.String("address") != "" {
				config.Server.Address = ctx.String("address")
			}

			server, err := NewServer(config, ctx.App.Writer, ctx.App.ErrWriter)
			if err != nil {
				return cli.Exit(fmt.Sprintf("failed to create server: %s", err), exitServerCreateFailed)
			}

			if err := server.Start(); err != nil {
				return cli.Exit(err.Error(), exitServerStartFailed)
			}

			fmt.Fprintln(ctx.App.ErrWriter, "Server stopped")
			return nil
		},
	}
}
