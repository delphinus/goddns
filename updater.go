package main

import (
	"fmt"
	"net/http"
	"net/url"

	"golang.org/x/xerrors"
)

var updaterUrl = "https://domains.google.com/nic/update"

type Updater interface {
	Update() (Result, error)
}

type Updaters struct {
	domain *Domain
	ip     string
}

func NewUpdater(domain *Domain, ip string) Updater { return &Updaters{domain, ip} }

func (u *Updaters) Update() (Result, error) {
	url, err := u.url()
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	resp, err := http.Get(url)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	if resp.StatusCode/100 != 2 {
		return nil, xerrors.New(fmt.Sprintf("%s returned %s", updaterUrl, resp.Status))
	}
	result, err := NewResult(resp.Body)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	return result, nil
}

func (u *Updaters) url() (string, error) {
	v := url.Values{}
	v.Set("hostname", u.domain.Hostname)
	v.Set("myip", u.ip)
	urls, err := url.Parse(updaterUrl)
	if err != nil {
		return "", xerrors.Errorf(": %w", err)
	}
	urls.User = url.UserPassword(u.domain.Username, u.domain.Password)
	urls.RawQuery = v.Encode()
	return urls.String(), nil
}
