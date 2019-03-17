package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/BurntSushi/toml"
	"golang.org/x/xerrors"
)

var cacheDir = "/usr/local/var/cache/goddns"
var timeNow = time.Now // for testing
var updateIntervalSeconds = time.Duration(60*5) * time.Second
var writeFile = ioutil.WriteFile // for testing

type Cache interface {
	CanUpdate() error
	Filename() string
	IsSame(ip string) bool
	Save(ip string) error
}

type Caches struct {
	IP           string    `toml:"ip" validate:"ip,required"`
	CreatedAt    time.Time `toml:"createdAt" validate:"required"`
	UpdatedAt    time.Time `toml:"updatedAt" validate:"required"`
	CanUpdatedIn time.Time `toml:"canUpdatedIn" validate:"required"`
	domain       *Domain
	filename     string
}

func NewCache(domain *Domain) (Cache, error) {
	cache := &Caches{
		domain:   domain,
		filename: path.Join(cacheDir, domain.Hostname+".cache"),
	}
	if st, err := os.Stat(cache.filename); os.IsNotExist(err) || st.IsDir() {
		return cache, nil
	}
	if _, err := toml.DecodeFile(cache.filename, cache); err != nil {
		return nil, xerrors.Errorf("%s: %w", cache.filename, err)
	}
	if err := validate.Struct(cache); err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	return cache, nil
}

func (c *Caches) CanUpdate() error {
	if c.CanUpdatedIn.After(timeNow()) {
		return xerrors.New(fmt.Sprintf("%s cannot be updated in %s",
			c.domain.Hostname, c.CanUpdatedIn))
	}
	return nil
}

func (c *Caches) Filename() string { return c.filename }

func (c *Caches) IsSame(ip string) bool { return ip == c.IP }

func (c *Caches) Save(ip string) error {
	if st, err := os.Stat(cacheDir); os.IsNotExist(err) || !st.IsDir() {
		if err := os.MkdirAll(cacheDir, 0775); err != nil {
			return xerrors.Errorf(": %w", err)
		}
	}
	now := timeNow()
	c.IP = ip
	c.UpdatedAt = now
	c.CanUpdatedIn = now.Add(updateIntervalSeconds)
	if c.CreatedAt.IsZero() {
		c.CreatedAt = now
	}
	if err := validate.Struct(c); err != nil {
		return xerrors.Errorf(": %w", err)
	}
	w := bytes.NewBuffer(nil)
	if err := toml.NewEncoder(w).Encode(c); err != nil {
		return xerrors.Errorf(": %w", err)
	}
	if err := writeFile(c.filename, w.Bytes(), 0644); err != nil {
		return xerrors.Errorf(": %w", err)
	}
	return nil
}
