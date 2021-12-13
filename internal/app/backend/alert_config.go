package backend

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/quangnc2k/do-an-gang/internal/model"
	"github.com/quangnc2k/do-an-gang/internal/persistance"
	"github.com/quangnc2k/do-an-gang/pkg/hxxp"
	"net/http"
	"time"
)

func alertConfigList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	data, err := persistance.GetRepoContainer().AlertConfigRepository.GetAll(ctx)
	if err != nil {
		hxxp.RespondJson(w, 500, err.Error(), nil)
		return
	}

	hxxp.RespondJson(w, 200, "Success", data)
}

func alertConfigReadOne(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")
	if id == "" {
		hxxp.RespondJson(w, http.StatusBadRequest, "invalid id", nil)
		return
	}

	data, err := persistance.GetRepoContainer().AlertConfigRepository.FindOneByID(ctx, id)
	if err != nil {
		hxxp.RespondJson(w, 500, err.Error(), nil)
		return
	}

	hxxp.RespondJson(w, 200, "Success", data)
}

func alertConfigCreateOne(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var payload model.AlertConfig

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		hxxp.RespondJson(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	defer r.Body.Close()

	payload.ID = uuid.New().String()

	if payload.SeverityString == "" || len(payload.Recipients) == 0 || payload.Name == "" {
		hxxp.RespondJson(w, http.StatusBadRequest, "invalid input", nil)
		return
	}

	payload.SuppressFor, err = time.ParseDuration(payload.SuppressForString)
	if err != nil {
		hxxp.RespondJson(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	err = persistance.GetRepoContainer().AlertConfigRepository.Create(ctx, &payload)
	if err != nil {
		hxxp.RespondJson(w, 500, err.Error(), nil)
		return
	}

	hxxp.RespondJson(w, 200, "Success", nil)
}

func alertConfigUpdateOne(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var payload model.AlertConfig

	id := chi.URLParam(r, "id")
	if id == "" {
		hxxp.RespondJson(w, http.StatusBadRequest, "invalid id", nil)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		hxxp.RespondJson(w, http.StatusBadRequest, err.Error(), nil)
		return
	}
	defer r.Body.Close()

	if payload.SeverityString == "" || len(payload.Recipients) == 0 || payload.Name == "" {
		hxxp.RespondJson(w, http.StatusBadRequest, "invalid input", nil)
		return
	}

	payload.SuppressFor, err = time.ParseDuration(payload.SuppressForString)
	if err != nil {
		hxxp.RespondJson(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	err = persistance.GetRepoContainer().AlertConfigRepository.UpdateOneByID(ctx, payload, id)
	if err != nil {
		hxxp.RespondJson(w, 500, err.Error(), nil)
		return
	}

	hxxp.RespondJson(w, 200, "Success", nil)
}

func alertConfigDeleteOne(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")
	if id == "" {
		hxxp.RespondJson(w, http.StatusBadRequest, "invalid id", nil)
		return
	}

	err := persistance.GetRepoContainer().AlertConfigRepository.DeleteOneByID(ctx, id)
	if err != nil {
		hxxp.RespondJson(w, 500, err.Error(), nil)
		return
	}

	hxxp.RespondJson(w, 200, "Success", nil)
}
