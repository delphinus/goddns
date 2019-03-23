package main

import (
	"io/ioutil"
	"os"
	"time"

	validator "gopkg.in/go-playground/validator.v9"
)

const defaultConfigFilename = "/usr/local/etc/goddns.toml"

// TODO: This endpoint shows IPv6 address if available.  But I want IPv4
// Address.
// const defaultCheckIPURL = "https://domains.google.com/checkip"
const defaultCheckIPURL = "https://api.ipify.org"
const defaultUpdaterURL = "https://domains.google.com/nic/update"
const defaultCacheDir = "/usr/local/var/cache/goddns"
const defaultUpdateInterval = 300 * time.Second

// Env is used to switch contexts in testing
type Env struct {
	Sig            chan os.Signal
	Validate       *validator.Validate
	UpdateInterval time.Duration
	TimeNow        func() time.Time
	WriteFile      func(string, []byte, os.FileMode) error
	ConfigFilename string
	CheckIPURL     string
	UpdaterURL     string
	CacheDir       string
}

// NewEnv creates the default env settings
func NewEnv() *Env {
	return &Env{
		Sig:            make(chan os.Signal),
		Validate:       validator.New(),
		UpdateInterval: defaultUpdateInterval,
		TimeNow:        time.Now,
		WriteFile:      ioutil.WriteFile,
		ConfigFilename: defaultConfigFilename,
		CheckIPURL:     defaultCheckIPURL,
		UpdaterURL:     defaultUpdaterURL,
		CacheDir:       defaultCacheDir,
	}
}
