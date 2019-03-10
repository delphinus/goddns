package main

import (
	"log/syslog"

	"golang.org/x/xerrors"
	"gopkg.in/urfave/cli.v2"
)

type Loggers interface {
	Info(m string) error
	Warning(m string) error
	Crit(m string) error
}

var logger Loggers

func Logger(*cli.Context) (err error) {
	logger, err = syslog.New(syslog.LOG_INFO, "goddns")
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}
	return nil
}