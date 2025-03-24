package leaderboard

type Rank struct {
	Data                RankData            `json:"rankData"`
	LeaderboardPosition LeaderboardPosition `json:"leaderboardPosition"`
}

type RankData struct {
	Tier         string `json:"tier"`
	Rank         string `json:"rank"`
	LeaguePoints int    `json:"leaguePoints"`
}

type LeaderboardPosition struct {
	Leaderboard string  `json:"leaderboard"`
	Position    int     `json:"position"`
	Top         float32 `json:"top"`
}
