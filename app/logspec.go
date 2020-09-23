// Copyright IBM Corp. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package app

import (
	"errors"

	cli "github.com/urfave/cli/v2"
	"go.uber.org/zap"

	"github.com/sykesm/batik/pkg/log"
)

func logspecCommand() *cli.Command {
	return &cli.Command{
		Name:        "logspec",
		Description: "Dynamically change the log level of the enabled logger.",
		Usage:       "change the logspec of the logger leveler to any supported log level (eg. debug, info)",
		Action: func(ctx *cli.Context) error {
			leveler, err := GetLeveler(ctx)
			if err != nil {
				return cli.Exit(err, exitChangeLogspecFailed)
			}

			if atomicLevel, ok := leveler.(zap.AtomicLevel); ok {
				atomicLevel.SetLevel(log.NameToLevel(ctx.Args().Get(0)))
				return nil
			}

			return cli.Exit(errors.New("LevelEnabler is not an AtomicLevel"), exitChangeLogspecFailed)
		},
	}
}
