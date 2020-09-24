// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"fmt"

	cli "github.com/urfave/cli/v2"
	"github.com/sykesm/batik/pkg/options"
	"github.com/sykesm/batik/pkg/store"
)

func dbCommand(config *options.Batik) *cli.Command {
	return &cli.Command{
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

					val, err := db.Get(key)
					if err != nil {
						fmt.Fprintf(ctx.App.ErrWriter, "%s\n", err)
					} else {
						fmt.Fprintf(ctx.App.ErrWriter, "%s\n", val)
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

					if err := db.Put(key, val); err != nil {
						fmt.Fprintf(ctx.App.ErrWriter, "%s\n", err)
					}

					return nil
				},
			},
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
