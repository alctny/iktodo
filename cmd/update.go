package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/alctny/iktodo/db"
	"github.com/urfave/cli/v2"
)

func UpdateCmd() *cli.Command {
	return &cli.Command{
		Name:      "update",
		Usage:     "Update the iktodo database",
		UsageText: "iktodo update <id> <name>",
		Before:    db.InitDB,
		Action: func(c *cli.Context) error {
			if c.NArg() < 2 {
				fmt.Fprintln(os.Stderr, "Error: missing arguments")
				cli.ShowSubcommandHelpAndExit(c, 1)
			}

			idStr := c.Args().First()
			id, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				return err
			}

			args := c.Args().Slice()[1:]
			name := strings.Join(args, " ")

			return db.Update(int(id), map[string]any{"name": name})
		},
	}
}
