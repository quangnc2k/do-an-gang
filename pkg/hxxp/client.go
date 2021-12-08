package hxxp

import (
	"context"
	"github.com/go-chi/chi/v5"
	"net/http"
	"time"
)

func NewHTTPClient() *http.Client {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	return client
}

func ChiRootContext(ctx context.Context) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			routeCtx := chi.RouteContext(r.Context())
			if routeCtx == nil {
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
			next.ServeHTTP(w, r.WithContext(context.WithValue(ctx, chi.RouteCtxKey, routeCtx)))
		})
	}
}
