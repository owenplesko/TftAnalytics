package schedule

import (
	"context"

	"github.com/go-co-op/gocron/v2"
	"github.com/owenplesko/TftAnalytics/internal/services"
)

type Schedule struct {
	scheduler gocron.Scheduler
	service   *services.Service
}

func New(service *services.Service) *Schedule {
	ctx := context.Background()

	scheduler, err := gocron.NewScheduler()
	if err != nil {
		panic(err)
	}

	_, _ = scheduler.NewJob(
		gocron.DailyJob(1, gocron.NewAtTimes(gocron.NewAtTime(0, 0, 0))),
		gocron.NewTask(service.RefreshUnitPlacement, ctx),
		gocron.WithSingletonMode(gocron.LimitModeReschedule),
	)

	return &Schedule{
		scheduler: scheduler,
		service:   service,
	}
}

func (s *Schedule) Start() {
	s.scheduler.Start()
}

func (s *Schedule) Shutdown() {
	s.scheduler.Shutdown()
}
