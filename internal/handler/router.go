package handler

import (
	"time"

	"github.com/disdreamq/BlogApi/internal/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
)

// Добавить rate limmiter и почти приехал
func NewRouter(
	rdb *redis.Client,
	userCtrl *UserController,
	postCtrl *PostController,
	authCtrl *AuthController,
	secret string,
	expiry time.Duration,
	PublicRPM int,
	ProtectedPRM int,
) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RecoveryMiddleware)
	r.Use(middleware.LoggingMiddleware)

	// Public routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.NewRateLimitMiddleware(rdb, PublicRPM).Limit)
		r.Post("/register", userCtrl.Create)
		r.Post("/login", authCtrl.Login)
	})

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.NewRateLimitMiddleware(rdb, ProtectedPRM).Limit)
		// r.Use(middleware.NewAuthMiddleware(jwt.NewProvider(secret, expiry)).Authenticate)

		r.Route("/users", func(r chi.Router) {
			r.Get("/{userID}", userCtrl.GetByID)
			r.Get("/{email}", userCtrl.GetByEmail)
			r.Put("/{userID}", userCtrl.Update)
			r.Delete("/{userID}", userCtrl.Delete)
		})
		r.Route("/posts", func(r chi.Router) {
			r.Post("/", postCtrl.Create)
			r.Get("/{postID}", postCtrl.GetByID)
			r.Get("/{title}", postCtrl.GetByTitle)
			r.Put("/{postID}", postCtrl.Update)
			r.Delete("/{postID}", postCtrl.Delete)
		})
	})
	return r
}
