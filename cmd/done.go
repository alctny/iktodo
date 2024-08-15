package cmd

import (
	"errors"

	"github.com/alctny/iktodo/common"
	"github.com/alctny/iktodo/db"
	"github.com/urfave/cli/v2"
)

func DoneCommand() *cli.Command {
	return &cli.Command{
		Name:      "done",
		Usage:     "set task to finished",
		UsageText: "iktodo done <id1 id2 ...>",
		Before:    db.InitDB,
		Action:    doneAction,
	}
}

func doneAction(ctx *cli.Context) error {
	ids, err := common.StringsToInts(ctx.Args().Slice())
	if err != nil {
		return errors.Join(errors.New("id errror"), err)
	}

	return db.DoneTask(ids)
}
