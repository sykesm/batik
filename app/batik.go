// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

	"github.com/bytecodealliance/wasmtime-go"
	"github.com/pkg/errors"
	zaplogfmt "github.com/sykesm/zap-logfmt"
	cli "github.com/urfave/cli/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/crypto/ssh/terminal"
	"gopkg.in/yaml.v3"

	"github.com/sykesm/batik/pkg/atexit"
	"github.com/sykesm/batik/pkg/buildinfo"
	"github.com/sykesm/batik/pkg/conf"
	"github.com/sykesm/batik/pkg/log"
	"github.com/sykesm/batik/pkg/log/pretty"
	"github.com/sykesm/batik/pkg/namespace"
	"github.com/sykesm/batik/pkg/options"
	"github.com/sykesm/batik/pkg/repl"
	"github.com/sykesm/batik/pkg/store"
	"github.com/sykesm/batik/pkg/submit"
	"github.com/sykesm/batik/pkg/validator"
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
		&cli.BoolFlag{
			Name:   "show-config",
			Hidden: true,
			Usage:  "Show the yaml serialization of the active configuration",
		},
	}
	app.Flags = append(app.Flags, config.Flags()...)
	app.Flags = append(app.Flags, config.Logging.Flags()...)
	app.Commands = []*cli.Command{
		startCommand(config, false),
		dbCommand(config),
	}

	// Sort the flags and commands to make it easier to find things.
	// https://github.com/urfave/cli/blob/master/docs/v2/manual.md#ordering
	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	// Setup hooks to show configuration instead of running commands.
	for i, cmd := range app.Commands {
		app.Commands[i].Before = showConfigBefore(cmd.Before, config)
	}

	app.Before = func(ctx *cli.Context) error {
		err := resolveConfig(ctx, config)
		if err != nil {
			return cli.Exit(errors.WithMessage(err, "unable to read config"), exitConfigLoadFailed)
		}

		encoder, writer, leveler := newBatikLoggerComponents(ctx, config.Logging)
		logger := log.NewLogger(encoder, writer, leveler).Named(ctx.App.Name)

		atexit.Register(func() { logger.Sync() })
		SetLogger(ctx, logger)
		SetLeveler(ctx, leveler)

		validators, err := newBatikValidatorComponents(config.Validators)
		if err != nil {
			return cli.Exit(err, exitConfigLoadFailed)
		}

		namespaces, err := newBatikNamespaceComponents(ctx, config.Namespaces, validators)
		if err != nil {
			return cli.Exit(err, exitConfigLoadFailed)
		}

		SetNamespaces(ctx, namespaces)
		// TODO safely shut down the DB
		// atexit.Register(func() { namespaces.Close() })

		return nil
	}

	app.After = func(ctx *cli.Context) error {
		defer atexit.Exit()
		return nil
	}

	// The default action starts the interactive console shell.
	app.Action = func(ctx *cli.Context) error {
		// Handle --show-config without a subcommand
		if err := showConfigBefore(nil, config)(ctx); err != nil {
			return err
		}

		// User provided a command that was not found
		if ctx.Args().Present() {
			arg := ctx.Args().First()
			ctx.App.CommandNotFound(ctx, arg)
			return cli.Exit("", exitCommandNotFound)
		}

		// Interactive shell
		sa, err := shellApp(ctx, config)
		if err != nil {
			return cli.Exit(err, exitShellSetupFailed)
		}

		return repl.New(sa, repl.WithStdin(stdin), repl.WithStdout(stdout), repl.WithStderr(stderr)).Run(ctx.Context)
	}

	return app
}

// resolveConfig is another hack related to the single configuration object
// pattern we're using.
//
// Top-level flags are used to point to the location of the configuration file
// and to setup things like logging and the data directory. When a value is
// provided with a flag, it is supposed to take priority over values loaded
// from the confguration file. The configuration file can't be loaded until
// after the flags are processed (because `--config` is a flag) and loading the
// configuration file will overwite the values set by the flags...
//
// So, the hack:
//   1. Save the values set by the local flags
//   2. Load the configuration file
//   3. Reapply the values to the flags
//   4. Rerun tag resolution on the configuration file
func resolveConfig(ctx *cli.Context, config *options.Batik) error {
	// Save the top-level flags to push back to the config
	flags := map[string]interface{}{}
	for _, name := range ctx.LocalFlagNames() {
		flags[name] = ctx.Value(name)
	}

	// Find the configuration file to use
	configPath := ctx.String("config")
	if configPath == "" {
		cf, err := conf.File(ctx.App.Name)
		if err != nil {
			return err
		}
		configPath = cf
	}

	// Load the configuration file.
	if configPath != "" {
		err := conf.LoadFile(configPath, config)
		if err != nil {
			return err
		}
	}

	// Re-apply the top level flags to the config as they take priority over
	// values loaded from the configuration file.
	for name, val := range flags {
		ctx.Set(name, fmt.Sprintf("%v", val))
	}

	// Resolve the tags after re-applying flags
	tr := &conf.TagResolver{SourcePath: filepath.Dir(configPath)}
	return tr.Resolve(config)
}

