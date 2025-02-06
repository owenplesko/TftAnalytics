package api

import (
	"context"
	"net/http"

	"github.com/jackc/pgx/v5"
)

func (env Controller) updateSummoner(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	w.Header().Set("Access-Control-Allow-Origin", "*")

	puuid := r.PathValue("puuid")

	err := env.Service.UpdateAllSummonerInfo(ctx, puuid)
	if err == pgx.ErrNoRows {
		http.Error(w, "summoner not found", 404)
	} else if err != nil {
		http.Error(w, "something went wrong", 500)
	}
}
