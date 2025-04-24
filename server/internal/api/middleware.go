package api

import (
	"net/http"

	"github.com/owenplesko/TftAnalytics/pkg/riot/scheduler"
)

// TODO: if admin priority = 2
func setSchedulerPriority(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = scheduler.WithPriority(ctx, 1)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func setCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}
