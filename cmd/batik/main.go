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
	cli.VersionPrinter = func(ctx *cli.Context) {
		fmt.Printf("Version:    %s\n", ctx.App.Version)
		fmt.Printf("Go Version: %s\n", runtime.Version())
		fmt.Printf("Git Commit: %s\n", buildinfo.GitCommit)
		fmt.Printf("OS/Arch:    %s\n", runtime.GOARCH)
		fmt.Printf("Built:      %s\n", ctx.App.Compiled.Format(time.ANSIC))
	}

	app.Batik(os.Args, os.Stdin, os.Stdout, os.Stderr).Run(os.Args)
}
