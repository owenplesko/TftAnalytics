package schedule

import (
	"github.com/owenplesko/TftAnalytics/internal/services"
	"github.com/robfig/cron/v3"
)

type Schedule struct {
	c       *cron.Cron
	service *services.Service
}

func New(service *services.Service) *Schedule {
	c := cron.New()

	return &Schedule{
		c:       c,
		service: service,
	}
}

func (s *Schedule) Start() {
	s.c.Start()
}

func (s *Schedule) Stop() {
	s.c.Stop()
}
