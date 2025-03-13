package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5"
)

func (controller Controller) getSummonerByPuuid(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	puuid := r.PathValue("puuid")

	summoner, err := controller.Service.GetSummonerByPuuid(ctx, puuid)
	if err == pgx.ErrNoRows {
		http.Error(w, "summoner not found", 404)
		return
	}
	if err != nil {
		log.Printf("Service.GetSummonerByPuuid failed with err: %v\n", err)
		http.Error(w, "something went wrong", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	bytes, _ := json.Marshal(summoner)
	w.Write(bytes)
}
