package config

import (
	"github.com/kelseyhightower/envconfig"
)

func Parse[T any](prefix string) (*T, error) {
	var conf T

	if err := envconfig.Process(prefix, &conf); err != nil {
		return nil, err
	}

	return &conf, nil
}
