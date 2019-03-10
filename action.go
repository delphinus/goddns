package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gopkg.in/urfave/cli.v2"
)

var tickIntervalSeconds = 60 * time.Second

type results struct {
	result Result
	err    error
}

var sig = make(chan os.Signal)

func Action(*cli.Context) error {
	_ = logger.Info("start")
	exit := make(chan int)
	resultsChan := make(chan results)
	go tick(exit, resultsChan)
	signal.Notify(sig, syscall.SIGINT)
LOOP:
	for {
		select {
		case s := <-sig:
			if s == syscall.SIGINT {
				_ = logger.Warning("SIGINT received. exiting...")
				exit <- 1
				break LOOP
			} else {
				_ = logger.Warning(fmt.Sprintf("unknwon signal: %s received.", s))
			}
		case results := <-resultsChan:
			if results.err != nil {
				_ = logger.Warning(fmt.Sprintf("error occurred. trying again later: %+v", results.err))
			} else if results.result.IsCritical() {
				_ = logger.Crit(fmt.Sprintf("critical error occurred. exiting...: %+v", results.result))
				exit <- 1
				break LOOP
			}
		}
	}
	return nil
}

func tick(exit <-chan int, resultsChan chan<- results) {
	t := time.NewTicker(tickIntervalSeconds)
LOOP:
	for {
		select {
		case <-exit:
			break LOOP
		case <-t.C:
			_ = logger.Info(fmt.Sprintf("loading %s", configFilename))
			if err := LoadConfig(); err != nil {
				resultsChan <- results{err: err}
				continue LOOP
			}
			for _, domain := range Config.Domains {
				_ = logger.Info(fmt.Sprintf("starting: %s", domain.Hostname))
				result, err := Start(domain)
				_ = logger.Info(fmt.Sprintf("result: %s", result))
				resultsChan <- results{result: result, err: err}
			}
		}
	}
	t.Stop()
}
