package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v3"
)

var config = struct {
	short bool
}{
	short: false,
}

func main() {
	inputs := struct {
		primary   string
		secondary string
	}{}

	cmd := &cli.Command{
		Name:                   "ipt",
		Usage:                  "IP Cli",
		UseShortOptionHandling: true,
		EnableShellCompletion:  true,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "short",
				Usage:       "Short output",
				Aliases:     []string{"s"},
				Destination: &config.short,
			},
		},
		Commands: []*cli.Command{
			{
				UseShortOptionHandling: true,
				Name:                   "inrange",
				Usage:                  "Check if IP is in range",
				Arguments: []cli.Argument{
					&cli.StringArg{
						Name:        "ip",
						Destination: &inputs.primary,
					},
					&cli.StringArg{
						Name:        "ranges",
						Destination: &inputs.secondary,
					},
				},
				Action: func(context.Context, *cli.Command) error {
					result := IPInRange(inputs.primary, inputs.secondary)
					return handleResult(&result)
				},
			},
			{
				UseShortOptionHandling: true,
				Name:                   "cidrange",
				Usage:                  "given a CIDR, return the range",
				Arguments: []cli.Argument{
					&cli.StringArg{
						Name:        "cidr",
						Destination: &inputs.primary,
					},
				},
				Action: func(context.Context, *cli.Command) error {
					result := CIDRBoundaries(inputs.primary)
					return handleResult(&result)
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func handleResult(result Result) error {
	if result.Error() != nil {
		return result.Error()
	}

	if config.short {
		fmt.Println(result.Short())
		return nil
	}

	fmt.Println(result.Result())
	return nil
}
