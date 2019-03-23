package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/google/logger"
)

func TestMain(m *testing.M) {
	_ = logger.Init("goddns", true, false, ioutil.Discard)
	os.Exit(m.Run())
}
