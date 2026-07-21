package config

import (
	"fmt"
	"log"
	"net/url"
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
	RedisPort     int    `envconfig:"REDIS_PORT" required:"true"`
	RedisUsername string `envconfig:"REDIS_USERNAME" required:"true"`
	RedisPassword string `envconfig:"REDIS_PASSWORD" required:"true"`

	// Auth
	SecretKey string `envconfig:"SECRET_KEY" required:"true"`
	Expiry    int    `envconfig:"EXPIRY" required:"true"`

	// PRM
	PublicRPM    int `envconfig:"PUBLIC_RPM" default:"5"`
	ProtectedRPM int `envconfig:"PROTECTED_RPM" default:"10"`

	// CORS
	CORSOrigins []string `envconfig:"CORS_ORIGINS" required:"true"`

	//Logging
	LogLevel string `envconfig:"LOG_LEVEL" default:"error"`
}

func (c *Config) PostgresDSN() string {
	user := url.PathEscape(c.PostgresUser)
	password := url.PathEscape(c.PostgresPassword)

	dsn := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(user, password),
		Host:   fmt.Sprintf("%s:%d", c.PostgresHost, c.PostgresPort),
		Path:   c.PostgresDB,
	}

	params := url.Values{}
	params.Add("sslmode", "disable")
	dsn.RawQuery = params.Encode()

	return dsn.String()
}

func (c *Config) RedisAddr() string {
	return fmt.Sprintf("%s:%d", c.RedisHost, c.RedisPort)
}

func Load() *Config {
	once.Do(func() {
		var err error
		err = envconfig.Process("", &cfg)
		if err != nil {
			log.Fatalf("config error: %v", err)
		}
	})
	return &cfg
}
