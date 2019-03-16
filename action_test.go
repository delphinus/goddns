package main

import (
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/urfave/cli.v2"
)

func TestTick(t *testing.T) {
	a := assert.New(t)
	defer prepareConfig(t)()
	defer prepareAddressOK(t, "192.168.100.100")()
	defer prepareCacheOK(t)()
	defer prepareUpdaterOK(t)()
	defer prepareTick(t)()
	exit := make(chan int)
	resultsChan := make(chan results)
	go tick(exit, resultsChan)
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
	defer prepareConfig(t)()
	defer prepareAddressOK(t, "192.168.100.100")()
	defer prepareCacheOK(t)()
	defer prepareUpdaterOK(t)()
	defer prepareTick(t)()
	go func() {
		time.Sleep(1500 * time.Millisecond)
		t.Logf("sending sig")
		sig <- syscall.SIGINT
	}()
	a.NoError(Action(&cli.Context{}))
	time.Sleep(1 * time.Second)
}

func prepareTick(t *testing.T) func() {
	original := tickIntervalSeconds
	tickIntervalSeconds = time.Second
	return func() {
		tickIntervalSeconds = original
	}
}
