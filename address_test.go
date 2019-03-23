package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddress(t *testing.T) {
	a := assert.New(t)
	for _, c := range []struct {
		statusOK bool
		hasError bool
		body     string
	}{
		{statusOK: false, body: "", hasError: true},
		{statusOK: true, body: "hoge", hasError: true},
		{statusOK: true, body: "192.168.1.1"},
		{statusOK: true, body: "fd3a:6175:2c72:b94f::1"},
	} {
		func() {
			env := NewEnv()
			defer prepareAddressDetail(t, env, c.statusOK, c.body)()
			ip, err := Address(env)
			if c.hasError {
				a.Error(err)
				t.Logf("err: %v", err)
			} else {
				a.Equal(c.body, ip)
				a.NoError(err)
			}
		}()
	}
}

func prepareAddressOK(t *testing.T, env *Env, ip string) func() {
	return prepareAddressDetail(t, env, true, ip)
}

func prepareAddressNG(t *testing.T, env *Env) func() {
	return prepareAddressDetail(t, env, false, "")
}

func prepareAddressDetail(
	t *testing.T,
	env *Env,
	statusOK bool,
	body string,
) func() {
	a := assert.New(t)
	ts := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if !statusOK {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			_, err := io.WriteString(w, body)
			a.NoError(err)
		}))
	env.CheckIPURL = ts.URL
	return ts.Close
}
