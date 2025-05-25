package schedule

import (
	"context"

	"github.com/owenplesko/TftAnalytics/internal/services"
	"github.com/owenplesko/TftAnalytics/pkg/riot"
	"github.com/robfig/cron/v3"
)

type Schedule struct {
	c       *cron.Cron
	service *services.Service
}

func New(service *services.Service) *Schedule {
	ctx := context.Background()

	c := cron.New()
	c.AddFunc("0 0 * * *", func() { service.RefreshUnitPlacement(ctx) })

	for region := range riot.RegionToCluster {
		c.AddFunc("*/20 * * * *", func() { service.CollectAllRankEntries(ctx, region) })
	}

	return &Schedule{
		c:       c,
		service: service,
	}
}

func (s *Schedule) Start() {
	s.c.Start()
}

func (s *Schedule) Shutdown() {
	s.c.Stop()
}
