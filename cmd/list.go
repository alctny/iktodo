package cmd

import (
	"fmt"

	"github.com/alctny/iktodo/db"
	"github.com/alctny/iktodo/task"
	"github.com/urfave/cli/v2"
)

// ListCmd 构建查询任务的 command
func ListCmd() *cli.Command {
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

			&cli.StringFlag{
				Name:    "color",
				Usage:   "color output",
				Value:   "all",
				Aliases: []string{"C"},
			},

			&cli.BoolFlag{
				Name:    "reverse",
				Usage:   "reverse list",
				Value:   false,
				Aliases: []string{"r"},
			},

			&cli.UintFlag{
				Name:    "size",
				Usage:   "set page size, default 9 when page is set",
				Aliases: []string{"s"},
			},

			&cli.UintFlag{
				Name:    "page",
				Usage:   "select the page",
				Aliases: []string{"offset", "p"},
				Action: func(ctx *cli.Context, u uint) error {
					if !ctx.IsSet("size") {
						return ctx.Set("size", "9")
					}
					return nil
				},
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

	page := ctx.Uint("page")
	size := ctx.Uint("size")
	var offset uint
	if ctx.IsSet("page") {
		offset = (page - 1) * size
	}

	ts, err := db.ListTask(query, size, offset)
	if err != nil {
		return err
	}

	// 一个勉强还算有创意，但并没有带来改善的写法
	// 而且无法解决有任务置顶的情况，所以必然要被砍掉
	length := len(ts)
	rev := func(i int) int { return length - i - 1 }
	sort := func(i int) int { return i }
	if ctx.IsSet("r") {
		sort = rev
	}

	// TODO 在进行分页之后输出结果中需要包含分页情况 [page/total_page]
	for i := range ts {
		fmt.Println(ts[sort(i)].ColorString(task.ColorWhen(ctx.String("color"))))
	}
	return nil
}
