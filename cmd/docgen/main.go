// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
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

	batik := app.Batik(nil, ioutil.NopCloser(&bytes.Buffer{}), ioutil.Discard, ioutil.Discard)

	doctype := "markdown"
	if len(os.Args) >= 2 {
		doctype = os.Args[1]
	}

	var doc string
	var err error

	switch doctype {
	case "markdown", "md":
		doc, err = batik.ToMarkdown()
	case "man":
		doc, err = batik.ToMan()
	default:
		log.Fatalf("unknown doc type: %s", doctype)
	}
	if err != nil {
		log.Fatalf("docgen failed: %v", err)
	}

	fmt.Println(doc)
}
