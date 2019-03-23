package main

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewCache(t *testing.T) {
	a := assert.New(t)
	for _, c := range []struct {
		noCache  bool
		hasError bool
		cache    string
		pattern  string
	}{
		{noCache: true, pattern: "no cache"},
		{cache: `hoge`, hasError: true, pattern: "invalid cache"},
		{cache: `
			ip = 'hoge'
			createdAt = 2019-01-02T12:00:00Z
			updatedAt = 2019-01-02T12:00:00Z
			canUpdatedIn = 2019-01-02T12:05:00Z
		`, hasError: true, pattern: "invalid IP"},
		{cache: `
			ip = '192.168.1.1'
			createdAt = 2019-01-02T12:00:00Z
			updatedAt = 2019-01-02T12:00:00Z
			canUpdatedIn = 2019-01-02T12:05:00Z
		`, pattern: "cache with IPv4"},
		{cache: `
			ip = 'fd3a:6175:2c72:b94f::1'
			createdAt = 2019-01-02T12:00:00Z
			updatedAt = 2019-01-02T12:00:00Z
			canUpdatedIn = 2019-01-02T12:05:00Z
		`, pattern: "cache with IPv6"},
	} {
		func() {
			defer prepareCacheDetail(t, c.cache)()
			cache, err := NewCache(&Domain{Hostname: "example.com"})
			switch {
			case c.noCache:
				a.NoError(err, c.pattern)
				a.Implements((*Cache)(nil), cache, c.pattern)
			case c.hasError:
				a.Error(err, c.pattern)
				t.Logf("err: %v", err)
			default:
				a.NoError(err, c.pattern)
				a.Implements((*Cache)(nil), cache, c.pattern)
			}
		}()
	}
}

func TestCacheCanUpdate(t *testing.T) {
	a := assert.New(t)
	defer prepareCacheDetail(t, `
		ip = '192.168.1.1'
		createdAt = 2019-01-02T12:00:00Z
		updatedAt = 2019-01-02T12:00:00Z
		canUpdatedIn = 2019-12-01T12:00:00Z
	`)()
	cache, err := NewCache(&Domain{Hostname: "example.com"})
	a.NoError(err)
	a.Implements((*Cache)(nil), cache)
	err = cache.CanUpdate()
	a.Error(err)
	t.Logf("err: %v", err)
}

func TestCacheSave(t *testing.T) {
	a := assert.New(t)
	for _, c := range []struct {
		ip       string
		hasCache bool
		hasError bool
		pattern  string
	}{
		{ip: "hoge", hasError: true, pattern: "invalid IP"},
		{ip: "192.168.1.1", pattern: "valid IP with no cache"},
		{hasCache: true, ip: "192.168.1.1", pattern: "valid IP with cache"},
	} {
		func() {
			if c.hasCache {
				defer prepareCacheOK(t)()
			} else {
				defer prepareCacheDetail(t, "")()
			}
			cache, err := NewCache(&Domain{Hostname: "example.com"})
			a.NoError(err, c.pattern)
			a.Implements((*Cache)(nil), cache, c.pattern)
			if c.hasError {
				err := cache.Save(c.ip)
				a.Error(err, c.pattern)
				t.Logf("err: %v", err)
			} else {
				a.NoError(cache.Save(c.ip), c.pattern)
			}
		}()
	}
}

func prepareCacheOK(t *testing.T) func() {
	return prepareCacheDetail(t, `
			ip = '192.168.1.1'
			createdAt = 2019-01-02T12:00:00Z
			updatedAt = 2019-01-02T12:00:00Z
			canUpdatedIn = 2019-01-02T12:05:00Z
	`)
}

func prepareCacheNG(t *testing.T) func() {
	return prepareCacheDetail(t, `
			ip = 'hoge'
			createdAt = 2019-01-02T12:00:00Z
			updatedAt = 2019-01-02T12:00:00Z
			canUpdatedIn = 2019-01-02T12:05:00Z
	`)
}

func prepareCacheDetail(
	t *testing.T,
	content string,
) func() {
	a := assert.New(t)
	tmpDir, err := ioutil.TempDir("", "")
	a.NoError(err)
	if content != "" {
		filename := path.Join(tmpDir, "example.com.cache")
		a.NoError(ioutil.WriteFile(filename, []byte(content), 0600))
	}
	original := cacheDir
	cacheDir = tmpDir
	tt, err := time.Parse(time.RFC3339, "2019-01-02T12:05:00Z")
	a.NoError(err)
	timeNow = func() time.Time { return tt }
	return func() {
		os.RemoveAll(tmpDir)
		cacheDir = original
		timeNow = time.Now
	}
}
