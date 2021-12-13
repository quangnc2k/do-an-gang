package backend

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/quangnc2k/do-an-gang/internal/persistance"
	"github.com/quangnc2k/do-an-gang/pkg/hxxp"
	"net/http"
)

func alertsList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	all := r.URL.Query().Has("all")

	data, err := persistance.GetRepoContainer().AlertRepository.FindAll(ctx, all)
	if err != nil {
		hxxp.RespondJson(w, 500, err.Error(), nil)
		return
	}

	hxxp.RespondJson(w, 200, "", data)
}

func alertResolveMultiple(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := hxxp.UserID(r.Context())
	if !ok || userID == "" {
		hxxp.RespondJson(w, http.StatusForbidden, "missing JWT", nil)
		return
	}

	payload := resolveAlertPayload{}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		hxxp.RespondJson(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	rows, err := persistance.GetRepoContainer().AlertRepository.Resolve(ctx, payload.Resolved, userID, payload.IDs...)
	if err != nil || rows == 0 {
		hxxp.RespondJson(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	hxxp.RespondJson(w, 200, "Success", nil)
}

func alertResolve(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")
	if id == "" {
		hxxp.RespondJson(w, http.StatusBadRequest, "invalid id", nil)
		return
	}

	userID, ok := hxxp.UserID(r.Context())
	if !ok || userID == "" {
		hxxp.RespondJson(w, http.StatusForbidden, "missing JWT", nil)
		return
	}

	payload := resolveAlertPayload{}
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		hxxp.RespondJson(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	rows, err := persistance.GetRepoContainer().AlertRepository.Resolve(ctx, payload.Resolved, userID, id)
	if err != nil || rows == 0 {
		hxxp.RespondJson(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	hxxp.RespondJson(w, 200, "Success", nil)
}

type resolveAlertPayload struct {
	IDs      []string `json:"ids"`
	Resolved bool     `json:"resolved"`
}
