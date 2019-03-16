package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/logger"
	"gopkg.in/urfave/cli.v2"
)

type results struct {
	result Result
	err    error
}

var sig = make(chan os.Signal)

func Action(*cli.Context) error {
	logger.Info("start")
	exit := make(chan int)
	resultsChan := make(chan results)
	go tick(exit, resultsChan)
	signal.Notify(sig, syscall.SIGINT)
LOOP:
	for {
		select {
		case s := <-sig:
			if s == syscall.SIGINT {
				logger.Warning("SIGINT received. exiting...")
				exit <- 1
				break LOOP
			} else {
				logger.Warningf("unknwon signal: %s received.", s)
			}
		case results := <-resultsChan:
			if results.err != nil {
				logger.Warningf("error occurred. trying again later: %v, %+v", results.err, results.err)
			} else if results.result.IsCritical() {
				logger.Errorf("critical error occurred. exiting...: %+v", results.result)
				exit <- 1
				break LOOP
			}
		}
	}
	return nil
}

func tick(exit <-chan int, resultsChan chan<- results) {
	t := time.NewTicker(time.Duration(Config.Interval) * time.Second)
	process(exit, resultsChan)
LOOP:
	for {
		select {
		case <-exit:
			break LOOP
		case <-t.C:
			process(exit, resultsChan)
		}
	}
	t.Stop()
}

func process(exit <-chan int, resultsChan chan<- results) {
	logger.Infof("loading %s", configFilename)
	if err := LoadConfig(); err != nil {
		resultsChan <- results{err: err}
		return
	}
	for _, domain := range Config.Domains {
		logger.Infof("starting: %s", domain.Hostname)
		result, err := Start(domain)
		if result != nil {
			logger.Infof("result: %s", result)
		}
		resultsChan <- results{result: result, err: err}
	}
}
