// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"text/tabwriter"

	cli "github.com/urfave/cli/v2"

	"github.com/sykesm/batik/pkg/buildinfo"
	"github.com/sykesm/batik/pkg/config"
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
	app.Commands = []*cli.Command{
		startCommand(),
		statusCommand(),
	}

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "config",
			Aliases: []string{"c"},
			Usage:   "Path to yaml file to load configuration parameters from",
			EnvVars: []string{"BATIK_CFG_PATH"},
		},
	}

	app.Before = func(c *cli.Context) error {
		// Load config file
		cfgPath := c.String("config")

		var cfg Config
		err := config.Load(cfgPath, config.EnvironLookuper(), &cfg)
		if err != nil {
			return cli.Exit(fmt.Sprintf("failed loading batik config: %s", err), 3)
		}

		server, err := NewServer(cfg)
		if err != nil {
			return cli.Exit(fmt.Sprintf("failed to create server: %s", err), 2)
		}

		app.Metadata = map[string]interface{}{
			"config": cfg,
			"server": server,
		}

		return nil
	}

	// setup flags for the ledger
	app.Action = func(c *cli.Context) error {
		if c.Args().Present() {
			arg := c.Args().First()
			if c.App.CommandNotFound == nil {
				return cli.Exit(fmt.Sprintf("%[1]s: '%[2]s' is not a %[1]s command. See `%[1]s --help`.\n", c.App.Name, arg), 3)
			}
			c.App.CommandNotFound(c, arg)
			return cli.Exit("", 3)
		}

		server := c.App.Metadata["server"].(*BatikServer)
		sa, err := shellApp(server)
		if err != nil {
			return cli.Exit(err, 3)
		}
		repl := repl.New(sa)
		return repl.Run(c.Context)
	}

	// Sort the flags and commands to make it easier to find things.
	// https://github.com/urfave/cli/blob/master/docs/v2/manual.md#ordering
	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	return app
}

func shellApp(server *BatikServer) (*cli.App, error) {
	app := cli.NewApp()
	app.Name = "batik"
	app.HideVersion = true
	app.UsageText = "command [arguments...]"
	app.CommandNotFound = func(c *cli.Context, name string) {
		fmt.Fprintf(c.App.ErrWriter, "Unknown command: %s\n", name)
	}

	app.Commands = []*cli.Command{
		{
			Name:        "exit",
			Description: "exit the shell",
			Action: func(ctx *cli.Context) error {
				return cli.Exit(repl.ErrExit, 0)
			},
		},
		startCommand(),
		statusCommand(),
	}
	app.Metadata = map[string]interface{}{
		"server": server,
	}

	sort.Sort(cli.CommandsByName(app.Commands))

	// Generate the help message
	s := strings.Builder{}
	s.WriteString("Commands:\n")
	w := tabwriter.NewWriter(&s, 0, 0, 1, ' ', 0)

	for _, c := range app.VisibleCommands() {
		_, err := fmt.Fprintf(w,
			"    %s %s\t%s\n",
			c.Name, c.Usage,
			c.Description,
		)
		if err != nil {
			return nil, err
		}
	}

	if err := w.Flush(); err != nil {
		return nil, err
	}

	app.CustomAppHelpTemplate = s.String()

	return app, nil
}
