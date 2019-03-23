package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/logger"
	"golang.org/x/xerrors"
	"gopkg.in/urfave/cli.v2"
)

type results struct {
	result Result
	err    error
}

// Action is the main logic for the app
func Action(sig chan os.Signal) func(*cli.Context) error {
	if sig == nil {
		sig = make(chan os.Signal)
	}
	return func(*cli.Context) error {
		logger.Info("start")
		config, err := LoadConfig()
		if err != nil {
			return xerrors.Errorf(": %w", err)
		}
		newConfig := make(chan *Configs)
		exit := make(chan int)
		resultsChan := make(chan results)
		go tick(config, newConfig, exit, resultsChan)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGHUP)
	LOOP:
		for {
			select {
			case s := <-sig:
				if manageSig(s, newConfig) {
					break LOOP
				}
			case r := <-resultsChan:
				if manageResults(r) {
					break LOOP
				}
			}
		}
		exit <- 1
		return nil
	}
}

func manageSig(s os.Signal, newConfig chan<- *Configs) bool {
	switch s {
	case syscall.SIGINT:
		logger.Warning("SIGINT received. exiting...")
		return true
	case syscall.SIGHUP:
		logger.Warning("SIGHUP received. reloading configs...")
		config, err := LoadConfig()
		if err != nil {
			logger.Warningf("%s", xerrors.Errorf("config has errors: %w", err))
		} else {
			newConfig <- config
		}
	default:
		logger.Warningf("unknwon signal: %s received.", s)
	}
	return false
}

func manageResults(r results) bool {
	if r.err != nil {
		logger.Warningf("error occurred. trying again later: %v, %+v", r.err, r.err)
		return false
	}
	if r.result.IsCritical() {
		logger.Errorf("critical error occurred. exiting...: %+v", r.result)
		return true
	}
	return false
}

func tick(
	config *Configs,
	newConfig <-chan *Configs,
	exit <-chan int,
	resultsChan chan<- results,
) {
	t := newTicker(config)
	process(config, exit, resultsChan)
LOOP:
	for {
		select {
		case c := <-newConfig:
			config = c
			t.Stop()
			t = newTicker(config)
		case <-exit:
			break LOOP
		case <-t.C:
			process(config, exit, resultsChan)
		}
	}
	t.Stop()
}

func newTicker(config *Configs) *time.Ticker {
	return time.NewTicker(time.Duration(config.Interval) * time.Second)
}

func process(config *Configs, exit <-chan int, resultsChan chan<- results) {
	logger.Infof("loading %s", configFilename)
	for _, domain := range config.Domains {
		logger.Infof("starting: %s", domain.Hostname)
		result, err := Start(domain)
		if result != nil {
			logger.Infof("result: %s", result)
		}
		resultsChan <- results{result: result, err: err}
	}
}
