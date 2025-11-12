package api

import (
	"log"
	"net/http"
)

func (controller Api) refreshUnitStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	err := controller.Service.RefreshUnitPlacement(ctx)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("refreshed"))
}
