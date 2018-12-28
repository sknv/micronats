package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	NatsAddr string `envconfig:"NATS_ADDR" default:"nats://127.0.0.1:4222"`
	Secret   string `envconfig:"SECRET" default:"SxvYeR7gxScmY9pjrW63ZZj7KGXDxGn9"`
}

func Parse() *Config {
	cfg := new(Config)
	envconfig.MustProcess("", cfg)
	return cfg
}
