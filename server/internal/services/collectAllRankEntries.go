package services

import (
	"context"
	"fmt"
	"log"

	"github.com/owenplesko/TftAnalytics/pkg/dedupe"
	"github.com/owenplesko/TftAnalytics/pkg/riot"
)

func (service *Service) CollectAllRankEntries(ctx context.Context, region string) error {
	return dedupe.Run(service.deduplicator, allRankEntriesTask{
		service: service,
		ctx:     ctx,
		region:  region,
	}).Await()
}

type allRankEntriesTask struct {
	service *Service
	ctx     context.Context
	region  string
}

func (task allRankEntriesTask) ID() string {
	return fmt.Sprintf("ALL_RANK_ENTRIES_%s", task.region)
}

// TODO: error handling can be better here
func (task allRankEntriesTask) Run() error {
	for _, tier := range riot.ApexTiers {
		err := task.service.CollectApexRankEntries(task.ctx, task.region, tier)
		if err != nil {
			log.Printf("CollectApexRankEntries failed with err: %v", err)
		}
	}
	for _, tier := range riot.Tiers {
		for _, division := range riot.Divisions {
			err := task.service.CollectRankEntries(task.ctx, task.region, tier, division)
			if err != nil {
				log.Printf("CollectRankEntries failed with err: %v", err)
			}
		}
	}

	return nil
}
