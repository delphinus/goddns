package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/xerrors"
)

func TestStartAddress(t *testing.T) {
	a := assert.New(t)
	prepareLogger(t)
	defer prepareAddressNG(t)()
	_, err := Start(&Domain{Hostname: "example.com"})
	a.Error(err)
	t.Logf("err: %v", err)
}

func TestStartCache(t *testing.T) {
	a := assert.New(t)
	prepareLogger(t)
	defer prepareAddressOK(t, "192.168.1.1")()
	defer prepareCacheNG(t)()
	_, err := Start(&Domain{Hostname: "example.com"})
	a.Error(err)
	t.Logf("err: %v", err)
}

func TestStartIPIsSame(t *testing.T) {
	a := assert.New(t)
	prepareLogger(t)
	defer prepareAddressOK(t, "192.168.1.1")()
	defer prepareCacheOK(t)()
	result, err := Start(&Domain{Hostname: "example.com"})
	a.NoError(err)
	a.Implements((*Result)(nil), result)
	a.True(result.IsSuccessful())
	a.False(result.IsCritical())
	t.Logf("result: %s", result)
}

func TestStartUpdate(t *testing.T) {
	a := assert.New(t)
	prepareLogger(t)
	defer prepareAddressOK(t, "192.168.100.100")()
	defer prepareCacheOK(t)()
	defer prepareUpdaterNG(t)()
	_, err := Start(&Domain{Hostname: "example.com"})
	a.Error(err)
	t.Logf("err: %v", err)
}

func TestStartCritical(t *testing.T) {
	a := assert.New(t)
	prepareLogger(t)
	defer prepareAddressOK(t, "192.168.100.100")()
	defer prepareCacheOK(t)()
	defer prepareUpdaterCritical(t)()
	_, err := Start(&Domain{Hostname: "example.com"})
	a.Error(err)
	t.Logf("err: %v", err)
}

func TestStartSave(t *testing.T) {
	a := assert.New(t)
	prepareLogger(t)
	defer prepareAddressOK(t, "192.168.100.100")()
	defer prepareCacheOK(t)()
	defer prepareUpdaterOK(t)()
	defer prepareCacheSaveNG(t)()
	_, err := Start(&Domain{Hostname: "example.com"})
	a.Error(err)
	t.Logf("err: %v", err)
}

func TestStart(t *testing.T) {
	a := assert.New(t)
	prepareLogger(t)
	defer prepareAddressOK(t, "192.168.100.100")()
	defer prepareCacheOK(t)()
	defer prepareUpdaterOK(t)()
	result, err := Start(&Domain{Hostname: "example.com"})
	a.NoError(err)
	a.Implements((*Result)(nil), result)
	a.True(result.IsSuccessful())
	a.False(result.IsCritical())
	t.Logf("result: %s", result)
}

func prepareCacheSaveNG(t *testing.T) func() {
	original := osOpenFile
	osOpenFile = func(string, int, os.FileMode) (*os.File, error) {
		return nil, xerrors.New("dummy")
	}
	return func() {
		osOpenFile = original
	}
}
