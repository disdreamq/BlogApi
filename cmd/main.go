package main

import (
	"net/http"
	"os"
	"time"

	"github.com/disdreamq/BlogApi/config"
	"github.com/disdreamq/BlogApi/internal/handler"
	"github.com/disdreamq/BlogApi/internal/infra/hasher"
	"github.com/disdreamq/BlogApi/internal/infra/jwt"
	"github.com/disdreamq/BlogApi/internal/repository/postgres"
	"github.com/disdreamq/BlogApi/internal/repository/redis"
	"github.com/disdreamq/BlogApi/internal/service"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

// логирование и запуск.
func main() {

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger().Output(zerolog.ConsoleWriter{Out: os.Stdout})
	logger.Info().Msg("Starting the application.")

	err := godotenv.Load()
	if err != nil {
		logger.Fatal().Err(err).Msg("Fatal error during parse env file.")
	}
	cfg := config.Load()
	rdb, err := redis.RedisConnect(cfg)
	if err != nil {
		logger.Err(err).Str("component", "Redis").Msg("Redis could not connect to db.")
	}
	cache := redis.NewRedisCache(rdb)

	DB, err := postgres.NewPostgresDB(cfg)
	if err != nil {
		logger.Fatal().Err(err).Str("component", "Postgres").Msg("Postgres could not connect to db.")
	}

	hasher := hasher.NewBcryptHasher(0)
	userRepo := postgres.NewUserRepository(DB)
	userSVC := service.NewUserService(userRepo, hasher)
	userCtrl := handler.NewUserController(userSVC)

	postRepo := postgres.NewPostRepository(DB)
	postSVC := service.NewPostService(postRepo, cache)
	postCtrl := handler.NewPostController(postSVC)

	prov := jwt.NewProvider(cfg.SecretKey, time.Duration(cfg.Expiry))
	authSVC := service.NewAuthService(userSVC, hasher, prov)
	authCtrl := handler.NewAuthController(authSVC)

	r := handler.NewRouter(rdb, userCtrl, postCtrl, authCtrl, cfg.SecretKey, time.Duration(cfg.Expiry), cfg.PublicRPM, cfg.ProtectedRPM)

	logger.Info().Msg("Server started.")
	http.ListenAndServe(":8080", r)
}
