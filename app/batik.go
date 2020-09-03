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

	"github.com/sykesm/batik/pkg/buildinfo"
	"github.com/sykesm/batik/pkg/config"
	"github.com/sykesm/batik/pkg/log"
	"github.com/sykesm/batik/pkg/repl"
)

func Batik(args []string, stdin io.ReadCloser, stdout, stderr io.Writer) *cli.App {
	app := cli.NewApp()
	app.Copyright = fmt.Sprintf("© Copyright IBM Corporation %04d. All rights reserved.", buildinfo.Built().Year())
	app.Name = "batik"
	app.Usage = "track some assets on the ledger"
	app.Compiled = buildinfo.Built()
	app.Version = buildinfo.FullVersion()
	app.Writer = stdout
	app.ErrWriter = stderr
	app.EnableBashCompletion = true
	app.CommandNotFound = func(c *cli.Context, name string) {
		fmt.Fprintf(c.App.ErrWriter, "%[1]s: '%[2]s' is not a %[1]s command. See `%[1]s --help`.\n", c.App.Name, name)
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
		&cli.StringFlag{
			Name:  "log-output-file",
			Usage: "Log output file path, if stdout, logs will go to stdout",
			Value: "stdout",
		},
		&cli.StringFlag{
			Name:  "errlog-output-file",
			Usage: "ErrLog output file path, if stderr, err logs will go to stderr",
			Value: "stderr",
		},
	}
	app.Commands = []*cli.Command{
		startCommand(false),
		statusCommand(),
	}

	app.Metadata = make(map[string]interface{})

	app.Before = func(c *cli.Context) error {
		logLevel := c.String("log-level")
		// logPath := c.String("log-output-file")
		// errLogPath := c.String("errlog-output-file")

		// w, err := log.NewWriter(logPath)
		// if err != nil {
		// 	return cli.Exit(errors.Wrap(err, "failed creating log writer"), exitLoggerCreateFailed)
		// }
		logger, err := log.NewLogger(log.Config{
			Name:    "batik",
			LogSpec: logLevel,
			Writer:  app.ErrWriter,
			// Format:  "logfmt",
		})
		if err != nil {
			return cli.Exit(errors.Wrap(err, "failed creating new logger"), exitLoggerCreateFailed)
		}

		SetLogger(c, logger)

		configPath := c.String("config")

		var cfg Config
		if err := config.Load(configPath, config.EnvironLookuper(), &cfg); err != nil {
			return cli.Exit(errors.Wrap(err, "failed loading batik config"), exitConfigLoadFailed)
		}

		SetConfig(c, cfg)

		return nil
	}

	app.After = func(c *cli.Context) error {
		logger, err := GetLogger(c)
		if err != nil {
			return cli.Exit(errors.Wrap(err, "failed to retrieve logger"), exitAppShutdownFailed)
		}
		logger.Sync()

		return nil
	}

	// setup flags for the ledger
	app.Action = func(c *cli.Context) error {
		if c.Args().Present() {
			arg := c.Args().First()
			c.App.CommandNotFound(c, arg)
			return cli.Exit("", exitCommandNotFound)
		}

		sa, err := shellApp(c)
		if err != nil {
			return cli.Exit(err, exitShellSetupFailed)
		}

		return repl.New(sa, repl.WithStdin(stdin), repl.WithStdout(stdout), repl.WithStderr(stderr)).Run(c.Context)
	}

	// Sort the flags and commands to make it easier to find things.
	// https://github.com/urfave/cli/blob/master/docs/v2/manual.md#ordering
	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	return app
}

func shellApp(ctx *cli.Context) (*cli.App, error) {
	app := cli.NewApp()
	app.Name = "batik"
	app.HideVersion = true
	app.Writer = ctx.App.Writer
	app.ErrWriter = ctx.App.ErrWriter
	app.CommandNotFound = func(c *cli.Context, name string) {
		fmt.Fprintf(c.App.ErrWriter, "Unknown command: %s\n", name)
	}
	app.ExitErrHandler = func(c *cli.Context, err error) {}
	app.Metadata = make(map[string]interface{})

	app.Before = func(c *cli.Context) error {
		SetConfig(c, GetConfig(ctx))
		logger, err := GetLogger(ctx)
		if err != nil {
			fmt.Fprintf(c.App.ErrWriter, "Failed setup: %s\n", err)
		}
		SetLogger(c, logger)

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
		startCommand(true),
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
