package ratelimiter

import (
	"time"

	"github.com/callicoder/go-commons/redis"
	"github.com/go-redis/redis_rate/v7"
)

// RateLimiter should implement
type RateLimiter interface {
	RateLimit(key string, maxAttempts int, windowInSeconds int) (int, time.Duration, bool)
}

// RedisRateLimiter implements ratelimiter using redis db to store key values
type RedisRateLimiter struct {
	limiter *redis_rate.Limiter
}

// NewRateLimiter gives a ratelimiter with given redis pool and configuration
func NewRateLimiter(redisClient redis.Client) *RedisRateLimiter {
	var limiter *redis_rate.Limiter
	limiter = redis_rate.NewLimiter(redisClient)

	return &RedisRateLimiter{
		limiter: limiter,
	}
}

func (redisRL *RedisRateLimiter) RateLimit(key string, maxAttempts int, windowInSeconds int) (int, time.Duration, bool) {
	count, delay, allow := redisRL.limiter.Allow(key, int64(maxAttempts), time.Duration(windowInSeconds)*time.Second)
	return int(count), delay, allow
}
