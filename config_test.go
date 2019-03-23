package main

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfigNotTOML(t *testing.T) {
	a := assert.New(t)
	defer prepareConfigDetail(t, `hogehogeo`)()
	_, err := LoadConfig()
	a.Error(err)
	t.Logf("err: %s", err)
}

func TestLoadConfigInvalid(t *testing.T) {
	a := assert.New(t)
	defer prepareConfigDetail(t, `
[[domains]]
username = 'hogehogeo'
password = 'fugafugao'
hostname = '123.invalid.example.com'
`)()
	_, err := LoadConfig()
	a.Error(err)
	t.Logf("err: %s", err)
}

func TestLoadConfigValid(t *testing.T) {
	a := assert.New(t)
	defer prepareConfig(t)()
	_, err := LoadConfig()
	a.NoError(err)
}

func prepareConfig(t *testing.T) func() {
	return prepareConfigDetail(t, `
interval = 1

[[domains]]
username = 'hogehogeo'
password = 'fugafugao'
hostname = 'example.com'

[[domains]]
username = 'hogehogeo2'
password = 'fugafugao2'
hostname = 'host2.example.com'
`)
}

func prepareConfigDetail(t *testing.T, config string) func() {
	a := assert.New(t)
	tmpDir, err := ioutil.TempDir("", "")
	a.NoError(err)
	original := configFilename
	configFilename = path.Join(tmpDir, "goddns.toml")
	a.NoError(ioutil.WriteFile(configFilename, []byte(config), 0600))
	return func() {
		a.NoError(os.RemoveAll(tmpDir))
		configFilename = original
	}
}
