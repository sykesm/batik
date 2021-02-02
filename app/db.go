// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"encoding/hex"
	"fmt"
	"sort"
	"strconv"

	"github.com/hokaccha/go-prettyjson"
	cli "github.com/urfave/cli/v2"
	"github.com/sykesm/batik/pkg/options"
	"github.com/sykesm/batik/pkg/transaction"
)

func dbCommand(config *options.Batik) *cli.Command {

	command := &cli.Command{
		Name:  "db",
		Usage: "perform operations against a kv store",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "namespace",
				Usage:    "target namespace for the subcommand",
				Required: true,
			},
		},
		Subcommands: []*cli.Command{
			getSubcommand(),
			keysSubcommand(),
			putSubcommand(),
		},
	}

	sort.Sort(cli.CommandsByName(command.Subcommands))

	return command
}

func getSubcommand() *cli.Command {
	command := &cli.Command{
		Name:  "get",
		Usage: "get a value from the db",
		Subcommands: []*cli.Command{
			getTransactionSubcommand(),
			getStateSubcommand(),
		},
		Action: func(ctx *cli.Context) error {
			ns, err := GetCurrentNamespace(ctx)
			if err != nil {
				fmt.Fprintln(ctx.App.ErrWriter, err)
				return nil
			}

			key, err := hex.DecodeString(ctx.Args().First())
			if err != nil {
				fmt.Fprintln(ctx.App.ErrWriter, err)
				return nil
			}

			val, err := ns.LevelDB.Get(key)
			if err != nil {
				fmt.Fprintln(ctx.App.ErrWriter, err)
				return nil
			}

			fmt.Fprintln(ctx.App.ErrWriter, hex.Dump(val))
			return nil
		},
	}

	sort.Sort(cli.CommandsByName(command.Subcommands))

	return command
}

func getTransactionSubcommand() *cli.Command {
	return &cli.Command{
		Name:  "tx",
		Usage: "get a transaction from the db",
		Action: func(ctx *cli.Context) error {
			ns, err := GetCurrentNamespace(ctx)
			if err != nil {
				fmt.Fprintln(ctx.App.ErrWriter, err)
				return nil
			}

			txID, err := hex.DecodeString(ctx.Args().First())
			if err != nil {
				fmt.Fprintln(ctx.App.ErrWriter, err)
				return nil
			}

			val, err := ns.TxRepo.GetTransaction(txID)
			if err != nil {
				fmt.Fprintln(ctx.App.ErrWriter, err)
				return nil
			}

			jsonOut, err := prettyjson.Marshal(val)
			if err != nil {
				fmt.Fprintln(ctx.App.ErrWriter, err)
				return nil
			}

			fmt.Fprintln(ctx.App.ErrWriter, string(jsonOut))
			return nil
		},
	}
}

func getStateSubcommand() *cli.Command {
	return &cli.Command{
		Name:  "state",
		Usage: "get a state from the db",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "consumed",
				Usage: "fetch a consumed state",
			},
		},
		Action: func(ctx *cli.Context) error {
			ns, err := GetCurrentNamespace(ctx)
			if err != nil {
				fmt.Fprintln(ctx.App.ErrWriter, err)
				return nil
			}

			txID, err := hex.DecodeString(ctx.Args().First())
			if err != nil {
				fmt.Fprintln(ctx.App.ErrWriter, err)
				return nil
			}
			outputIndex, err := strconv.ParseUint(ctx.Args().Get(1), 10, 64)
			if err != nil {
				fmt.Fprintln(ctx.App.ErrWriter, err)
				return nil
			}

			stateID := transaction.StateID{
				TxID:        txID,
				OutputIndex: outputIndex,
			}
			val, err := ns.TxRepo.GetState(stateID, ctx.Bool("consumed"))
			if err != nil {
				fmt.Fprintln(ctx.App.ErrWriter, err)
				return nil
			}

			jsonOut, err := prettyjson.Marshal(val)
			if err != nil {
				fmt.Fprintln(ctx.App.ErrWriter, err)
				return nil
			}

			fmt.Fprintln(ctx.App.ErrWriter, string(jsonOut))
			return nil
		},
	}
}

func keysSubcommand() *cli.Command {
	return &cli.Command{
		Name:  "keys",
		Usage: "dump all keys in the db",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "prefix",
				Usage:       "prefix to range over",
				DefaultText: "",
			},
		},
		Action: func(ctx *cli.Context) error {
			ns, err := GetCurrentNamespace(ctx)
			if err != nil {
				fmt.Fprintln(ctx.App.ErrWriter, err)
				return nil
			}

			prefix, err := hex.DecodeString(ctx.String("prefix"))
			if err != nil {
				fmt.Fprintln(ctx.App.ErrWriter, err)
				return nil
			}

			iter := ns.LevelDB.NewIterator(prefix, nil)
			keys, err := iter.Keys()
			if err != nil {
				fmt.Fprintln(ctx.App.ErrWriter, err)
				return nil
			}

			for _, k := range keys {
				fmt.Fprintln(ctx.App.ErrWriter, hex.EncodeToString(k))
			}
			return nil
		},
	}
}

func putSubcommand() *cli.Command {
	return &cli.Command{
		Name:  "put",
		Usage: "store a value in the db",
		Action: func(ctx *cli.Context) error {
			ns, err := GetCurrentNamespace(ctx)
			if err != nil {
				fmt.Fprintln(ctx.App.ErrWriter, err)
				return nil
			}

			key, err := hex.DecodeString(ctx.Args().Get(0))
			if err != nil {
				fmt.Fprintln(ctx.App.ErrWriter, err)
				return nil
			}
			val, err := hex.DecodeString(ctx.Args().Get(1))
			if err != nil {
				fmt.Fprintln(ctx.App.ErrWriter, err)
				return nil
			}

			if err := ns.LevelDB.Put(key, val); err != nil {
				fmt.Fprintln(ctx.App.ErrWriter, err)
			}

			return nil
		},
	}
}
