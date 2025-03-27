package api

import (
	"TFTAnalyticsServer/internal/db"
	"encoding/json"
	"log"
	"net/http"
)

func (controller Api) getSummonerStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	puuid := r.PathValue("puuid")

	stats, err := controller.Queries.GetSummonerStats(ctx, puuid)
	if err != nil {
		log.Printf("Queries.GetSummonerStats failed with err: %v\n", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	// return empty array not null
	if stats == nil {
		stats = []db.GetSummonerStatsRow{}
	}

	w.Header().Set("Content-Type", "application/json")
	bytes, _ := json.Marshal(stats)
	w.Write(bytes)
}
