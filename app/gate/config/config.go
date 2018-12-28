package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	GateAddr     string `envconfig:"GATE_ADDR" default:":8080"`
	NatsAddr     string `envconfig:"NATS_ADDR" default:"nats://127.0.0.1:4222"`
	ResourcesDir string `envconfig:"RESOURCES_DIR" default:"resources"`
}

func Parse() *Config {
	cfg := new(Config)
	envconfig.MustProcess("", cfg)
	return cfg
}

func (c *Config) SchemaDir() string {
	return fmt.Sprintf("%s/schema", c.ResourcesDir)
}
