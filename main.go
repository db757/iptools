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
	cmd := &cli.Command{
		UseShortOptionHandling: true,
		Name:                   "ipt",
		Usage:                  "IP Cli",
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
					if IPInRange(ip, ranges) {
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
					fmt.Printf("from: %s\n", from)
					fmt.Printf("to: %s\n", to)
					return nil
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
