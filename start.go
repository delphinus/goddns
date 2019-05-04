package main

import (
	"github.com/google/logger"
	"golang.org/x/xerrors"
)

// Start starts the main logic
func Start(env *Env, config *Configs, domain *Domain) (Result, error) {
	ip, err := Address(config)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	logger.Infof("detected: %s", ip)
	cache, err := NewCache(env, domain)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	logger.Infof("cache detected: %s", cache.Filename())
	if cache.IsSame(ip) {
		return NoNeedToUpdate(), nil
	}
	if err = cache.CanUpdate(); err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	updater := NewUpdater(env, domain, ip)
	result, err := updater.Update()
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	if result.IsSuccessful() {
		if err := cache.Save(ip); err != nil {
			return nil, xerrors.Errorf(": %w", err)
		}
	}
	return result, nil
}
