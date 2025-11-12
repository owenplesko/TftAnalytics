package api

import (
	"net/http"
)

func (controller Api) refreshUnitStats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	err := controller.Service.RefreshUnitPlacement(ctx)
	if err != nil {
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("refreshed"))
}
