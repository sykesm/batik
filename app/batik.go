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
	zaplogfmt "github.com/sykesm/zap-logfmt"
	cli "github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/sykesm/batik/pkg/atexit"
	"github.com/sykesm/batik/pkg/buildinfo"
	"github.com/sykesm/batik/pkg/conf"
	"github.com/sykesm/batik/pkg/log"
	"github.com/sykesm/batik/pkg/log/pretty"
	"github.com/sykesm/batik/pkg/options"
	"github.com/sykesm/batik/pkg/repl"
)

func Batik(args []string, stdin io.ReadCloser, stdout, stderr io.Writer) *cli.App {
	atexit := atexit.New()
	config := options.BatikDefaults()

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
			Usage:   "Path to yaml configuration file",
			EnvVars: []string{"BATIK_CFG_PATH"},
		},
	}
	app.Flags = append(app.Flags, config.Logging.Flags()...)
	app.Commands = []*cli.Command{
		startCommand(config, false),
	}

	// Sort the flags and commands to make it easier to find things.
	// https://github.com/urfave/cli/blob/master/docs/v2/manual.md#ordering
	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	app.Before = func(ctx *cli.Context) error {
		err := resolveConfig(ctx, config)
		if err != nil {
			return cli.Exit(errors.WithMessage(err, "unable to read config"), exitConfigLoadFailed)
		}

		encoder, writer, leveler := newBatikLoggerComponents(ctx, config.Logging)
		logger := log.NewLogger(encoder, writer, leveler).Named(ctx.App.Name)

		atexit.Register(func() { logger.Sync() })
		SetLogger(ctx, logger)

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

	return nil
}

func newBatikLoggerComponents(ctx *cli.Context, config options.Logging) (zapcore.Encoder, zapcore.WriteSyncer, zapcore.LevelEnabler) {
	var encoder zapcore.Encoder
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.EpochNanosTimeEncoder

	w := ctx.App.ErrWriter
	f, ok := w.(*os.File)
	switch {
	case config.Color == "yes", config.Color == "auto" && ok && terminal.IsTerminal(int(f.Fd())):
		w = pretty.NewWriter(w, encoderConfig, pretty.ParseUnixTime)
		encoder = zaplogfmt.NewEncoder(encoderConfig)
	case config.Format == "json":
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	default:
		encoder = zaplogfmt.NewEncoder(encoderConfig)
	}

	return encoder, log.NewWriteSyncer(w), log.NewLeveler(config.LogSpec)
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
		startCommand(config, true),
		exitCommand(),
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

func exitCommand() *cli.Command {
	return &cli.Command{
		Name:  "exit",
		Usage: "exit the shell",
		Action: func(ctx *cli.Context) error {
			return repl.ErrExit
		},
	}
}
