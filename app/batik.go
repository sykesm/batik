// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"fmt"
	"io"
	"sort"

	"github.com/sykesm/batik/pkg/buildinfo"
	cli "github.com/urfave/cli/v2"
)

func Batik(args []string, stdin io.ReadCloser, stdout io.Writer) *cli.App {
	app := cli.NewApp()
	app.Copyright = fmt.Sprintf("© Copyright IBM Corporation %04d. All rights reserved.", buildinfo.Built().Year())
	app.Usage = "track some assets on the ledger"
	app.Compiled = buildinfo.Built()
	app.Version = buildinfo.FullVersion()
	app.EnableBashCompletion = true

	// setup flags for the ledger

	// Sort the flags and commands to make it easier to find things.
	// https://github.com/urfave/cli/blob/master/docs/v2/manual.md#ordering
	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	return app
}
