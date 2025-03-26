package api

import (
	"TFTAnalyticsServer/internal/db"
	"encoding/json"
	"log"
	"net/http"
)

func (controller Api) getSummonerAggregates(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	puuid := r.PathValue("puuid")

	aggregates, err := controller.Queries.GetSummonerAggregates(ctx, puuid)
	if err != nil {
		log.Printf("Queries.GetSummonerAggregates failed with err: %v\n", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	// return empty array not null
	if aggregates == nil {
		aggregates = []db.GetSummonerAggregatesRow{}
	}

	w.Header().Set("Content-Type", "application/json")
	bytes, _ := json.Marshal(aggregates)
	w.Write(bytes)
}
