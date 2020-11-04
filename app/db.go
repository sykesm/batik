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
	"github.com/sykesm/batik/pkg/store"
	"github.com/sykesm/batik/pkg/transaction"
)

func dbCommand(config *options.Batik) *cli.Command {
	command := &cli.Command{
		Name:  "db",
		Usage: "perform operations against a kv store",
		Subcommands: []*cli.Command{
			getSubcommand(config.Ledger.DataDir),
			keysSubcommand(config.Ledger.DataDir),
			putSubcommand(config.Ledger.DataDir),
		},
	}

	sort.Sort(cli.CommandsByName(command.Subcommands))

	return command
}

func getSubcommand(dataDir string) *cli.Command {
	command := &cli.Command{
		Name:  "get",
		Usage: "get a value from the db",
		Subcommands: []*cli.Command{
			getTransactionSubcommand(dataDir),
			getStateSubcommand(dataDir),
		},
		Action: func(ctx *cli.Context) error {
			db, err := levelDB(ctx, dataDir)
			if err != nil {
				fmt.Fprintln(ctx.App.ErrWriter, err)
				return nil
			}

			key, err := hex.DecodeString(ctx.Args().First())
			if err != nil {
				fmt.Fprintln(ctx.App.ErrWriter, err)
				return nil
			}

			val, err := db.Get(key)
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

func getTransactionSubcommand(dataDir string) *cli.Command {
	return &cli.Command{
		Name:  "tx",
		Usage: "get a transaction from the db",
		Action: func(ctx *cli.Context) error {
			db, err := levelDB(ctx, dataDir)
			if err != nil {
				fmt.Fprintln(ctx.App.ErrWriter, err)
				return nil
			}

			transactionRepo := store.NewRepository(db)

			txID, err := hex.DecodeString(ctx.Args().First())
			if err != nil {
				fmt.Fprintln(ctx.App.ErrWriter, err)
				return nil
			}

			val, err := transactionRepo.GetTransaction(txID)
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

func getStateSubcommand(dataDir string) *cli.Command {
	return &cli.Command{
		Name:  "state",
		Usage: "get a state from the db",
		Action: func(ctx *cli.Context) error {
			db, err := levelDB(ctx, dataDir)
			if err != nil {
				fmt.Fprintln(ctx.App.ErrWriter, err)
				return nil
			}

			transactionRepo := store.NewRepository(db)

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
			val, err := transactionRepo.GetState(stateID)
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

func keysSubcommand(dataDir string) *cli.Command {
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
			db, err := levelDB(ctx, dataDir)
			if err != nil {
				fmt.Fprintln(ctx.App.ErrWriter, err)
				return nil
			}

			prefix, err := hex.DecodeString(ctx.String("prefix"))
			if err != nil {
				fmt.Fprintln(ctx.App.ErrWriter, err)
				return nil
			}

			iter := db.NewIterator(prefix, nil)
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

func putSubcommand(dataDir string) *cli.Command {
	return &cli.Command{
		Name:  "put",
		Usage: "store a value in the db",
		Action: func(ctx *cli.Context) error {
			db, err := levelDB(ctx, dataDir)
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

			if err := db.Put(key, val); err != nil {
				fmt.Fprintln(ctx.App.ErrWriter, err)
			}

			return nil
		},
	}
}

func levelDB(ctx *cli.Context, dir string) (*store.LevelDBKV, error) {
	var err error
	db := GetKV(ctx)
	if db == nil {
		db, err = store.NewLevelDB(dir)
		if err != nil {
			return nil, err
		}
	}

	return db.(*store.LevelDBKV), nil
}
