package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfigNotTOML(t *testing.T) {
	a := assert.New(t)
	env := NewEnv()
	defer prepareConfigDetail(t, env, `hogehogeo`)()
	_, err := LoadConfig(env)
	a.Error(err)
	t.Logf("err: %s", err)
}

func TestLoadConfigInvalid(t *testing.T) {
	a := assert.New(t)
	env := NewEnv()
	defer prepareConfigDetail(t, env, `
[[domains]]
username = 'hogehogeo'
password = 'fugafugao'
hostname = '123.invalid.example.com'
`)()
	_, err := LoadConfig(env)
	a.Error(err)
	t.Logf("err: %s", err)
}

func TestLoadConfigValid(t *testing.T) {
	a := assert.New(t)
	env := NewEnv()
	defer prepareConfig(t, env, "")()
	_, err := LoadConfig(env)
	a.NoError(err)
}

func prepareConfig(t *testing.T, env *Env, checkIPURL string) func() {
	return prepareConfigDetail(t, env, fmt.Sprintf(`
interval = 1
check_ip_url = '%s'

[[domains]]
username = 'hogehogeo'
password = 'fugafugao'
hostname = 'example.com'

[[domains]]
username = 'hogehogeo2'
password = 'fugafugao2'
hostname = 'host2.example.com'
`, checkIPURL))
}

func prepareConfigDetail(t *testing.T, env *Env, config string) func() {
	a := assert.New(t)
	tmpDir, err := ioutil.TempDir("", "")
	a.NoError(err)
	env.ConfigFilename = path.Join(tmpDir, "goddns.toml")
	a.NoError(ioutil.WriteFile(env.ConfigFilename, []byte(config), 0600))
	return func() {
		a.NoError(os.RemoveAll(tmpDir))
	}
}
