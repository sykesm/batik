// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/pkg/errors"
	cli "github.com/urfave/cli/v2"

	"github.com/sykesm/batik/pkg/atexit"
	"github.com/sykesm/batik/pkg/buildinfo"
	"github.com/sykesm/batik/pkg/conf"
	"github.com/sykesm/batik/pkg/log"
	"github.com/sykesm/batik/pkg/options"
	"github.com/sykesm/batik/pkg/repl"
)

func Batik(args []string, stdin io.ReadCloser, stdout, stderr io.Writer) *cli.App {
	config := options.BatikDefaults()
	atexit := atexit.New()

	app := cli.NewApp()
	app.Copyright = fmt.Sprintf("© Copyright IBM Corporation %04d. All rights reserved.", buildinfo.Built().Year())
	app.Name = "batik"
	app.Usage = "track some assets on the ledger"
	app.Compiled = buildinfo.Built()
	app.Version = buildinfo.FullVersion()
	app.Writer = stdout
	app.ErrWriter = stderr
	app.EnableBashCompletion = true
	app.CommandNotFound = func(ctx *cli.Context, name string) {
		fmt.Fprintf(ctx.App.ErrWriter, "%[1]s: '%[2]s' is not a %[1]s command. See `%[1]s --help`.\n", ctx.App.Name, name)
	}
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "config",
			Aliases: []string{"c"},
			Usage:   "Path to yaml file to load configuration parameters from",
			EnvVars: []string{"BATIK_CFG_PATH"},
		},
		&cli.StringFlag{
			Name:    "log-level",
			Usage:   "Log level",
			Value:   "info",
			EnvVars: []string{"BATIK_LOG_LEVEL"},
		},
	}
	app.Commands = []*cli.Command{
		startCommand(config, false),
		statusCommand(),
	}

	// Sort the flags and commands to make it easier to find things.
	// https://github.com/urfave/cli/blob/master/docs/v2/manual.md#ordering
	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	app.Before = func(ctx *cli.Context) error {
		logLevel := ctx.String("log-level")
		logger, err := log.NewLogger(log.Config{
			Name:    app.Name,
			LogSpec: logLevel,
			Writer:  app.ErrWriter,
			Format:  "logfmt",
		})
		if err != nil {
			return cli.Exit(errors.Wrap(err, "failed creating new logger"), exitLoggerCreateFailed)
		}

		atexit.Register(func() { logger.Sync() })
		SetLogger(ctx, logger)

		err = resolveConfig(ctx, config)
		if err != nil {
			return cli.Exit(errors.WithMessage(err, "unable to read config"), exitConfigLoadFailed)
		}

		return nil
	}

	app.After = func(ctx *cli.Context) error {
		defer atexit.Exit()
		return nil
	}

	// The default action starts the interactive console shell.
	app.Action = func(ctx *cli.Context) error {
		if ctx.Args().Present() {
			arg := ctx.Args().First()
			ctx.App.CommandNotFound(ctx, arg)
			return cli.Exit("", exitCommandNotFound)
		}

		sa, err := shellApp(ctx, config)
		if err != nil {
			return cli.Exit(err, exitShellSetupFailed)
		}

		return repl.New(sa, repl.WithStdin(stdin), repl.WithStdout(stdout), repl.WithStderr(stderr)).Run(ctx.Context)
	}

	return app
}

func resolveConfig(ctx *cli.Context, config *options.Batik) error {
	configPath := ctx.String("config")
	if configPath == "" {
		cf, err := conf.File(ctx.App.Name)
		if err != nil {
			return err
		}
		configPath = cf
	}

	if configPath != "" {
		err := conf.LoadFile(configPath, config)
		if err != nil {
			return err
		}
	}

	return config.ApplyDefaults()
}

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
		{
			Name:        "exit",
			Description: "exit the shell",
			Action: func(ctx *cli.Context) error {
				return repl.ErrExit
			},
		},
		startCommand(config, true),
		statusCommand(),
	}

	sort.Sort(cli.CommandsByName(app.Commands))

	// Generate the help message
	s := strings.Builder{}
	s.WriteString("Commands:\n")
	w := tabwriter.NewWriter(&s, 0, 0, 1, ' ', 0)
	for _, c := range app.VisibleCommands() {
		if _, err := fmt.Fprintf(w, "    %s %s\t%s\n", c.Name, c.Usage, c.Description); err != nil {
			return nil, err
		}
	}
	if err := w.Flush(); err != nil {
		return nil, err
	}

	app.CustomAppHelpTemplate = s.String()

	return app, nil
}
