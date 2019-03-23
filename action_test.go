package main

import (
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/urfave/cli.v2"
)

func TestTick(t *testing.T) {
	a := assert.New(t)
	env := NewEnv()
	defer prepareConfig(t, env)()
	defer prepareAddressOK(t, env, "192.168.100.100")()
	defer prepareCacheOK(t, env)()
	defer prepareUpdaterOK(t, env)()
	newConfig := make(chan *Configs)
	exit := make(chan int)
	resultsChan := make(chan results)
	config, err := LoadConfig(env)
	a.NoError(err)
	go tick(env, config, newConfig, exit, resultsChan)
	r1 := <-resultsChan
	r2 := <-resultsChan
	exit <- 1
	a.NoError(r1.err)
	a.NoError(r2.err)
	a.True(r1.result.IsSuccessful())
	a.True(r2.result.IsSuccessful())
	t.Logf("r1: %s", r1.result)
	t.Logf("r2: %s", r2.result)
}

func TestAction(t *testing.T) {
	a := assert.New(t)
	env := NewEnv()
	defer prepareConfig(t, env)()
	defer prepareAddressOK(t, env, "192.168.100.100")()
	defer prepareCacheOK(t, env)()
	defer prepareUpdaterOK(t, env)()
	env.Sig = make(chan os.Signal)
	go func() {
		time.Sleep(1500 * time.Millisecond)
		t.Log("sending sig")
		env.Sig <- syscall.SIGINT
	}()
	a.NoError(Action(&cli.Context{}, env))
	time.Sleep(1 * time.Second)
}

func TestActionReloadConfig(t *testing.T) {
	a := assert.New(t)
	env := NewEnv()
	defer prepareConfig(t, env)()
	defer prepareAddressOK(t, env, "192.168.100.100")()
	defer prepareCacheOK(t, env)()
	defer prepareUpdaterOK(t, env)()
	env.Sig = make(chan os.Signal)
	go func() {
		time.Sleep(1500 * time.Millisecond)
		t.Log("sending SIGHUP")
		env.Sig <- syscall.SIGHUP
		time.Sleep(2500 * time.Millisecond)
		t.Log("sending SIGINT")
		env.Sig <- syscall.SIGINT
	}()
	a.NoError(Action(&cli.Context{}, env))
	time.Sleep(3 * time.Second)
}
