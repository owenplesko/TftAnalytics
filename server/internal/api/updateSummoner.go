package api

import (
	"context"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5"
)

func (controller Api) updateSummoner(w http.ResponseWriter, r *http.Request) {
	ctx := context.WithoutCancel(r.Context())

	puuid := r.PathValue("puuid")

	err := controller.Service.UpdateSummoner(ctx, puuid)
	if err == pgx.ErrNoRows {
		http.Error(w, "summoner not found", http.StatusNotFound)
	} else if err != nil {
		log.Printf("Service.UpdateSummonerInfo failed with err: %v\n", err)
		http.Error(w, "something went wrong", http.StatusFailedDependency)
		return
	}

	w.Write([]byte("updated"))
}
