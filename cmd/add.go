package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/alctny/iktodo/db"
	"github.com/alctny/iktodo/task"
	"github.com/urfave/cli/v2"
)

func AddCommand() *cli.Command {
	return &cli.Command{
		Name:   "add",
		Before: db.InitDB,
		Action: func(ctx *cli.Context) error {
			if !ctx.Args().Present() {
				return fmt.Errorf("iktodo add <task>")
			}

			name := strings.Join(ctx.Args().Slice(), " ")
			t := &task.Task{
				Status:   0,
				Name:     name,
				CreateAt: time.Now(),
			}

			return db.SaveTask(t)
		},
	}
}
