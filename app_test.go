package main

import (
	"io/ioutil"

	"github.com/google/logger"
)

func init() {
	_ = logger.Init("goddns", true, false, ioutil.Discard)
}
