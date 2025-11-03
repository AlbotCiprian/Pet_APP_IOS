package config

import (
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port        int           `envconfig:"WORKER_PORT" default:"8090"`
	RedisAddr   string        `envconfig:"REDIS_ADDR" default:"redis:6379"`
	ShutdownTTL time.Duration `envconfig:"SHUTDOWN_TTL" default:"10s"`
}

func Load() Config {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	return cfg
}
