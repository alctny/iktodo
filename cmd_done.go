package main

import (
	"errors"
	"strconv"

	"github.com/urfave/cli/v2"
)

func DoneCommand() *cli.Command {
	return &cli.Command{
		Name:   "done",
		Usage:  "set task to finished",
		Before: initDB,
		Action: action,
	}
}

func action(ctx *cli.Context) error {
	idStr := ctx.Args().First()
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return errors.Join(errors.New("id error"))
	}
	return DoneTask(int(id))
}