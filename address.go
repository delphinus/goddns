package main

import (
	"io/ioutil"
	"net"
	"net/http"

	"golang.org/x/xerrors"
)

var checkipUrl = "https://domains.google.com/checkip"

func Address() (string, error) {
	resp, err := http.Get(checkipUrl)
	if err != nil {
		return "", xerrors.Errorf(": %w", err)
	}
	if resp.StatusCode/100 != 2 {
		return "", xerrors.Errorf("%s returned %s", checkipUrl, resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", xerrors.Errorf(": %w", err)
	}
	ip := string(body)
	if net.ParseIP(ip) == nil {
		return "", xerrors.Errorf("body is not IP Address: %s", body)
	}
	return ip, nil
}
