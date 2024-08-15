package main

import (
	"fmt"
	"os"

	"github.com/alctny/iktodo/cmd"
	"github.com/urfave/cli/v2"
)

// get from build flags
var Version string = "alpha"

func NewApp() *cli.App {
	app := cli.App{
		Name:           "iktodo",
		Version:        Version,
		Authors:        []*cli.Author{{Name: "Alctny", Email: "ltozvxe@gmail.com"}},
		DefaultCommand: "list",
		Commands: []*cli.Command{
			cmd.AddCommand(),
			cmd.ListCommand(),
			cmd.DoneCommand(),
			cmd.RemoveCmd(),
		},
	}

	return &app
}

func main() {
	app := NewApp()
	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
