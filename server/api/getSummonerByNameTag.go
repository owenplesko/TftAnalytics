package api

import (
	"TFTAnalyticsServer/riot"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

func (controller Api) getSummonerByNameTag(w http.ResponseWriter, r *http.Request) {
	ctx := context.WithoutCancel(r.Context())

	cluster := "americas"
	name := r.PathValue("name")
	tag := r.PathValue("tag")

	summoner, err := controller.Service.GetOrCollectSummonerByNameTag(ctx, cluster, name, tag)
	if errors.Is(err, riot.ErrNotFound) {
		http.Error(w, "summoner not found", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("Service.GetOrCollectSummonerByNameTag failed with err: %v\n", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	bytes, _ := json.Marshal(summoner)
	w.Write(bytes)
}
