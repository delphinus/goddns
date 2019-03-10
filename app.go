package main

import (
	"fmt"

	"gopkg.in/urfave/cli.v2"
)

func NewApp() *cli.App {
	return &cli.App{
		Usage:   "Update entries in Google Domains",
		Version: Version,
		Authors: []*cli.Author{
			{Name: "JINNOUCHI Yasushi", Email: "me@delphinus.dev"},
		},
		Before: handlExit(Before),
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "verbose",
				Usage: "Show logs in addition to syslog",
			},
		},
		Action: handlExit(Action),
	}
}

func Before(c *cli.Context) error {
	if err := Logger(c); err != nil {
		return err
	}
	return nil
}

func handlExit(handler func(*cli.Context) error) func(*cli.Context) error {
	return func(c *cli.Context) error {
		if err := handler(c); err != nil {
			logger.Warning(fmt.Sprintf("error received: %+v", err))
			return cli.Exit(fmt.Sprintf("%+v", err), 1)
		}
		return nil
	}
}
