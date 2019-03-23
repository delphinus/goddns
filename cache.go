package main

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/BurntSushi/toml"
	"golang.org/x/xerrors"
)

// Cache is an interface to deal with caches
type Cache interface {
	CanUpdate() error
	Filename() string
	IsSame(ip string) bool
	Save(ip string) error
}

// Caches is an implementation of Cache
type Caches struct {
	IP           string    `toml:"ip" validate:"ip,required"`
	CreatedAt    time.Time `toml:"createdAt" validate:"required"`
	UpdatedAt    time.Time `toml:"updatedAt" validate:"required"`
	CanUpdatedIn time.Time `toml:"canUpdatedIn" validate:"required"`
	env          *Env
	domain       *Domain
	filename     string
}

// NewCache creates Cache
func NewCache(env *Env, domain *Domain) (Cache, error) {
	cache := &Caches{
		env:      env,
		domain:   domain,
		filename: path.Join(env.CacheDir, domain.Hostname+".cache"),
	}
	if st, err := os.Stat(cache.filename); os.IsNotExist(err) || st.IsDir() {
		return cache, nil
	}
	if _, err := toml.DecodeFile(cache.filename, cache); err != nil {
		return nil, xerrors.Errorf("%s: %w", cache.filename, err)
	}
	if err := env.Validate.Struct(cache); err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	return cache, nil
}

// CanUpdate detects if the cache can be updated
func (c *Caches) CanUpdate() error {
	if c.CanUpdatedIn.After(c.env.TimeNow()) {
		return xerrors.New(fmt.Sprintf("%s cannot be updated in %s",
			c.domain.Hostname, c.CanUpdatedIn))
	}
	return nil
}

// Filename returns the filename of caches
func (c *Caches) Filename() string { return c.filename }

// IsSame detects if the cache is the same as supplied IPs
func (c *Caches) IsSame(ip string) bool { return ip == c.IP }

// Save saves caches into files
func (c *Caches) Save(ip string) error {
	if st, err := os.Stat(c.env.CacheDir); os.IsNotExist(err) || !st.IsDir() {
		if err := os.MkdirAll(c.env.CacheDir, 0750); err != nil {
			return xerrors.Errorf(": %w", err)
		}
	}
	now := c.env.TimeNow()
	c.IP = ip
	c.UpdatedAt = now
	c.CanUpdatedIn = now.Add(c.env.UpdateInterval)
	if c.CreatedAt.IsZero() {
		c.CreatedAt = now
	}
	if err := c.env.Validate.Struct(c); err != nil {
		return xerrors.Errorf(": %w", err)
	}
	w := bytes.NewBuffer(nil)
	if err := toml.NewEncoder(w).Encode(c); err != nil {
		return xerrors.Errorf(": %w", err)
	}
	if err := c.env.WriteFile(c.filename, w.Bytes(), 0644); err != nil {
		return xerrors.Errorf(": %w", err)
	}
	return nil
}
