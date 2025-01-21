package services

import (
	"context"
	"time"
)

func (env ServiceEnv) RegionCollectionLoop(ctx context.Context, region string) {
	backoffTicker := time.NewTicker(time.Second * 5)

	for range backoffTicker.C {
		puuids, _ := env.Queries.GetPuuidsWithNullSummonerData(ctx, 100)

		for _, puuid := range puuids {
			env.CollectSummonerDetails(ctx, region, puuid)
		}
	}
}

func (env ServiceEnv) AccountDataCollectionLoop(ctx context.Context, cluster string) {
	backoffTicker := time.NewTicker(time.Second * 5)

	for range backoffTicker.C {
		puuids, _ := env.Queries.GetPuuidsWithNullAccountData(ctx, 100)
		for _, puuid := range puuids {
			env.CollectAccountByPuuid(ctx, cluster, puuid)
		}
	}
}

func (env ServiceEnv) MatchHistoryCollectionLoop(ctx context.Context, cluster string) {
	backoffTicker := time.NewTicker(time.Second * 5)

	for range backoffTicker.C {
		rows, _ := env.Queries.GetOldestMatchHistories(ctx, 1)
		for _, row := range rows {
			env.CollectMatchHistory(ctx, cluster, row.Puuid, time.Now().Add(-time.Hour*24*3))
		}
	}
}
