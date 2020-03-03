// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"os"

	"github.com/sykesm/batik/app"
)

func main() {
	app := app.Batik(os.Args, os.Stdin, os.Stdout)

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "command failed: %s", err)
		os.Exit(2)
	}
}
