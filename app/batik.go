// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"fmt"
	"io"
	"net"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	cli "github.com/urfave/cli/v2"
	"github.com/sykesm/batik/pkg/buildinfo"
	tb "github.com/sykesm/batik/pkg/pb/transaction"
	"github.com/sykesm/batik/pkg/repl"
	"github.com/sykesm/batik/pkg/transaction"
	"google.golang.org/grpc"
)

var statusCommand = &cli.Command{
	Name:        "status",
	Description: "check status of server",
	Action: func(ctx *cli.Context) error {
		address := ctx.String("address")
		if err := checkStatus(address); err != nil {
			return cli.Exit(fmt.Sprintf("Server not running at %s", address), 1)

		}
		return cli.Exit("Server running", 0)
	},
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "address",
			Aliases: []string{"a"},
			Usage:   "Listen address for the grpc server",
		},
	},
}

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
		os.Exit(3)
	}
	app.Commands = []*cli.Command{
		{
			Name:        "start",
			Description: "start the grpc server",
			Action: func(ctx *cli.Context) error {
				address := ctx.String("address")
				if err := startServer(address); err != nil {
					return cli.Exit(err.Error(), 2)
				}

				return cli.Exit("Server started", 0)
			},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "address",
					Aliases:  []string{"a"},
					Usage:    "Listen address for the grpc server",
					Required: true,
				},
			},
		},
		statusCommand,
	}

	// setup flags for the ledger
	app.Action = func(c *cli.Context) error {
		if c.Args().Present() {
			arg := c.Args().First()
			if c.App.CommandNotFound != nil {
				c.App.CommandNotFound(c, arg)
			} else {
				return cli.Exit(fmt.Sprintf("%[1]s: '%[2]s' is not a %[1]s command. See `%[1]s --help`.\n", c.App.Name, arg), 3)
			}
		}

		sa, err := shellApp()
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

func shellApp() (*cli.App, error) {
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
		statusCommand,
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

func startServer(address string) error {
	fmt.Printf("Starting server at %s\n", address)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	server := grpc.NewServer()

	encodeTxSvc := &transaction.EncodeService{}
	tb.RegisterEncodeTransactionAPIServer(server, encodeTxSvc)

	return server.Serve(listener)
}

func checkStatus(address string) error {
	fmt.Printf("Checking status of server at %s\n", address)

	//create GRPC client conn
	clientConn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer clientConn.Close()

	//TODO add client healthcheck to verify grpc server status

	return nil
}
