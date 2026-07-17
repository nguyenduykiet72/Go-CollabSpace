package middleware

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"

	"Go-CollabSpace/pkg/httpx"
)

// RateLimiter is a small Gin middleware backed by go-redis_rate (GCRA / token
// bucket). It is intentionally simple: pass a key extractor and a per-route
// limit. Returns 429 with a Retry-After header when exhausted.
//
// Example:
//
//	limiter := middleware.NewRateLimiter(redisClient)
//	router.POST("/auth/login",
//	    limiter.Limit("auth:login", middleware.LimitByIP, redis_rate.PerMinute(10)),
//	    handler)
type RateLimiter struct {
	rl *redis_rate.Limiter
}

func NewRateLimiter(client *redis.Client) *RateLimiter {
	return &RateLimiter{rl: redis_rate.NewLimiter(client)}
}

type KeyFunc func(*gin.Context) string

func LimitByIP(c *gin.Context) string {
	return c.ClientIP()
}

func (r *RateLimiter) Limit(namespace string, keyFn KeyFunc, limit redis_rate.Limit) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := keyFn(c)
		if key == "" {
			c.Next()
			return
		}

		bucket := namespace + ":" + key
		res, err := r.rl.Allow(c.Request.Context(), bucket, limit)
		if err != nil {
			if errors.Is(err, redis.ErrClosed) {
				// Redis closed during shutdown; fail open to avoid blocking
				// in-flight requests.
				c.Next()
				return
			}
			// Treat Redis errors as fail-open too: rate limiting must not
			// become an availability hazard.
			c.Next()
			return
		}

		c.Header("X-RateLimit-Limit", strconv.Itoa(limit.Rate))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(res.Remaining))

		if res.Allowed <= 0 {
			retryAfter := int(res.RetryAfter / time.Second)
			if retryAfter < 1 {
				retryAfter = 1
			}
			c.Header("Retry-After", strconv.Itoa(retryAfter))
			httpx.WriteJSON(c, http.StatusTooManyRequests, nil, "Too many requests, please retry later")
			c.Abort()
			return
		}

		c.Next()
	}
}
