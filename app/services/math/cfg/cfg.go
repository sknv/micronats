package cfg

import (
	"os"

	"github.com/sknv/micronats/app/lib/xflags"
)

type Config struct {
	NatsAddr string `long:"nats-addr" env:"NATS_ADDR" default:"nats://localhost:4222" description:"nats address"`
}

func Parse() *Config {
	cfg := new(Config)
	if _, err := xflags.ParseArgs(os.Args[1:], cfg); err != nil {
		os.Exit(1)
	}
	return cfg
}
