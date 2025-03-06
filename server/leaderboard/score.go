package leaderboard

import "TFTAnalyticsServer/types"

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

func getRankScore(rank types.RankData) int {
	rankScore := tierScoreMap[rank.Tier]

	if rankScore < tierScoreMap["MASTER"] {
		rankScore += divisionScoreMap[rank.Rank]
	}

	rankScore += int(rank.LeaguePoints)

	return int(rank.LeaguePoints) + tierScoreMap[rank.Tier] + tierScoreMap[rank.Rank]
}