// showConfigBefore returns a function that dumps the configuration and
// terminates the command when the `show-config` flag is used. If the flag is
// not set, cli.BeforeFunc `b` is called instead.
//
// This is intented to be a test utility and debug aid.
func showConfigBefore(b cli.BeforeFunc, config *options.Batik) cli.BeforeFunc {
	return func(ctx *cli.Context) error {
		if ctx.Bool("show-config") {
			if err := yaml.NewEncoder(ctx.App.Writer).Encode(config); err != nil {
				return cli.Exit(err, exitConfigEncodeFailed)
			}
			return cli.Exit("", exitOkay)
		}
		if b != nil {
			return b(ctx)
		}
		return nil
	}
}

func newBatikLoggerComponents(ctx *cli.Context, config options.Logging) (zapcore.Encoder, zapcore.WriteSyncer, zap.AtomicLevel) {
	var encoder zapcore.Encoder
	encoderConfig := zap.NewProductionEncoderConfig()

	w := ctx.App.ErrWriter
	f, ok := w.(*os.File)
	switch {
	case config.Format == "json":
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	case config.Color == "yes", config.Color == "auto" && ok && terminal.IsTerminal(int(f.Fd())):
		w = pretty.NewWriter(w, encoderConfig, pretty.ParseUnixTime)
		encoder = zaplogfmt.NewEncoder(encoderConfig)
	default:
		encoder = zaplogfmt.NewEncoder(encoderConfig)
	}

	return encoder, log.NewWriteSyncer(w), log.NewLeveler(config.LogSpec)
}

func newBatikNamespaceComponents(ctx *cli.Context, config []options.Namespace, validators map[string]submit.Validator) (map[string]*namespace.Namespace, error) {
	logger, err := GetLogger(ctx)
	if err != nil {
		return nil, errors.WithMessage(err, "could not retrieve logger")
	}

	namespaces := map[string]*namespace.Namespace{}
	for _, ns := range config {
		namespaceLogger := logger.With(zap.String("namespace", ns.Name))

		dbPath := filepath.Join(ns.DataDir, ns.Name)
		namespaceLogger.Debug("initializing database", zap.String("data_dir", dbPath))
		db, err := store.NewLevelDB(dbPath)
		if err != nil {
			return nil, err
		}

		v, ok := validators[ns.Validator]
		if !ok {
			return nil, errors.Errorf("namespace %q requires validator %q which is not defined", ns.Name, ns.Validator)
		}

		namespaces[ns.Name] = namespace.New(namespaceLogger, db, v)
	}
	return namespaces, nil
}

func newBatikValidatorComponents(config []options.Validator) (map[string]submit.Validator, error) {
	var wasmEngine *wasmtime.Engine
	result := map[string]submit.Validator{}
	for i, validatorConf := range config {
		if validatorConf.Name == "" {
			return nil, errors.Errorf("validator at position %d in config has no name", i)
		}
		var v submit.Validator
		switch validatorConf.Type {
		case "builtin":
			if validatorConf.Name != "signature-builtin" {
				return nil, errors.Errorf("validator %q is not a known builtin validator", validatorConf.Name)
			}
			v = validator.NewSignature()
		case "wasm":
			wasmBin, err := ioutil.ReadFile(validatorConf.Path)
			if err != nil {
				return nil, errors.Wrapf(err, "could not load wasm binary for validator %q at %q", validatorConf.Name, validatorConf.Path)
			}

			if wasmEngine == nil {
				wasmEngine = wasmtime.NewEngine()
			}

			v, err = validator.NewWASM(wasmEngine, wasmBin)
			if err != nil {
				return nil, errors.WithMessagef(err, "could not create wasm validator for %q", validatorConf.Name)
			}
		default:
			return nil, errors.Errorf("validator %q has unknown type %q, must be wasm or builtin", validatorConf.Name, validatorConf.Type)
		}

		result[validatorConf.Name] = v
	}

	return result, nil
}
