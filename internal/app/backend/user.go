package backend

import (
	"encoding/json"
	"fmt"
	"git.cyradar.com/atd/atd/pkg/httputil"
	"git.cyradar.com/atd/atd/pkg/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/quangnc2k/do-an-gang/internal/persistance"
	"github.com/quangnc2k/do-an-gang/pkg/hxxp"
	"net/http"
	"regexp"
)

type createUserPayload struct {
	Email    string   `json:"email"`
	Password string   `json:"password"`
	Scopes   []string `json:"scopes"`
}

type updateUserPayload struct {
	Password    string   `json:"password"`
	OldPassword string   `json:"old_password,omitempty"`
	Scopes      []string `json:"scopes"`
}

func usersList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	data, err := persistance.GetRepoContainer().UserRepository.List(ctx)
	if err != nil {
		hxxp.RespondJson(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	hxxp.RespondJson(w, 200, "", data)
}

func usersCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	payload := createUserPayload{}

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		hxxp.RespondJson(w, 500, "", nil)
		return
	}

	email := payload.Email
	password := payload.Password

	if email == "" || password == "" {
		hxxp.RespondJson(w, http.StatusUnprocessableEntity, "invalid input", nil)
		return
	}

	re := regexp.MustCompile("^(([^<>()[\\]\\\\.,;:\\s@\"]+(\\.[^<>()[\\]\\\\.,;:\\s@\"]+)*)|(\".+\"))@((\\[[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\])|(([a-zA-Z\\-0-9]+\\.)+[a-zA-Z]{2,}))$")
	if !re.MatchString(email) {
		hxxp.RespondJson(w, http.StatusUnprocessableEntity, "invalid email", nil)
		return
	}

	id, err := persistance.GetRepoContainer().UserRepository.Create(
		ctx,
		uuid.New().String(),
		payload.Email,
		payload.Password,
		payload.Scopes)
	if err != nil {
		hxxp.RespondJson(w, 500, "Save failed", nil)
		return
	}

	hxxp.RespondJson(w, 200, "Success", id)
}

func usersUpdate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")
	if id == "" {
		hxxp.RespondJson(w, http.StatusBadRequest, "invalid id", nil)
		return
	}

	payload := updateUserPayload{}

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		hxxp.RespondJson(w, 500, "", nil)
		return
	}

	if payload.Password == "" || payload.OldPassword == "" {
		hxxp.RespondJson(w, http.StatusUnprocessableEntity, "invalid password", nil)
		return
	}

	currentUserID, ok := middleware.UserID(r.Context())
	if !ok {
		hxxp.RespondJson(w, nil, fmt.Errorf("%w: invalid user id", httputil.ErrRequestValidation))
		return
	}

	if currentUserID == id {
		err := persistance.GetRepoContainer().UserRepository.UpdateValidate(ctx, id, payload.OldPassword)
		if err != nil {
			httputil.Respond(ctx, w, r, nil, err)
			return
		}
	}

	_id, err := persistance.GetRepoContainer().UserRepository.Update(ctx, id, payload.Password, payload.Scopes)
	if err != nil {
		httputil.Respond(ctx, w, r, nil, httputil.ErrResourceNotUpdated)
		return
	}

	httputil.Respond(ctx, w, r, _id, nil)
}

func usersDelete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")
	if id == "" {
		httputil.Respond(ctx, w, r, nil, fmt.Errorf("%w: invalid user id", httputil.ErrRequestValidation))
		return
	}

	rows, err := persistance.GetRepoContainer().UserRepository.Delete(ctx, id)
	if err != nil || rows == 0 {
		httputil.Respond(ctx, w, r, nil, httputil.ErrResourceNotDeleted)
		return
	}

	httputil.Respond(ctx, w, r, nil, nil)
}
