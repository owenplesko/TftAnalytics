package api

import (
	"TFTAnalyticsServer/internal/db"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func (controller Api) getSummonerStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	puuid := r.PathValue("puuid")
	setNumber, err := strconv.Atoi(r.PathValue("setNumber"))
	if err != nil {
		http.Error(w, "failed to parse setNumber to int", http.StatusBadRequest)
		return
	}

	stats, err := controller.Queries.GetSummonerStats(ctx, db.GetSummonerStatsParams{
		SummonerPuuid: puuid,
		SetNumber:     int32(setNumber),
	})
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
