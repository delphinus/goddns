package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/xerrors"
)

func TestStartAddress(t *testing.T) {
	a := assert.New(t)
	env := NewEnv()
	defer prepareAddressNG(t, env)()
	config, err := LoadConfig(env)
	a.NoError(err)
	_, err = Start(env, config, &Domain{Hostname: "example.com"})
	a.Error(err)
	t.Logf("err: %v", err)
}

func TestStartCache(t *testing.T) {
	a := assert.New(t)
	env := NewEnv()
	defer prepareAddressOK(t, env, "192.168.1.1")()
	config, err := LoadConfig(env)
	a.NoError(err)
	defer prepareCacheNG(t, env)()
	_, err = Start(env, config, &Domain{Hostname: "example.com"})
	a.Error(err)
	t.Logf("err: %v", err)
}

func TestStartIPIsSame(t *testing.T) {
	a := assert.New(t)
	env := NewEnv()
	defer prepareAddressOK(t, env, "192.168.1.1")()
	config, err := LoadConfig(env)
	a.NoError(err)
	defer prepareCacheOK(t, env)()
	result, err := Start(env, config, &Domain{Hostname: "example.com"})
	a.NoError(err)
	a.Implements((*Result)(nil), result)
	a.True(result.IsSuccessful())
	a.False(result.IsCritical())
	t.Logf("result: %s", result)
}

func TestStartUpdate(t *testing.T) {
	a := assert.New(t)
	env := NewEnv()
	defer prepareAddressOK(t, env, "192.168.100.100")()
	config, err := LoadConfig(env)
	a.NoError(err)
	defer prepareCacheOK(t, env)()
	defer prepareUpdaterNG(t, env)()
	_, err = Start(env, config, &Domain{Hostname: "example.com"})
	a.Error(err)
	t.Logf("err: %v", err)
}

func TestStartCritical(t *testing.T) {
	a := assert.New(t)
	env := NewEnv()
	defer prepareAddressOK(t, env, "192.168.100.100")()
	config, err := LoadConfig(env)
	a.NoError(err)
	defer prepareCacheOK(t, env)()
	defer prepareUpdaterCritical(t, env)()
	result, err := Start(env, config, &Domain{Hostname: "example.com"})
	a.NoError(err)
	a.Implements((*Result)(nil), result)
	a.False(result.IsSuccessful())
	a.True(result.IsCritical())
	t.Logf("result: %v", result)
}

func TestStartSave(t *testing.T) {
	a := assert.New(t)
	env := NewEnv()
	defer prepareAddressOK(t, env, "192.168.100.100")()
	config, err := LoadConfig(env)
	a.NoError(err)
	defer prepareCacheOK(t, env)()
	defer prepareUpdaterOK(t, env)()
	prepareCacheSaveNG(t, env)
	_, err = Start(env, config, &Domain{Hostname: "example.com"})
	a.Error(err)
	t.Logf("err: %v", err)
}

func TestStart(t *testing.T) {
	a := assert.New(t)
	env := NewEnv()
	defer prepareAddressOK(t, env, "192.168.100.100")()
	config, err := LoadConfig(env)
	a.NoError(err)
	defer prepareCacheOK(t, env)()
	defer prepareUpdaterOK(t, env)()
	result, err := Start(env, config, &Domain{Hostname: "example.com"})
	a.NoError(err)
	a.Implements((*Result)(nil), result)
	a.True(result.IsSuccessful())
	a.False(result.IsCritical())
	t.Logf("result: %s", result)
}

func prepareCacheSaveNG(t *testing.T, env *Env) {
	env.WriteFile = func(string, []byte, os.FileMode) error {
		return xerrors.New("dummy")
	}
}
