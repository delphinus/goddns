package main

import (
	"testing"
)

func prepareLogger(t *testing.T) {
	logger = testLogger{t: t}
}

type testLogger struct {
	t *testing.T
}

func (l testLogger) Info(m string) error {
	l.t.Logf("Info: %s", m)
	return nil
}

func (l testLogger) Warning(m string) error {
	l.t.Logf("Warning: %s", m)
	return nil
}

func (l testLogger) Crit(m string) error {
	l.t.Logf("Crit: %s", m)
	return nil
}
