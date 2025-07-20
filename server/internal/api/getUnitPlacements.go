package api

import (
	"encoding/json"
	"log"
	"net/http"
)

func (controller Api) getUnitPlacements(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	placements, err := controller.Service.GetUnitPlacments(ctx)
	if err != nil {
		log.Printf("Service.GetUnitPlacements failed with err: %v\n", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	bytes, _ := json.Marshal(placements)
	w.Write(bytes)
}
