package driver

import (
	"log"

	"go.risoftinc.com/goresponse"
	"go.risoftinc.com/xarch/config"
)

func ResponseManager(cfg config.ResponseManager) (*goresponse.ResponseConfig, error) {
	return goresponse.LoadConfig(goresponse.ConfigSource{
		Method: cfg.Method,
		Path:   cfg.Path,
	})
}

func ResponseManagerAsync(cfg config.ResponseManager) (*goresponse.AsyncConfigManager, error) {
	asyncManager := goresponse.NewAsyncConfigManager(goresponse.ConfigSource{
		Method: cfg.Method,
		Path:   cfg.Path,
	}, cfg.Interval)

	asyncManager.AddCallback(func(oldConfig, newConfig *goresponse.ResponseConfig) {
		log.Printf("Config response manager updated!")
	})

	if err := asyncManager.Start(); err != nil {
		return nil, err
	}

	return asyncManager, nil
}
