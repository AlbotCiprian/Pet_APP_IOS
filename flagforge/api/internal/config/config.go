package config

import (
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config defines the runtime configuration for the API service.
type Config struct {
	Port        int           `envconfig:"API_PORT" default:"8080"`
	PostgresDSN string        `envconfig:"POSTGRES_DSN" required:"true"`
	RedisAddr   string        `envconfig:"REDIS_ADDR" default:"redis:6379"`
	JWTSecret   string        `envconfig:"JWT_SECRET" default:"change_me"`
	Migrations  string        `envconfig:"MIGRATIONS_DIR" default:"../deploy/migrations"`
	ShutdownTTL time.Duration `envconfig:"SHUTDOWN_TTL" default:"10s"`
}

// Load parses environment variables into Config.
func Load() Config {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	return cfg
}
