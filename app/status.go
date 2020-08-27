// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"fmt"

	cli "github.com/urfave/cli/v2"
)

func statusCommand() *cli.Command {
	return &cli.Command{
		Name:        "status",
		Description: "check status of server",
		Action: func(ctx *cli.Context) error {
			server := ctx.App.Metadata["server"].(*BatikServer)

			if err := server.Status(); err != nil {
				return cli.Exit(fmt.Sprintf("Server not running at %s", server.address), 1)

			}
			return cli.Exit("Server running", 0)
		},
	}
}
