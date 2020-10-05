// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"text/tabwriter"

	cli "github.com/urfave/cli/v2"
	"go.uber.org/zap/zapcore"

	"github.com/sykesm/batik/pkg/log"
	"github.com/sykesm/batik/pkg/options"
	"github.com/sykesm/batik/pkg/repl"
)

// shellApp is the interactive console application.
func shellApp(parentCtx *cli.Context, config *options.Batik) (*cli.App, error) {
	app := cli.NewApp()
	app.Name = "batik"
	app.HideVersion = true
	app.Writer = parentCtx.App.Writer
	app.ErrWriter = parentCtx.App.ErrWriter
	app.CommandNotFound = func(ctx *cli.Context, name string) {
		fmt.Fprintf(ctx.App.ErrWriter, "Unknown command: %s\n", name)
	}
	app.ExitErrHandler = func(ctx *cli.Context, err error) {}

	app.Before = func(ctx *cli.Context) error {
		logger, err := GetLogger(parentCtx)
		if err != nil {
			fmt.Fprintf(ctx.App.ErrWriter, "Failed setup: %s\n", err)
		}
		SetLogger(ctx, logger)

		return nil
	}

	app.Commands = []*cli.Command{
		dbCommand(config),
		exitCommand(),
		logspecCommand(),
		startCommand(config, true),
	}

	sort.Sort(cli.CommandsByName(app.Commands))

	// Generate the help message
	s := strings.Builder{}
	s.WriteString("Commands:\n")
	w := tabwriter.NewWriter(&s, 8, 8, 2, ' ', 0)
	for _, c := range app.VisibleCommands() {
		if _, err := fmt.Fprintf(w, "    %s\t%s\n", c.Name, c.Usage); err != nil {
			return nil, err
		}
	}
	if err := w.Flush(); err != nil {
		return nil, err
	}

	app.CustomAppHelpTemplate = s.String()

	return app, nil
}

func logspecCommand() *cli.Command {
	return &cli.Command{
		Name:        "logspec",
		Description: "Dynamically change the log level of the enabled logger.",
		Usage:       "change the logspec of the logger leveler to any supported log level (eg. debug, info)",
		Action: func(ctx *cli.Context) error {
			leveler, err := GetLeveler(ctx)
			if err != nil {
				return cli.Exit(err, exitChangeLogspecFailed)
			}
			type levelSetter interface{ SetLevel(zapcore.Level) }
			ls, ok := leveler.(levelSetter)
			if !ok {
				return cli.Exit(errors.New("log level cannot be changed"), exitChangeLogspecFailed)
			}
			ls.SetLevel(log.NameToLevel(ctx.Args().First()))
			return nil
		},
	}
}

func exitCommand() *cli.Command {
	return &cli.Command{
		Name:  "exit",
		Usage: "exit the shell",
		Action: func(ctx *cli.Context) error {
			return repl.ErrExit
		},
	}
}
