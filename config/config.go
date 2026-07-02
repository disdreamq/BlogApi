package config

import (
	"fmt"
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
	RedisPort     string `envconfig:"REDIS_PORT" required:"true"`
	RedisUsername string `envconfig:"REDIS_USERNAME" required:"true"`
	RedisPassword string `envconfig:"REDIS_PASSWORD" required:"true"`

	// Auth
	SecretKey string `envconfig:"SECRET_KEY" required:"true"`
	Alghoritm string `envconfig:"ALGHORITM" required:"true"`

	// CORS
	CORSOrigins []string `envconfig:"CORS_ORIGINS" required:"true"`

	//Logging
	LogLevel string `envconfig:"LOG_LEVEL" default:"error"`
}

func (c *Config) PostgresDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.PostgresHost, c.PostgresPort, c.PostgresUser, c.PostgresPassword, c.PostgresDB,
	)
}

func (c *Config) RedisAddr() string {
	return fmt.Sprintf("%s:%d", c.RedisHost, c.RedisPort)
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
