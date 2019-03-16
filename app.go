package main

import (
	"fmt"
	"io/ioutil"

	"github.com/google/logger"
	"gopkg.in/urfave/cli.v2"
)

// NewApp creates app
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
		Action: handlExit(Action(nil)),
	}
}

// Before sets up logger
func Before(c *cli.Context) error {
	_ = logger.Init("goddns", c.Bool("verbose"), true, ioutil.Discard)
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
