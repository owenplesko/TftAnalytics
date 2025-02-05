package api

import (
	"context"
	"encoding/json"
	"net/http"
)

func (env ApiEnv) getMatchComps(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	w.Header().Set("Access-Control-Allow-Origin", "*")

	matchId := r.PathValue("matchid")

	comps, err := env.ServiceEnv.GetMatchComps(ctx, matchId)
	if err != nil {
		http.Error(w, "something went wrong", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	bytes, _ := json.Marshal(comps)
	w.Write(bytes)
}
