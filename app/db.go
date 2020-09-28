// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"encoding/hex"
	"fmt"
	"sort"

	cli "github.com/urfave/cli/v2"
	"github.com/sykesm/batik/pkg/options"
	"github.com/sykesm/batik/pkg/store"
)

func dbCommand(config *options.Batik) *cli.Command {
	command := &cli.Command{
		Name:  "db",
		Usage: "perform operations against a kv store",
		Subcommands: []*cli.Command{
			{
				Name:  "get",
				Usage: "get a value from the db",
				Action: func(ctx *cli.Context) error {
					key := []byte(ctx.Args().First())

					db, err := levelDB(ctx, config.Ledger.DataDir)
					if err != nil {
						fmt.Fprintf(ctx.App.ErrWriter, "%s\n", err)
					}

					decodedKey := make([]byte, hex.DecodedLen(len(key)))
					if _, err := hex.Decode(decodedKey, key); err != nil {
						fmt.Fprintf(ctx.App.ErrWriter, "%s\n", err)
					}

					val, err := db.Get(decodedKey)
					if err != nil {
						fmt.Fprintf(ctx.App.ErrWriter, "%s\n", err)
					} else {
						fmt.Fprintf(ctx.App.ErrWriter, "%s\n", hex.Dump(val))
					}

					return nil
				},
			},
			{
				Name:  "put",
				Usage: "store a value in the db",
				Action: func(ctx *cli.Context) error {
					key := []byte(ctx.Args().Get(0))
					val := []byte(ctx.Args().Get(1))

					db, err := levelDB(ctx, config.Ledger.DataDir)
					if err != nil {
						fmt.Fprintf(ctx.App.ErrWriter, "%s\n", err)
					}

					decodedKey := make([]byte, hex.DecodedLen(len(key)))
					if _, err := hex.Decode(decodedKey, key); err != nil {
						fmt.Fprintf(ctx.App.ErrWriter, "%s\n", err)
					}
					decodedVal := make([]byte, hex.DecodedLen(len(val)))
					if _, err := hex.Decode(decodedVal, val); err != nil {
						fmt.Fprintf(ctx.App.ErrWriter, "%s\n", err)
					}

					if err := db.Put(decodedKey, decodedVal); err != nil {
						fmt.Fprintf(ctx.App.ErrWriter, "%s\n", err)
					}

					return nil
				},
			},
			{
				Name:  "keys",
				Usage: "dump all keys in the db",
				Action: func(ctx *cli.Context) error {
					db, err := levelDB(ctx, config.Ledger.DataDir)
					if err != nil {
						fmt.Fprintf(ctx.App.ErrWriter, "%s\n", err)
					}

					iter := db.NewIterator(nil, nil)
					keys, err := iter.Keys()
					if err != nil {
						fmt.Fprintf(ctx.App.ErrWriter, "%s\n", err)
					} else {
						for _, k := range keys {
							fmt.Fprintf(ctx.App.ErrWriter, "%s\n", hex.EncodeToString(k))
						}
					}

					return nil
				},
			},
		},
	}

	sort.Sort(cli.CommandsByName(command.Subcommands))

	return command
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
