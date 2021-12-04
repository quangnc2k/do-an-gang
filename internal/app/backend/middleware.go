package dashboard

import (
	"context"
	"net/http"

	"git.cyradar.com/atd/atd/pkg/middleware"
)

const (
	ctxKeyStatsOrigin   ctxKey = "stats-origin"
	ctxKeyStatsSensorID ctxKey = "stats-sensor-id"
)

func middlewareStats(ctx context.Context) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return middleware.TimeRange(ctx)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctxx := r.Context()

			origin := r.URL.Query().Get("origin")
			if origin == "" {
				origin = "*"
			}

			sensorID := r.URL.Query().Get("sensorID")

			ctxx = context.WithValue(ctxx, ctxKeyStatsOrigin, origin)
			ctxx = context.WithValue(ctxx, ctxKeyStatsSensorID, sensorID)
			next.ServeHTTP(w, r.WithContext(ctxx))
		}))
	}
}
