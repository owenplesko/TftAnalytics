package riot

import (
	"context"
	"fmt"
	"time"

	"github.com/owenplesko/TftAnalytics/pkg/riot/scheduler"
)

type RateLimiter map[string]*scheduler.Scheduler

func newRateLimiter(rateDuration time.Duration) map[string]*scheduler.Scheduler {
	limiters := make(map[string]*scheduler.Scheduler)

	for cluster := range ClusterToRegions {
		limiters[cluster] = scheduler.New(rateDuration)
	}
	for region := range RegionToCluster {
		limiters[region] = scheduler.New(rateDuration)
	}

	return limiters
}

func (limters RateLimiter) Wait(ctx context.Context, server string) error {
	limiter, ok := limters[server]
	if !ok {
		return fmt.Errorf("no limiter found for server: %v", server)
	}

	<-limiter.Wait(ctx)

	return nil
}
