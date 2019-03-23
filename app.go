package main

import (
	"fmt"
	"io/ioutil"

	"github.com/google/logger"
	"gopkg.in/urfave/cli.v2"
)

// NewApp creates app
func NewApp() *cli.App {
	env := NewEnv()
	return &cli.App{
		Usage:   "Update entries in Google Domains",
		Version: Version,
		Authors: []*cli.Author{
			{Name: "JINNOUCHI Yasushi", Email: "me@delphinus.dev"},
		},
		Before: handlExit(env, Before),
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "verbose",
				Usage: "Show logs in addition to syslog",
			},
		},
		Action: handlExit(env, Action),
	}
}

// Before sets up logger
func Before(c *cli.Context, env *Env) error {
	_ = logger.Init("goddns", c.Bool("verbose"), true, ioutil.Discard)
	return nil
}

func handlExit(
	env *Env,
	handler func(*cli.Context, *Env) error,
) func(*cli.Context) error {
	return func(c *cli.Context) error {
		if err := handler(c, env); err != nil {
			logger.Warning(fmt.Sprintf("error received: %+v", err))
			return cli.Exit(fmt.Sprintf("%+v", err), 1)
		}
		return nil
	}
}
