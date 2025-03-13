package api

import (
	"encoding/json"
	"log"
	"net/http"
)

func (controller Controller) getMatchComps(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	matchId := r.PathValue("matchid")

	comps, err := controller.Service.GetMatchComps(ctx, matchId)
	if err != nil {
		log.Printf("Service.GetMatchComps failed with err: %v\n", err)
		http.Error(w, "something went wrong", 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	bytes, _ := json.Marshal(comps)
	w.Write(bytes)
}
