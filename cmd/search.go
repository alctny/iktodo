package cmd

import (
	"fmt"

	"github.com/alctny/iktodo/db"
	"github.com/urfave/cli/v2"
)

func FindCmd() *cli.Command {
	return &cli.Command{
		Name:      "find",
		Usage:     "search task",
		UsageText: "iktodo find <key words>",
		Before:    db.InitDB,
		Action: func(ctx *cli.Context) error {
			kw := ctx.Args().First()
			ts, err := db.Search(kw)
			if err != nil {
				return err
			}

			for _, t := range ts {
				fmt.Println(t.ID, t.Name)
			}
			return nil
		},
	}
}
