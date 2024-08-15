package cmd

import (
	"fmt"

	"github.com/alctny/iktodo/db"
	"github.com/urfave/cli/v2"
)

// ListCommand 构建查询任务的 command
func ListCommand() *cli.Command {
	return &cli.Command{
		Name:      "list",
		Usage:     "list task",
		UsageText: "iktodo list [option]",
		Before:    db.InitDB,
		Action:    listAction,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "all",
				Usage:   "list all",
				Value:   false,
				Aliases: []string{"a"},
			},

			&cli.BoolFlag{
				Name:    "done",
				Usage:   "list finished task",
				Value:   false,
				Aliases: []string{"d"},
			},
		},
	}
}

func listAction(ctx *cli.Context) error {
	query := map[string]any{"status": 0}

	if ctx.IsSet("done") {
		query["status"] = -1
	}

	if ctx.IsSet("all") {
		query = map[string]any{}
	}

	ts, err := db.ListTask(query)
	if err != nil {
		return err
	}
	for _, v := range ts {
		fmt.Printf("%02d %s\n", v.ID, v.Name)
	}

	return nil

}
