// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/pkg/errors"
	cli "github.com/urfave/cli/v2"
	yaml "gopkg.in/yaml.v3"

	"github.com/sykesm/batik/app/options"
	"github.com/sykesm/batik/pkg/buildinfo"
	"github.com/sykesm/batik/pkg/log"
	"github.com/sykesm/batik/pkg/repl"
)

func Batik(args []string, stdin io.ReadCloser, stdout, stderr io.Writer) *cli.App {
	config := options.ConfigDefaults()

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

	app.Before = func(ctx *cli.Context) error {
		logLevel := ctx.String("log-level")
		logger, err := log.NewLogger(log.Config{
			Name:    "batik",
			LogSpec: logLevel,
			Writer:  app.ErrWriter,
			Format:  "logfmt",
		})
		if err != nil {
			return cli.Exit(errors.Wrap(err, "failed creating new logger"), exitLoggerCreateFailed)
		}

		configPath := ctx.String("config")
		if configPath != "" {
			cf, err := os.Open(configPath)
			if err != nil {
				return cli.Exit(errors.Wrap(err, "unable to read config"), exitConfigLoadFailed)
			}

			decoder := yaml.NewDecoder(cf)
			err = decoder.Decode(config)
			if err != nil {
				return cli.Exit(errors.Wrap(err, "unable to read config"), exitConfigDecodeFailed)
			}
		}

		err = config.ApplyDefaults()
		if err != nil {
			panic(err)
		}

		SetLogger(ctx, logger)

		return nil
	}

	app.After = func(ctx *cli.Context) error {
		logger, _ := GetLogger(ctx)
		if logger != nil {
			logger.Sync()
		}

		return nil
	}

	// setup flags for the ledger
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

	// Sort the flags and commands to make it easier to find things.
	// https://github.com/urfave/cli/blob/master/docs/v2/manual.md#ordering
	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	return app
}

func shellApp(parentCtx *cli.Context, config *options.Config) (*cli.App, error) {
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
