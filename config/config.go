package config

import (
	"log"
	"sync"

	"github.com/kelseyhightower/envconfig"
)

var (
	cfg  Config
	once sync.Once
)

type Config struct {
	// PostgreSQL
	PostgresDB       string `envconfig:"POSTGRES_DB" required:"true"`
	PostgresHost     string `envconfig:"POSTGRES_HOST" required:"true"`
	PostgresPort     int    `envconfig:"POSTGRES_PORT" required:"true"`
	PostgresUser     string `envconfig:"POSTGRES_USER" required:"true"`
	PostgresPassword string `envconfig:"POSTGRES_PASSWORD" required:"true"`

	// Redis
	RedisDatabase int    `envconfig:"REDIS_DATABASE" required:"true"`
	RedisHost     string `envconfig:"REDIS_HOST" required:"true"`
	RedisUsername string `envconfig:"REDIS_PORT" required:"true"`
	RedisPassword string `envconfig:"REDIS_USERNAME" required:"true"`

	// Auth
	SecretKey string `envconfig:"SECRET_KEY" required:"true"`
	Alghoritm string `envconfig:"ALGHORITM" required:"true"`

	// CORS
	CORSOrigins []string `envconfig:"CORS_ORIGINS" required:"true"`

	//Logging
	LogLevel string `envconfig:"LOG_LEVEL" default:"error"`
}

func Load() (*Config, error) {
	once.Do(func() {
		var err error
		err = envconfig.Process("", &cfg)
		if err != nil {
			log.Fatalf("config error: %v", err)
		}
	})
	return &cfg, nil
}
