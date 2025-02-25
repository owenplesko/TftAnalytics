package riot

import (
	"TFTAnalyticsServer/riot/limiter"
	"time"
)

type RateLimiter map[string]*limiter.Limiter

func NewRateLimiter(rateDuration time.Duration) map[string]*limiter.Limiter {
	limiters := make(map[string]*limiter.Limiter)

	for cluster := range ClusterToRegions {
		limiters[cluster] = limiter.New(rateDuration)
	}

	return limiters
}
