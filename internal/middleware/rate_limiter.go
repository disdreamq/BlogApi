package middleware

import (
	"net/http"
	"strconv"

	"github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"
)

type RateLimitMiddleware struct {
	limiter *redis_rate.Limiter
	rpm     int
}

func NewRateLimitMiddleware(rdb *redis.Client, rpm int) *RateLimitMiddleware {
	return &RateLimitMiddleware{
		limiter: redis_rate.NewLimiter(rdb),
		rpm:     rpm,
	}
}

func (l *RateLimitMiddleware) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		key := "rate_limit:" + ip

		res, err := l.limiter.Allow(r.Context(), key, redis_rate.PerMinute(l.rpm))
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		if res.Allowed == 0 {
			w.Header().Set("Retry-After", strconv.Itoa((int(res.RetryAfter.Seconds()))))
			http.Error(w, "too many requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
