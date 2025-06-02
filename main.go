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
		count     int
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
						Name:        "<ip> ",
						Destination: &inputs.primary,
					},
					&cli.StringArg{
						Name:        "<ranges>",
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
				Usage:                  "Given a CIDR, return the range",
				Arguments: []cli.Argument{
					&cli.StringArg{
						Name:        "<cidr>",
						Destination: &inputs.primary,
					},
				},
				Action: func(context.Context, *cli.Command) error {
					result := CIDRBoundaries(inputs.primary)
					return handleResult(&result)
				},
			},
			{
				UseShortOptionHandling: true,
				Name:                   "next",
				Usage:                  "Get next IP",
				Arguments: []cli.Argument{
					&cli.StringArg{
						Name:        "<ip>",
						Destination: &inputs.primary,
					},
				},
				Action: func(context.Context, *cli.Command) error {
					result := Next(inputs.primary)
					return handleResult(&result)
				},
			},
			{
				UseShortOptionHandling: true,
				Name:                   "prev",
				Usage:                  "Get previous IP",
				Arguments: []cli.Argument{
					&cli.StringArg{
						Name:        "<ip>",
						Destination: &inputs.primary,
					},
				},
				Action: func(context.Context, *cli.Command) error {
					result := Prev(inputs.primary)
					return handleResult(&result)
				},
			},
			{
				UseShortOptionHandling: true,
				Name:                   "getn",
				Usage:                  "Get N IPs from CIDR, not including the network",
				Arguments: []cli.Argument{
					&cli.StringArg{
						Name:        "<cidr>",
						Destination: &inputs.primary,
					},
					&cli.IntArg{
						Name:        "<count>",
						Destination: &inputs.count,
					},
				},
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:    "offset",
						Value:   0,
						Usage:   "Number of IPs to skip before starting to return results",
						Aliases: []string{"o"},
					},
					&cli.BoolFlag{
						Name:    "tail",
						Usage:   "Count backwards from the end of the range",
						Aliases: []string{"t"},
					},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					result := GetN(inputs.primary, inputs.count, cmd.Int("offset"), cmd.Bool("tail"))
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
