package services

import (
	"TFTAnalyticsServer/db"
	"TFTAnalyticsServer/leaderboard"
	"TFTAnalyticsServer/riot"
	"context"
	"log"
)

var tierScoreMap = map[string]int{
	"CHALLENGER":  2800,
	"GRANDMASTER": 2800,
	"MASTER":      2800,
	"DIAMOND":     2400,
	"EMERALD":     2000,
	"PLATINUM":    1600,
	"GOLD":        1200,
	"SILVER":      800,
	"BRONZE":      400,
	"IRON":        0,
}

var divisionScoreMap = map[string]int{
	"I":   300,
	"II":  200,
	"III": 100,
	"IV":  0,
}

func getRankScore(lp int, tier, division string) int {
	return lp + tierScoreMap[tier] + tierScoreMap[division]
}

func (service Service) GetSummonerRank(ctx context.Context, puuid string) (db.TftRank, error) {
	return service.Queries.GetSummonerRank(ctx, puuid)
}

func (service Service) CollectSummonerRank(ctx context.Context, region, summonerId string) error {
	rankEntry, err := riot.GetRank(region, summonerId)
	if err != nil {
		return err
	}

	return service.Queries.UpsertRank(ctx, db.UpsertRankParams{
		SummonerID:   rankEntry.SummonerId,
		Tier:         rankEntry.Tier,
		Rank:         rankEntry.Rank,
		LeaguePoints: rankEntry.LeaguePoints,
		Wins:         rankEntry.Wins,
		Losses:       rankEntry.Losses,
	})
}

func (service Service) CollectRankEntries(ctx context.Context, region, tier, division string) error {
	page := 0
	for {
		page++
		rankEntries, err := riot.GetRankEntries(region, tier, division, page)
		if err != nil {
			return err
		}
		if len(rankEntries) == 0 {
			break
		}

		// create datastore params
		//upsertPuuidParams := make([]db.BatchUpsertPuuidParams, len(rankEntries))
		upsertRankEntryParams := make([]db.BatchUpsertRankParams, len(rankEntries))
		setLeaderboardScoresParams := make([]leaderboard.SetLeaderboardScoresParams, len(rankEntries))

		for i, rankEntry := range rankEntries {
			//upsertPuuidParams[i] = db.BatchUpsertPuuidParams{
			//	Puuid:  rankEntry.Puuid,
			//	Region: region,
			//}

			setLeaderboardScoresParams[i] = leaderboard.SetLeaderboardScoresParams{
				Id:    rankEntry.SummonerId,
				Score: getRankScore(int(rankEntry.LeaguePoints), rankEntry.Tier, rankEntry.Rank),
			}

			upsertRankEntryParams[i] = db.BatchUpsertRankParams{
				SummonerID:   rankEntry.SummonerId,
				Tier:         rankEntry.Tier,
				Rank:         rankEntry.Rank,
				LeaguePoints: rankEntry.LeaguePoints,
				Wins:         rankEntry.Wins,
				Losses:       rankEntry.Losses,
			}
		}

		//service.Queries.BatchUpsertPuuid(ctx, upsertPuuidParams).Exec(nil)
		_ = service.Leaderboard.SetLeaderboardScores(ctx, setLeaderboardScoresParams...)
		service.Queries.BatchUpsertRank(ctx, upsertRankEntryParams).Exec(nil)

		log.Printf("Rank entries collected for %v %v %v page %v\n", region, tier, division, page)
	}

	return nil
}

func (service Service) CollectApexRankEntries(ctx context.Context, region, tier string) error {
	rankPage, err := riot.GetApexRankPage(region, tier)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	// create datastore params
	setLeaderboardScoresParams := make([]leaderboard.SetLeaderboardScoresParams, len(rankPage.Entries))
	upsertRankEntryParams := make([]db.BatchUpsertRankParams, len(rankPage.Entries))

	for i, rankEntry := range rankPage.Entries {
		setLeaderboardScoresParams[i] = leaderboard.SetLeaderboardScoresParams{
			Id:    rankEntry.SummonerId,
			Score: getRankScore(int(rankEntry.LeaguePoints), rankPage.Tier, rankEntry.Rank),
		}

		upsertRankEntryParams[i] = db.BatchUpsertRankParams{
			SummonerID:   rankEntry.SummonerId,
			Tier:         rankPage.Tier,
			Rank:         rankEntry.Rank,
			LeaguePoints: rankEntry.LeaguePoints,
			Wins:         rankEntry.Wins,
			Losses:       rankEntry.Losses,
		}
	}

	// batch upsert rank entries
	_ = service.Leaderboard.SetLeaderboardScores(ctx, setLeaderboardScoresParams...)
	service.Queries.BatchUpsertRank(ctx, upsertRankEntryParams).Exec(nil)

	log.Printf("Rank entries collected for %v %v\n", region, tier)
	return nil
}
