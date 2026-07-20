// @title           Blog API
// @version         1.0
// @description     REST API для блога с авторизацией, управлением пользователями и постами
// @host            localhost:8080
// @BasePath        /
// @schemes         http
// @securityDefinitions.apikey BearerAuth
// @in                        header
// @name                      Authorization
// @description               Введите токен в формате: Bearer <ваш_токен>

package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/disdreamq/BlogApi/docs"
	"github.com/fatih/color"

	"github.com/disdreamq/BlogApi/config"
	"github.com/disdreamq/BlogApi/internal/handler"
	"github.com/disdreamq/BlogApi/internal/infra/hasher"
	"github.com/disdreamq/BlogApi/internal/infra/jwt"
	"github.com/disdreamq/BlogApi/internal/repository/postgres"
	"github.com/disdreamq/BlogApi/internal/repository/redis"
	"github.com/disdreamq/BlogApi/internal/service"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	httpSwagger "github.com/swaggo/http-swagger"
)

func printBanner(cfg *config.Config) {
	color.NoColor = false // Force colors

	banner := `
██████╗ ██╗      ██████╗  ██████╗        █████╗ ██████╗ ██╗
██╔══██╗██║     ██╔═══██╗██╔════╝       ██╔══██╗██╔══██╗██║
██████╔╝██║     ██║   ██║██║  ███╗█████╗███████║██████╔╝██║
██╔══██╗██║     ██║   ██║██║   ██║╚════╝██╔══██║██╔═══╝ ██║
██████╔╝███████╗╚██████╔╝╚██████╔╝      ██║  ██║██║     ██║
╚═════╝ ╚══════╝ ╚═════╝  ╚═════╝       ╚═╝  ╚═╝╚═╝     ╚═╝ `

	red := color.New(color.FgRed, color.Bold)
	green := color.New(color.FgGreen, color.Bold)
	yellow := color.New(color.FgYellow, color.Bold)
	magenta := color.New(color.FgMagenta, color.Bold)
	white := color.New(color.FgWhite, color.Bold)

	red.Println(banner)
	println()

	port := "8080"

	yellow.Println("Configuration:")
	white.Printf("   Port:        %s\n", port)
	white.Printf("   JWT Expiry:  %s\n", time.Duration(cfg.Expiry).String())
	white.Printf("   Rate Limit:  Public=%d RPM, Protected=%d RPM\n", cfg.PublicRPM, cfg.ProtectedRPM)
	println()

	green.Println("Services:")
	green.Println("   ✓ PostgreSQL connected")
	green.Println("   ✓ Redis connected")
	green.Println("   ✓ JWT Auth middleware enabled")
	green.Println("   ✓ Rate limiting enabled")
	green.Println("   ✓ Logging middleware enabled")
	green.Println("   ✓ Recovery middleware enabled")
	println()

	yellow.Println("Authentication:")
	white.Println("   Public endpoints:  POST /register, POST /login")
	white.Println("   Protected:         PUT/DELETE users, POST/GET/PUT/DELETE posts")
	white.Println("   Header:            Authorization: Bearer <token>")
	println()

	magenta.Println("Endpoints:")
	white.Printf("   API Base:     http://localhost:%s/\n", port)
	white.Printf("   Swagger:      http://localhost:%s/swagger/\n", port)
	white.Printf("   Swagger JSON: http://localhost:%s/swagger/doc.json\n", port)
	println()

	yellow.Println("🚀 Server is running! Press CTRL+C to stop.")
}

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
	authSVC := service.NewAuthService(userRepo, hasher, prov)
	authCtrl := handler.NewAuthController(authSVC)

	r := handler.NewRouter(rdb, userCtrl, postCtrl, authCtrl, cfg.SecretKey, time.Duration(cfg.Expiry), cfg.PublicRPM, cfg.ProtectedRPM)

	// Swagger
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), // явно указываем URL документации
	))

	// Print startup banner

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal().
				Err(err).
				Msg("Critical error during starting server")
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	printBanner(cfg)
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
