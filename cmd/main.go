package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

func main() {

	// Logging
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger().Output(zerolog.ConsoleWriter{Out: os.Stdout})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	logger.Info().Msg("Starting the application.")

	// Load cfg
	err := godotenv.Load()
	if err != nil {
		logger.Fatal().Err(err).Msg("Fatal error during parse env file.")
	}
	cfg := config.Load()

	// connect to cache
	rdb, err := redis.RedisConnect(cfg)
	if err != nil {
		logger.Err(err).Str("component", "Redis").Msg("Redis could not connect to db.")
	}
	cache := redis.NewRedisCache(rdb)

	// connect to DB
	DB, err := postgres.NewPostgresDB(cfg)
	if err != nil {
		logger.Fatal().Err(err).Str("component", "Postgres").Msg("Postgres could not connect to db.")
	}

	// prepare user controller
	hasher := hasher.NewBcryptHasher(0)
	userRepo := postgres.NewUserRepository(DB)
	userSVC := service.NewUserService(userRepo, hasher)
	userCtrl := handler.NewUserController(userSVC)

	// prepare post controller
	postRepo := postgres.NewPostRepository(DB)
	postSVC := service.NewPostService(postRepo, cache)
	postCtrl := handler.NewPostController(postSVC)

	// prepare auth controller
	prov := jwt.NewProvider(cfg.SecretKey, time.Duration(cfg.Expiry))
	authSVC := service.NewAuthService(userSVC, hasher, prov)
	authCtrl := handler.NewAuthController(authSVC)

	r := handler.NewRouter(rdb, userCtrl, postCtrl, authCtrl, cfg.SecretKey, time.Duration(cfg.Expiry), cfg.PublicRPM, cfg.ProtectedRPM)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		logger.Info().
			Msg("Starting server on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal().
				Err(err).
				Msg("Critical error during starting server")
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	// Shutdown
	logger.Info().Msg("Shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Fatal().
			Err(err).
			Msg("Server shutdown failed")
	}
	DB.Close()
	rdb.Close()

}
