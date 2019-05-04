package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"

	"golang.org/x/xerrors"
)

// Address detects my IP address
func Address(config *Configs) (string, error) {
	resp, err := http.Get(config.CheckIPURL)
	if err != nil {
		return "", xerrors.Errorf(": %w", err)
	}
	if resp.StatusCode/100 != 2 {
		return "", xerrors.New(fmt.Sprintf("%s returned %s",
			config.CheckIPURL, resp.Status))
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", xerrors.Errorf(": %w", err)
	}
	ip := string(body)
	if net.ParseIP(ip) == nil {
		return "", xerrors.New(fmt.Sprintf("body is not IP Address: %s", body))
	}
	return ip, nil
}
