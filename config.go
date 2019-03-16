package main

import (
	"github.com/BurntSushi/toml"
	"golang.org/x/xerrors"
	validator "gopkg.in/go-playground/validator.v9"
)

const defaultIntervalSeconds = 60

var validate = validator.New()
var configFilename = "/usr/local/etc/goddns.toml"

type Configs struct {
	Interval int       `toml:"interval" validate:"gt=0"`
	Domains  []*Domain `toml:"domains" validate:"gt=0,dive,required"`
}

type Domain struct {
	Username string `toml:"username" validate:"required"`
	Password string `toml:"password" validate:"required"`
	Hostname string `toml:"hostname" validate:"fqdn,required"`
}

var Config *Configs

func LoadConfig() error {
	Config = &Configs{Interval: defaultIntervalSeconds}
	if _, err := toml.DecodeFile(configFilename, Config); err != nil {
		return xerrors.Errorf(": %w", err)
	}
	if err := validate.Struct(Config); err != nil {
		return xerrors.Errorf(": %w", err)
	}
	return nil
}
