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
			server := ctx.App.Metadata["server"]
			if server == nil {
				return cli.Exit("Server not running", 0)
			}

			bs := server.(*BatikServer)
			if err := bs.Status(); err != nil {
				return cli.Exit(fmt.Sprintf("Server not responding at %s", bs.address), 1)
			}

			return cli.Exit("Server running", 0)
		},
	}
}
