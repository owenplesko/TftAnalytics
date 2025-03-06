package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5"
)

func (env Controller) getSummonerByPuuid(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	w.Header().Set("Access-Control-Allow-Origin", "*")

	puuid := r.PathValue("puuid")

	summoner, err := env.Service.GetSummonerByPuuid(ctx, puuid)
	if err == pgx.ErrNoRows {
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
