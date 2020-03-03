// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"io"

	"github.com/urfave/cli"
	"github.com/sykesm/batik/pkg/buildinfo"
)

func Batik(args []string, stdin io.ReadCloser, stdout io.Writer) *cli.App {
	app := cli.NewApp()
	app.Author = "IBM Corporation"
	app.Usage = "track some assets on the ledger"
	app.Compiled = buildinfo.Built()
	app.Version = buildinfo.FullVersion()
	app.EnableBashCompletion = true

	return app
}
