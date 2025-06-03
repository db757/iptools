package main

import (
	"context"
	"log"
	"os"

	"github.com/db757/iptools/internal/handlers"
	"github.com/urfave/cli/v3"
)

func main() {
	app := handlers.NewApp()

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
				Destination: &app.Config.Short,
			},
		},
		Commands: []*cli.Command{
			{
				Name:      "inrange",
				Usage:     "Check if IP is in range",
				UsageText: "ipt inrange [ip] [ranges]",
				Arguments: []cli.Argument{
					&cli.StringArg{
						Name:        "ip",
						Destination: &app.Input.Primary,
					},
					&cli.StringArg{
						Name:        "ranges",
						Destination: &app.Input.Secondary,
					},
				},
				Action: app.InRangeHandler,
			},
			{
				Name:      "cidrange",
				Usage:     "Given a CIDR, return the range",
				UsageText: "ipt cidrange [cidr]",
				Arguments: []cli.Argument{
					&cli.StringArg{
						Name:        "cidr",
						Destination: &app.Input.Primary,
					},
				},
				Action: app.CIDRBoundariesHandler,
			},
			{
				Name:      "next",
				Usage:     "Get next IP",
				UsageText: "ipt next [ip]",
				Arguments: []cli.Argument{
					&cli.StringArg{
						Name:        "ip",
						Destination: &app.Input.Primary,
					},
				},
				Action: app.NextHandler,
			},
			{
				Name:      "prev",
				Usage:     "Get previous IP",
				UsageText: "ipt prev [ip]",
				Arguments: []cli.Argument{
					&cli.StringArg{
						Name:        "ip",
						Destination: &app.Input.Primary,
					},
				},
				Action: app.PrevHandler,
			},
			{
				UseShortOptionHandling: true,
				Name:                   "getn",
				Usage:                  "Get N IPs from CIDR, not including the network address and broadcast address",
				UsageText:              "ipt getn [cidr] [count] [--offset|-o offset] [--tail|-t]",
				Arguments: []cli.Argument{
					&cli.StringArg{
						Name:        "cidr",
						Destination: &app.Input.Primary,
					},
					&cli.IntArg{
						Name:        "count",
						Destination: &app.Input.Count,
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
				Action: app.GetNHandler,
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
