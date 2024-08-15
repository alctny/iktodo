package cmd

import (
	"errors"

	"github.com/alctny/iktodo/common"
	"github.com/alctny/iktodo/db"
	"github.com/urfave/cli/v2"
)

func RemoveCmd() *cli.Command {
	return &cli.Command{
		Name:      "remove",
		Usage:     "remove task",
		UsageText: "iktodo remove [subcommand] [id1 id2 ...]",
		Before:    db.InitDB,
		Subcommands: []*cli.Command{
			{
				Name:  "done",
				Usage: "remove all finished task",
				Action: func(ctx *cli.Context) error {
					return db.RemoveFinished()
				},
			},
		},
		Action: func(ctx *cli.Context) error {
			ids, err := common.StringsToInts(ctx.Args().Slice())
			if err != nil {
				return errors.Join(errors.New("id error"), err)
			}

			return db.DeleteTask(ids)
		},
	}
}
