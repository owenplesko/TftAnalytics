package services

import (
	"context"
	"log"
	"time"

	"github.com/owenplesko/TftAnalytics/internal/db"
)

func (service *Service) MatchCollectionLoop(ctx context.Context, region string) {
	backoffTime := time.NewTicker(time.Second * 5)

	for range backoffTime.C {
		summoners, err := service.queries.GetOldestMatchesAfter(ctx, db.GetOldestMatchesAfterParams{
			Limit:  100,
			Region: region,
		})
		if err != nil {
			log.Printf("in MatchCollectionLoop Queries.GetOldestMatchesAfter failed with err: %v", err)
			continue
		}

		for _, summoner := range summoners {
			err = service.CollectMatchHistory(ctx, region, summoner.Puuid, summoner.MatchesBeforeTimestamp.Time)
			if err != nil {
				log.Printf("error in MatchCollectionLoop service.CollectMatchHistory failed with err: %v", err)
			}
		}
	}
}
