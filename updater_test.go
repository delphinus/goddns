package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdater(t *testing.T) {
	a := assert.New(t)
	for _, c := range []struct {
		statusOK bool
		hasError bool
		body     string
	}{
		{statusOK: false, body: "", hasError: true},
		{statusOK: true, body: "hoge", hasError: true},
		{statusOK: true, body: "good 1.2.3.4"},
	} {
		func() {
			defer prepareUpdaterDetail(t, c.statusOK, c.body)()
			updater := NewUpdater(&Domain{Hostname: "example.com"}, "192.168.1.1")
			a.Implements((*Updater)(nil), updater)
			result, err := updater.Update()
			if c.hasError {
				a.Error(err)
				t.Logf("err: %v", err)
			} else {
				a.NoError(err)
				a.Implements((*Result)(nil), result)
			}
		}()
	}
}

func prepareUpdaterOK(t *testing.T) func() {
	return prepareUpdaterDetail(t, true, "good 1.2.3.4")
}

func prepareUpdaterNG(t *testing.T) func() {
	return prepareUpdaterDetail(t, false, "")
}

func prepareUpdaterCritical(t *testing.T) func() {
	return prepareUpdaterDetail(t, true, "nohost")
}

func prepareUpdaterDetail(t *testing.T, statusOK bool, body string) func() {
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
	original := updaterUrl
	updaterUrl = ts.URL
	return func() {
		ts.Close()
		updaterUrl = original
	}
}
