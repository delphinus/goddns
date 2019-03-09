package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gopkg.in/urfave/cli.v2"
)

func Action(c *cli.Context) error {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT)
	exit := make(chan int)
	go func() {
		for {
			s := <-sig
			switch s {
			case syscall.SIGINT:
				exit <- 1
			default:
				exit <- 2
			}
		}
	}()
	go tick(exit)
	if code := <-exit; code > 0 {
		return cli.Exit("died", code)
	}
	return nil
}

func tick(exit chan int) {
	t := time.NewTicker(Config.Interval * time.Second)
	for {
		select {
		case c := <-exit:
			break
		case <-t.C:
			fmt.Println("hoge")
		}
	}
	t.Stop()
}
