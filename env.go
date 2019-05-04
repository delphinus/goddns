package main

import (
	"io/ioutil"
	"os"
	"time"

	validator "gopkg.in/go-playground/validator.v9"
)

const defaultConfigFilename = "/usr/local/etc/goddns.toml"
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
		UpdaterURL:     defaultUpdaterURL,
		CacheDir:       defaultCacheDir,
	}
}
