package main

import (
	"fmt"

	"golang.org/x/xerrors"
)

func Start(domain *Domain) (Result, error) {
	ip, err := Address()
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	logger.Notice(fmt.Sprintf("detected: %s", ip))
	cache, err := NewCache(domain)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	logger.Notice(fmt.Sprintf("cache detected: %s", cache.Filename()))
	if cache.IsSame(ip) {
		return NoNeedToUpdate(), nil
	}
	if err := cache.CanUpdate(); err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	updater := NewUpdater(domain, ip)
	result, err := updater.Update()
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	if result.IsSuccessful() {
		if err := cache.Save(ip); err != nil {
			return nil, xerrors.Errorf(": %w", err)
		}
		return result, nil
	} else if result.IsCritical() {
		return nil, xerrors.New(fmt.Sprintf("Updater returned a critical error: %s", result))
	}
	return result, nil
}
