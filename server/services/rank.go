package services

import (
	"TFTAnalyticsServer/leaderboard"
	"TFTAnalyticsServer/types"
	"context"
	"log"
)

func (service Service) GetSummonerRank(ctx context.Context, region string, summonerId string) (types.Rank, error) {
	return service.Leaderboard.GetRank(ctx, region, summonerId)
}

func (service Service) CollectSummonerRank(ctx context.Context, region, summonerId string) error {
	rankEntry, err := service.Riot.GetRank(region, summonerId)
	if err != nil {
		return err
	}

	setRankParam := leaderboard.SetRankParams{
		SummonerId: rankEntry.SummonerId,
		RankData: types.RankData{
			Tier:         rankEntry.Tier,
			Rank:         rankEntry.Rank,
			LeaguePoints: int(rankEntry.LeaguePoints),
		}}

	err = service.Leaderboard.SetRank(ctx, region, setRankParam)

	return err
}

func (service Service) CollectRankEntries(ctx context.Context, region, tier, division string) error {
	page := 0
	for {
		page++
		rankEntries, err := service.Riot.GetRankEntries(region, tier, division, page)
		if err != nil {
			return err
		}
		if len(rankEntries) == 0 {
			break
		}

		// create datastore params
		//upsertPuuidParams := make([]db.BatchUpsertPuuidParams, len(rankEntries))
		setRankParams := make([]leaderboard.SetRankParams, len(rankEntries))

		for i, rankEntry := range rankEntries {
			//upsertPuuidParams[i] = db.BatchUpsertPuuidParams{
			//	Puuid:  rankEntry.Puuid,
			//	Region: region,
			//}

			setRankParams[i] = leaderboard.SetRankParams{
				SummonerId: rankEntry.SummonerId,
				RankData: types.RankData{
					Tier:         rankEntry.Tier,
					Rank:         rankEntry.Rank,
					LeaguePoints: int(rankEntry.LeaguePoints),
				}}
		}

		//service.Queries.BatchUpsertPuuid(ctx, upsertPuuidParams).Exec(nil)
		_ = service.Leaderboard.SetRank(ctx, region, setRankParams...)

		log.Printf("Rank entries collected for %v %v %v page %v\n", region, tier, division, page)
	}

	return nil
}

func (service Service) CollectApexRankEntries(ctx context.Context, region, tier string) error {
	rankPage, err := service.Riot.GetApexRankPage(region, tier)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	// create datastore params
	setRankParams := make([]leaderboard.SetRankParams, len(rankPage.Entries))

	for i, rankEntry := range rankPage.Entries {
		setRankParams[i] = leaderboard.SetRankParams{
			SummonerId: rankEntry.SummonerId,
			RankData: types.RankData{
				Tier:         tier,
				Rank:         rankEntry.Rank,
				LeaguePoints: int(rankEntry.LeaguePoints),
			},
		}
	}

	// batch upsert rank entries
	_ = service.Leaderboard.SetRank(ctx, region, setRankParams...)

	log.Printf("Rank entries collected for %v %v\n", region, tier)
	return nil
}
