package cfg

import (
	"os"

	"github.com/sknv/micronats/app/lib/xflags"
)

type Config struct {
	Addr       string `long:"rest-addr" env:"REST_ADDR" default:":8080" description:"rest api address"`
	NatsAddr   string `long:"nats-addr" env:"NATS_ADDR" default:"nats://127.0.0.1:4222" description:"nats address"`
	ConsulAddr string `long:"consul-addr" env:"CONSUL_ADDR" default:"127.0.0.1:8500" description:"consul address"`
}

func Parse() *Config {
	cfg := new(Config)
	if _, err := xflags.ParseArgs(os.Args[1:], cfg); err != nil {
		os.Exit(1)
	}
	return cfg
}
