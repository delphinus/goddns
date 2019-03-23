package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/xerrors"
)

// Updater is an interface to update
type Updater interface {
	Update() (Result, error)
}

// Updaters is an implementation of Updater
type Updaters struct {
	env    *Env
	domain *Domain
	ip     string
}

// NewUpdater creates Updater
func NewUpdater(env *Env, domain *Domain, ip string) Updater {
	return &Updaters{env, domain, ip}
}

// Update updates the IP
func (u *Updaters) Update() (Result, error) {
	req, err := u.req()
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	if resp.StatusCode/100 != 2 {
		return nil, xerrors.New(fmt.Sprintf("%s returned %s",
			u.env.UpdaterURL, resp.Status))
	}
	result, err := NewResult(resp.Body)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	return result, nil
}

func (u *Updaters) req() (*http.Request, error) {
	v := url.Values{}
	v.Set("hostname", u.domain.Hostname)
	v.Set("myip", u.ip)
	req, err := http.NewRequest(
		"POST",
		u.env.UpdaterURL,
		strings.NewReader(v.Encode()),
	)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	req.SetBasicAuth(u.domain.Username, u.domain.Password)
	return req, nil
}
