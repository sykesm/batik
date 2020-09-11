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
			// server, err := GetServer(ctx)
			// if err != nil {
			// 	return cli.Exit(err, exitServerStatusFailed)
			// }
			// if server == nil {
			fmt.Fprintln(ctx.App.Writer, "Server not running")
			return nil
			// }

			// if err := server.Status(); err != nil {
			// 	return cli.Exit(fmt.Errorf("server not responding at %s", server.address), exitServerStatusFailed)
			// }

			// fmt.Fprintln(ctx.App.Writer, "Server running")
			// return nil
		},
	}
}
