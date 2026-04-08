package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type rateLimiter struct {
	limiter *rate.Limiter
}

func RateLimit(refillPerSecond float64, requestLimit int) gin.HandlerFunc {
	rl := &rateLimiter{
		limiter: rate.NewLimiter(rate.Limit(refillPerSecond), requestLimit),
	}
	return rl.handle
}

func (rl *rateLimiter) handle(c *gin.Context) {
	if rl.limiter.Allow() {
		c.Next()
	} else {
		c.AbortWithStatus(http.StatusTooManyRequests)
	}
}
