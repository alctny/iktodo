package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
)

func AddCommand() *cli.Command {
	return &cli.Command{
		Name:   "add",
		Before: initDB,
		Action: func(ctx *cli.Context) error {
			if !ctx.Args().Present() {
				return fmt.Errorf("iktodo add <task>")
			}

			name := strings.Join(ctx.Args().Slice(), " ")
			t := &Task{
				Status:   0,
				Name:     name,
				CreateAt: time.Now(),
			}

			return SaveTask(t)
		},
	}
}
