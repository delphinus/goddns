package main

import (
	"github.com/BurntSushi/toml"
	"golang.org/x/xerrors"
	validator "gopkg.in/go-playground/validator.v9"
)

const defaultIntervalSeconds = 60

var validate = validator.New()
var configFilename = "/usr/local/etc/goddns.toml"

// Configs is a struct to define configuration of the app
type Configs struct {
	Interval int       `toml:"interval" validate:"gt=0"`
	Domains  []*Domain `toml:"domains" validate:"gt=0,dive,required"`
}

// Domain is a struct to store setting for a domain
type Domain struct {
	Username string `toml:"username" validate:"required"`
	Password string `toml:"password" validate:"required"`
	Hostname string `toml:"hostname" validate:"fqdn,required"`
}

// LoadConfig loads config from TOML
func LoadConfig() (*Configs, error) {
	config := &Configs{Interval: defaultIntervalSeconds}
	if _, err := toml.DecodeFile(configFilename, config); err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	if err := validate.Struct(config); err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	return config, nil
}
