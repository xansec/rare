package cmd

import (
	"fmt"

	"github.com/urfave/cli"
)

func filterFunction(c *cli.Context) error {
	fmt.Println("Howdy")

	writeLines := c.Bool("line")

	extractor := buildExtractorFromArguments(c)
	for {
		match, more := <-extractor.ReadChan
		if !more {
			break
		}
		if writeLines {
			fmt.Printf("%d: %s\n", match.LineNumber, match.Extracted)
		} else {
			fmt.Println(match.Extracted)
		}
	}
	return nil
}

// HistogramCommand Exported command
func FilterCommand() *cli.Command {
	return &cli.Command{
		Name:      "filter",
		Usage:     "Filter incoming results with search criteria, and output raw matches",
		Action:    filterFunction,
		ArgsUsage: "<-|filename>",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "line,l",
				Usage: "Output line numbers",
			},
		},
	}
}
