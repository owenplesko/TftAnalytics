package crawler

import (
	"context"
	"log"
	"time"

	"github.com/owenplesko/TftAnalytics/internal/services"
)

type MatchCrawler struct {
	service *services.Service
	region  string
}

func NewMatchCrawler(service *services.Service, region string) *MatchCrawler {
	return &MatchCrawler{
		service: service,
		region:  region,
	}
}

func (mc *MatchCrawler) Begin() {
	ctx := context.Background()

	backoffTime := time.NewTicker(time.Second * 5)

	for range backoffTime.C {
		summoners, err := mc.service.GetOldestMatchHistories(ctx, mc.region)
		if err != nil {
			log.Printf("in MatchCollectionLoop Queries.GetOldestMatchesAfter failed with err: %v", err)
			continue
		}

		for _, summoner := range summoners {
			err = mc.service.CollectMatchHistory(ctx, mc.region, summoner.Puuid, summoner.MatchesBeforeTimestamp.Time)
			if err != nil {
				log.Printf("error in MatchCollectionLoop service.CollectMatchHistory failed with err: %v", err)
			}
		}
	}
}
