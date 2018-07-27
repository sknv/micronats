package cfg

import (
	"os"

	flags "github.com/jessevdk/go-flags"
)

type Config struct {
	Addr    string `long:"math-addr" env:"MATH_ADDR" default:"localhost:8081" description:"math service address"`
	NatsURL string `long:"nats-url" env:"NATS_URL" default:"nats://localhost:4222" description:"nats server url"`
}

func Parse() *Config {
	config := &Config{}
	flagParser := flags.NewParser(config, flags.Default)
	if _, err := flagParser.ParseArgs(os.Args[1:]); err != nil {
		os.Exit(1)
	}
	return config
}
