package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	var ip, ranges string
	var cidr string
	var short bool
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
				Destination: &short,
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
						Destination: &ip,
					},
					&cli.StringArg{
						Name:        "ranges",
						Destination: &ranges,
					},
				},
				Action: func(context.Context, *cli.Command) error {
					isInRange := IPInRange(ip, ranges)
					if short {
						fmt.Println(isInRange)
						return nil
					}

					if isInRange {
						fmt.Printf("%s is in %s\n", ip, ranges)
					} else {
						fmt.Printf("%s is NOT in %s\n", ip, ranges)
					}

					return nil
				},
			},
			{
				UseShortOptionHandling: true,
				Name:                   "cidrange",
				Usage:                  "given a CIDR, return the range",
				Arguments: []cli.Argument{
					&cli.StringArg{
						Name:        "cidr",
						Destination: &cidr,
					},
				},
				Action: func(context.Context, *cli.Command) error {
					from, to := CIDRBoundaries(cidr)
					if short {
						fmt.Printf("%s-%s\n", from, to)
						return nil
					}

					fmt.Printf("from: %s\nto: %s\n", from, to)
					return nil
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
