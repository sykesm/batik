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
	"github.com/sykesm/batik/pkg/config"
)

func main() {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("Version:    %s\n", c.App.Version)
		fmt.Printf("Go Version: %s\n", runtime.Version())
		fmt.Printf("Git Commit: %s\n", buildinfo.GitCommit)
		fmt.Printf("OS/Arch:    %s\n", runtime.GOARCH)
		fmt.Printf("Built:      %s\n", c.App.Compiled.Format(time.ANSIC))
	}

	cfgPath := os.Getenv("BATIK_CFG_DIR")
	cfg, err := config.NewBatikConfig(cfgPath)
	if err != nil {
		errorExit(2, "failed loading batik config: %s", err)
	}

	app := app.Batik(os.Args, os.Stdin, os.Stdout, os.Stderr, cfg)

	if err := app.Run(os.Args); err != nil {
		errorExit(2, "command failed: %s", err)
	}
}

func errorExit(exitCode int, fmtMsg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, fmtMsg+"\n", args...)
	os.Exit(exitCode)
}
