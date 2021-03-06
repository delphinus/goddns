package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"

	"golang.org/x/xerrors"
)

// Result is an interface to deal with results
type Result interface {
	IsSuccessful() bool
	IsCritical() bool
	String() string
}

type noNeedToUpdate struct{}

func (r noNeedToUpdate) IsSuccessful() bool { return true }
func (r noNeedToUpdate) IsCritical() bool   { return false }
func (r noNeedToUpdate) String() string {
	return "IP is the same. No need to update."
}

// NoNeedToUpdate returns a result that means it does not need to be updated.
func NoNeedToUpdate() Result { return noNeedToUpdate{} }

// Results is a implementation of Result
type Results struct {
	code     string
	ip       string
	messages map[string]string
}

// NewResult parses bytes and make results
func NewResult(r io.Reader) (Result, error) {
	content, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	parts := bytes.SplitN(content, []byte{' '}, 2)
	result := &Results{
		code: string(parts[0]),
		messages: map[string]string{
			"nohost": "The hostname does not exist, " +
				"or does not have Dynamic DNS enabled.",
			"badauth": "The username / password combination " +
				"is not valid for the specified host.",
			"notfqdn": "The supplied hostname " +
				"is not a valid fully-qualified domain name.",
			"badagent": "Your Dynamic DNS client is making bad requests. " +
				"Ensure the user agent is set in the request.",
			"abuse": "Dynamic DNS access for the hostname has been blocked " +
				"due to failure to interpret previous responses correctly.",
			"911": "An error happened on our end. Wait 5 minutes and retry.",
		},
	}
	if len(parts) > 1 {
		result.ip = string(parts[1])
	}
	if !result.isValid() {
		return nil, xerrors.New(fmt.Sprintf("content: '%s' is invalid", content))
	}
	return result, nil
}

// IsSuccessful returns if the result is successful
func (r *Results) IsSuccessful() bool {
	return r.code == "good" || r.code == "nochg"
}

// IsCritical returns if the result is critical
func (r *Results) IsCritical() bool {
	return !r.IsSuccessful() && r.code != "911"
}

// String is needed for fmt.Stringer
func (r *Results) String() string {
	if r.IsSuccessful() {
		return fmt.Sprintf("Successful! code: %s, ip: %s", r.code, r.ip)
	}
	return fmt.Sprintf("Failed... code: %s, %s", r.code, r.messages[r.code])
}

func (r *Results) isValid() bool {
	if r.code == "good" || r.code == "nochg" {
		return net.ParseIP(r.ip) != nil
	}
	_, ok := r.messages[r.code]
	return ok
}
