package main

import (
	"os"

	"github.com/google/logger"
)

func init() {
	_ = logger.Init("goddns", false, false, os.Stderr)
}
