package middleware

import (
	"awesomeProject/internal/model"
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func RateLimitMiddleware(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		val, exists := c.Get("user_key")
		if !exists {
			c.Next()
			return
		}

		userKey := val.(*model.UserAPIKey)
		ctx := context.Background()

		key := fmt.Sprintf("ratelimit:%d:%d", userKey.ID, time.Now().Unix()/60)

		count, _ := rdb.Incr(ctx, key).Result()
		if count == 1 {
			rdb.Expire(ctx, key, 2*time.Minute)
		}

		limit := userKey.RateLimit
		if limit <= 0 {
			limit = 60
		}

		c.Header("X-RateLimit-Limit", strconv.Itoa(limit))

		if count > int64(limit) {
			c.Header("X-RateLimit-Remaining", "0")
			c.Header("Retry-After", "60")
			c.AbortWithStatusJSON(429, gin.H{
				"error": gin.H{"message": "rate limit exceeded", "type": "rate_limit_error", "code": "rate_limit_exceeded"},
			})
			return
		}

		remaining := int64(limit) - count
		if remaining < 0 {
			remaining = 0
		}
		c.Header("X-RateLimit-Remaining", strconv.FormatInt(remaining, 10))
		c.Next()
	}
}
