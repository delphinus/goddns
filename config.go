package main

import (
	"github.com/BurntSushi/toml"
	"golang.org/x/xerrors"
	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/urfave/cli.v2"
)

const filename = "/usr/local/etc/goddns.toml"

type Configs struct {
	Interval int `toml:"interval" validate:"gt=0,required"`
	Domains  []struct {
		Username string `toml:"username" validate:"required"`
		Password string `toml:"password" validate:"required"`
		Domain   string `toml:"domain" validate:"fqdn,required"`
	} `toml:"domains" validate:"gt=0,dive,required"`
}

var Config = &Configs{Interval: 60}

func LoadConfig(*cli.Context) error {
	if _, err := toml.DecodeFile(filename, Config); err != nil {
		return xerrors.Errorf(": %w", err)
	}
	validate := validator.New()
	if err := validate.Struct(Config); err != nil {
		return xerrors.Errorf(": %w", err)
	}
	return nil
}
