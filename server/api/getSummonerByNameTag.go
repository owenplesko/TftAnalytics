package api

import (
	"TFTAnalyticsServer/riot"
	"context"
	"encoding/json"
	"net/http"
)

func (env Controller) getSummonerByNameTag(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	w.Header().Set("Access-Control-Allow-Origin", "*")

	cluster := "americas"
	name := r.PathValue("name")
	tag := r.PathValue("tag")

	summoner, err := env.Service.GetOrCollectSummonerByNameTag(ctx, cluster, name, tag)
	if err == riot.ErrNotFound {
		http.Error(w, "summoner not found", 404)
		return
	}
	if err != nil {
		http.Error(w, "something went wrong", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	bytes, _ := json.Marshal(summoner)
	w.Write(bytes)
}
