package main

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	cmd := buildApp()
	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func buildApp() *cli.Command {
	return &cli.Command{
		Name:  "supervisor",
		Usage: "Multi-platform project tracking and reporting tool",
		Commands: []*cli.Command{
			{
				Name:  "diff",
				Usage: "Generate structured JSON diff between two git references",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "repo",
						Usage:    "Path to local git repository",
						Required: false,
					},
					&cli.StringFlag{
						Name:     "from",
						Usage:    "Starting git reference (tag/branch/commit)",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "to",
						Usage:    "Target git reference (tag/branch/commit)",
						Required: true,
					},
					&cli.StringSliceFlag{
						Name:  "exclude-suffix",
						Usage: "File suffixes to exclude (e.g., .png)",
					},
					&cli.StringSliceFlag{
						Name:  "exclude-path",
						Usage: "Path prefixes to exclude (e.g., vendor/)",
					},
				},
				Action: runDiff,
			},
		},
	}
}
