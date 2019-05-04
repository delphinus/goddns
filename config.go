package main

import (
	"github.com/BurntSushi/toml"
	"golang.org/x/xerrors"
)

// TODO: This endpoint shows IPv6 address if available.  But I want IPv4
// Address.
// const defaultCheckIPURL = "https://domains.google.com/checkip"
const defaultCheckIPURL = "https://api.ipify.org"
const defaultIntervalSeconds = 60

// Configs is a struct to define configuration of the app
type Configs struct {
	CheckIPURL string    `toml:"check_ip_url" validate:""`
	Interval   int       `toml:"interval" validate:"gt=0"`
	Domains    []*Domain `toml:"domains" validate:"gt=0,dive,required"`
}

// Domain is a struct to store setting for a domain
type Domain struct {
	Username string `toml:"username" validate:"required"`
	Password string `toml:"password" validate:"required"`
	Hostname string `toml:"hostname" validate:"fqdn,required"`
}

// LoadConfig loads config from TOML
func LoadConfig(env *Env) (*Configs, error) {
	config := &Configs{
		CheckIPURL: defaultCheckIPURL,
		Interval:   defaultIntervalSeconds,
	}
	if _, err := toml.DecodeFile(env.ConfigFilename, config); err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	if err := env.Validate.Struct(config); err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	return config, nil
}
