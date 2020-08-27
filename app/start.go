// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package app

import (
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
			server := ctx.App.Metadata["server"].(*BatikServer)
			if server == nil {
				return cli.Exit("server does not exist", 2)
			}
			if ctx.String("address") != "" {
				server.address = ctx.String("address")
			}
			if err := server.Start(); err != nil {
				return cli.Exit(err.Error(), 2)
			}

			return cli.Exit("Server stopped", 0)
		},
	}
}
