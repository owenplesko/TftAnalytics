package services

import (
	"TFTAnalyticsServer/db"
	"TFTAnalyticsServer/riot"
	"context"
	"time"
)

func (env Service) SummonerDataCollectionLoop(ctx context.Context, region string) {
	backoffTicker := time.NewTicker(time.Second * 5)

	for range backoffTicker.C {
		puuids, _ := env.Queries.GetPuuidsWithNullSummonerData(ctx, db.GetPuuidsWithNullSummonerDataParams{
			Limit:  100,
			Region: region,
		})

		for _, puuid := range puuids {
			env.CollectSummonerDetails(ctx, region, puuid)
		}
	}
}

func (env Service) SummonerRegionCollectionLoop(ctx context.Context) {
	backoffTicker := time.NewTicker(time.Second * 5)

	for range backoffTicker.C {
		puuids, _ := env.Queries.GetPuuidsWithNullRegion(ctx, 100)

		for _, puuid := range puuids {
			env.CollectSummonerRegion(ctx, puuid)
		}
	}
}

func (env Service) AccountDataCollectionLoop(ctx context.Context, cluster string) {
	backoffTicker := time.NewTicker(time.Second * 5)

	for range backoffTicker.C {
		puuids, _ := env.Queries.GetPuuidsWithNullAccountData(ctx, 100)
		for _, puuid := range puuids {
			env.CollectAccountByPuuid(ctx, cluster, puuid)
		}
	}
}

func (env Service) MatchHistoryCollectionLoop(ctx context.Context, cluster string) {
	backoffTicker := time.NewTicker(time.Second * 5)

	for range backoffTicker.C {
		for _, region := range riot.ClusterToRegions[cluster] {
			rows, _ := env.Queries.GetOldestMatchHistories(ctx, db.GetOldestMatchHistoriesParams{
				Region: region,
				Limit:  100,
			})

			for _, row := range rows {
				env.CollectMatchHistory(ctx, cluster, row.Puuid, time.Now().Add(-time.Hour*24*3))
			}
		}

	}
}

func (env Service) SummonerRankCollectionLoop(ctx context.Context, region string) {
	//for {
	//	for _, tier := range riot.Tiers {
	//
	//	}
	//}
}
