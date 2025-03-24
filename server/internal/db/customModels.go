package db

type CompData struct {
	Companion struct {
		ContentId string `json:"contentId"`
		ItemId    int32  `json:"itemId"`
		SkinId    int32  `json:"skinId"`
		Species   string `json:"species"`
	} `json:"companion"`
	RemainingGold     int32   `json:"goldLeft"`
	LastRound         int32   `json:"lastRound"`
	Level             int32   `json:"level"`
	Placement         int32   `json:"placement"`
	PlayersEliminated int32   `json:"playersEliminated"`
	Puuid             string  `json:"puuid"`
	TimeEliminated    float64 `json:"timeEliminated"`
	DamageToPlayers   int32   `json:"totalDamageToPlayers"`
	Traits            []struct {
		Name       string `json:"name"`
		NumUnits   int32  `json:"numUnits"`
		Style      int32  `json:"style"`
		TierActive int32  `json:"tierCurrent"`
		TierMax    int32  `json:"tierMax"`
	} `json:"traits"`
	Units []struct {
		CharacterId string   `json:"characterId"`
		ItemNames   []string `json:"itemNames"`
		Rarity      int32    `json:"rarity"`
		Tier        int32    `json:"tier"`
	} `json:"units"`
}
