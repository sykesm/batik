// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	cli "github.com/urfave/cli/v2"
	"github.com/sykesm/batik/app"
	"github.com/sykesm/batik/pkg/buildinfo"
)

func main() {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("Version:    %s\n", c.App.Version)
		fmt.Printf("Go Version: %s\n", runtime.Version())
		fmt.Printf("Git Commit: %s\n", buildinfo.GitCommit)
		fmt.Printf("OS/Arch:    %s\n", runtime.GOARCH)
		fmt.Printf("Built:      %s\n", c.App.Compiled.Format(time.ANSIC))
	}

	app := app.Batik(os.Args, os.Stdin, os.Stdout, os.Stderr)

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "command failed: %s", err)
		os.Exit(2)
	}
}
