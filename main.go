package main

import (
	"fmt"
	"log"
	"os"

	"rare/cmd"
	"rare/pkg/color"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Usage = "A regex parser and extractor"

	app.Version = fmt.Sprintf("%s, %s", version, buildSha)

	app.Description = `Aggregate and display information parsed from text files using
	regex and a simple handlebars-like expression syntax.
	
	https://github.com/zix99/rare`

	app.Copyright = `rare  Copyright (C) 2019 Chris LaPointe
    This program comes with ABSOLUTELY NO WARRANTY.
    This is free software, and you are welcome to redistribute it
	under certain conditions`

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "nocolor,nc",
			Usage: "Disables color output",
		},
		cli.BoolFlag{
			Name:  "color",
			Usage: "Force-enable color output",
		},
	}

	app.Commands = []cli.Command{
		*cmd.FilterCommand(),
		*cmd.HistogramCommand(),
		*cmd.HelpCommand(),
	}

	app.Before = cli.BeforeFunc(func(c *cli.Context) error {
		if c.Bool("nocolor") {
			color.Enabled = false
		} else if c.Bool("color") {
			color.Enabled = true
		}
		return nil
	})

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
