package main

import (
	"github.com/urfave/cli/v2"
)

// get from build flags
var Version string = "Get From Build Flag"

func NewApp() *cli.App {
	app := cli.App{
		Name:           "iktodo",
		Version:        Version,
		Authors:        []*cli.Author{{Name: "Alctny", Email: "ltozvxe@gmail.com"}},
		DefaultCommand: "list",
		Commands: []*cli.Command{
			AddCommand(),
			ListCommand(),
			DoneCommand(),
		},
	}

	return &app
}
