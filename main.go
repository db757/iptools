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
	cmd := &cli.Command{
		UseShortOptionHandling: true,
		Name:                   "ipt",
		Usage:                  "IP Cli",

		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "verbose",
				Usage:   "Enable verbose output",
				Aliases: []string{"v"},
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					if IPInRange(ip, ranges) {
						fmt.Printf("%s is in %s\n", ip, ranges)
					} else {
						fmt.Printf("%s is NOT in %s\n", ip, ranges)
					}
					return nil
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
