package cmd

import (
	"fmt"

	"github.com/alctny/iktodo/db"
	"github.com/urfave/cli/v2"
)

func ArrgregationCmd() *cli.Command {
	return &cli.Command{
		Name:      "aggregation",
		Usage:     "task finished status",
		UsageText: "iktodo aggregation",
		Aliases:   []string{"aggre"},
		Before:    db.InitDB,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "simple",
				Usage:   "output message clean",
				Aliases: []string{"s"},
				Value:   false,
			},
		},
		Action: func(ctx *cli.Context) error {
			result, err := db.Aggregate()
			if err != nil {
				return err
			}

			if ctx.IsSet("simple") {
				fmt.Println(result.Unfinish, result.Finished, result.Total)
				return nil
			}

			fmt.Println("Unfinish: ", result.Unfinish)
			fmt.Println("Finished: ", result.Finished)
			fmt.Println("Total:    ", result.Total)

			return nil
		},
	}
}
