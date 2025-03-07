package riot

import (
	"TFTAnalyticsServer/riot/limiter"
	"fmt"
	"time"
)

type RateLimiter map[string]*limiter.Limiter

func NewRateLimiter(rateDuration time.Duration) map[string]*limiter.Limiter {
	limiters := make(map[string]*limiter.Limiter)

	for cluster := range ClusterToRegions {
		limiters[cluster] = limiter.New(rateDuration)
	}
	for region := range RegionToCluster {
		limiters[region] = limiter.New(rateDuration)
	}

	return limiters
}

func (limters RateLimiter) Wait(server string, priority int) error {
	limiter, ok := limters[server]
	if !ok {
		return fmt.Errorf("no limiter found for server: %v", server)
	}

	<-limiter.Wait(priority)

	return nil
}
