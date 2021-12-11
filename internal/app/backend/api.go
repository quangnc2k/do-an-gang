package backend

import (
	"context"
	"errors"
	"github.com/quangnc2k/do-an-gang/internal/persistance"
	"github.com/quangnc2k/do-an-gang/pkg/hxxp"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/quangnc2k/do-an-gang/internal/config"
)

func ServeBackend(ctx context.Context, addr string) (err error) {
	jwt := jwtauth.New("HS256", []byte(config.Env.JWTSecret), nil)

	r := chi.NewRouter()
	r.Use(hxxp.ChiRootContext(ctx))
	m := hxxp.Authenticator(ctx, func(ctx context.Context, w http.ResponseWriter, r *http.Request, id string) error {
		u, err := persistance.GetRepoContainer().UserRepository.FindOneByID(ctx, id)
		if err != nil || u == nil {
			return errors.New("invalid jwt")
		}
		return nil
	})

	r.Route("/api/v1", func(r chi.Router) {
		//r.Use(chiMiddleware.CleanPath)
		r.Use(chiMiddleware.URLFormat)
		r.Use(chiMiddleware.Timeout(30 * time.Second))
		r.Use(chiMiddleware.Recoverer)
		r.Use(chiMiddleware.Logger)
		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				next.ServeHTTP(w, r)
			})
		})
		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(jwt))
			r.Use(m)

			r.Route("/users", func(r chi.Router) {
				r.Get("/", usersList)
				r.Post("/", usersCreate)
				r.Patch("/{id}", usersUpdate)
				r.Delete("/{id}", usersDelete)
			})

			r.Route("/threats", func(r chi.Router) {
				r.Get("/", threatsList)
				r.Get("/overview", threatOverview)
				r.Get("/recent-affected", recentAffectedHost)
				r.Get("/recent-phase", recentAttackPhase)
				
				r.Get("/stats-severity", threatStatsSeverity)
				r.Get("/stats-phase", threatStatsPhase)
				r.Get("/stats-affected", threatTopAffectedHost)
				r.Get("/stats-suspected", threatTopAttacker)
				r.Get("/histogram-affected", threatHistogramAffected)
			})

			r.Route("/alerts", func(r chi.Router) {
				r.Get("/", alertsList)
				r.Patch("/{id}/resolve", alertResolve)
				r.Patch("/resolve", alertResolveMultiple)
			})

			r.Route("/alert-configs", func(r chi.Router) {
				r.Get("/", alertsList)
				r.Post("/", alertConfigCreateOne)
				r.Get("/{id}", alertConfigReadOne)
				r.Patch("/{id}", alertConfigUpdateOne)
				r.Delete("/{id}", alertConfigDeleteOne)
			})
		})

		r.Post("/login", authenticatePasswordBased)
	})

	log.Printf("Listening at %s\n", addr)
	srv := &http.Server{Addr: addr, Handler: r}
	log.Fatal(srv.ListenAndServe())
	return
}
