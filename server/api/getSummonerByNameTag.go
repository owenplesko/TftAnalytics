package api

import (
	"TFTAnalyticsServer/riot"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

func (env Controller) getSummonerByNameTag(w http.ResponseWriter, r *http.Request) {
	ctx := context.WithoutCancel(r.Context())

	w.Header().Set("Access-Control-Allow-Origin", "*")

	cluster := "americas"
	name := r.PathValue("name")
	tag := r.PathValue("tag")

	summoner, err := env.Service.GetOrCollectSummonerByNameTag(ctx, cluster, name, tag)
	if errors.Is(err, riot.ErrNotFound) {
		http.Error(w, "summoner not found", 404)
		return
	}
	if err != nil {
		log.Printf("Service.GetOrCollectSummonerByNameTag failed with err: %v\n", err)
		http.Error(w, "something went wrong", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	bytes, _ := json.Marshal(summoner)
	w.Write(bytes)
}
