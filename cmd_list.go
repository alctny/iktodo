package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

// ListCommand 构建查询任务的 command
func ListCommand() *cli.Command {
	return &cli.Command{
		Name:   "list",
		Before: initDB,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "all",
				Usage:   "list all",
				Value:   false,
				Aliases: []string{"a"},
			},
		},
		Action: func(ctx *cli.Context) error {
			all := ctx.IsSet("all")

			ts, err := ListTask(all)
			if err != nil {
				return err
			}
			for _, v := range ts {
				fmt.Printf("%02d %s\n", v.ID, v.Name)
			}

			return nil
		},
	}
}
