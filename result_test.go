package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewResult(t *testing.T) {
	a := assert.New(t)
	for _, c := range []struct {
		content      string
		isSuccessful bool
		isCritical   bool
		hasError     bool
	}{
		{content: "good 1.2.3.4", isSuccessful: true},
		{content: "good fd3a:6175:2c72:b94f::1", isSuccessful: true},
		{content: "good hoge", hasError: true},
		{content: "good  1.2.3.4", hasError: true},
		{content: "good 1.2.3.4 extra", hasError: true},
		{content: "nochg 1.2.3.4", isSuccessful: true},
		{content: "nohost", isCritical: true},
		{content: "badauth", isCritical: true},
		{content: "notfqdn", isCritical: true},
		{content: "badagent", isCritical: true},
		{content: "abuse", isCritical: true},
		{content: "911"},
	} {
		r := bytes.NewBuffer([]byte(c.content))
		result, err := NewResult(r)
		if c.hasError {
			t.Logf("err: %v", err)
		} else {
			a.NoError(err, c.content)
			a.Implements((*Result)(nil), result, c.content)
			a.Equal(c.isSuccessful, result.IsSuccessful(), c.content)
			a.Equal(c.isCritical, result.IsCritical(), c.content)
			if c.isSuccessful {
				a.Contains(result.String(), "Successful!", c.content)
			} else {
				a.Contains(result.String(), "Failed...", c.content)
			}
		}
	}
}
