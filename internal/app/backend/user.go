package backend

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/quangnc2k/do-an-gang/internal/persistance"
	"github.com/quangnc2k/do-an-gang/pkg/hxxp"
)

type createUserPayload struct {
	Email    string   `json:"email"`
	Password string   `json:"password"`
	Scopes   []string `json:"scopes"`
}

type updateUserPayload struct {
	Password    string `json:"password"`
	OldPassword string `json:"old_password,omitempty"`
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
		hxxp.RespondJson(w, 500, err.Error(), nil)
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
		nil)
	if err != nil {
		hxxp.RespondJson(w, 500, fmt.Sprintf("Save failed:", err.Error()), nil)
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
		hxxp.RespondJson(w, http.StatusUnprocessableEntity, "invalid password/old password", nil)
		return
	}

	currentUserID, ok := hxxp.UserID(r.Context())
	if !ok {
		hxxp.RespondJson(w, http.StatusUnprocessableEntity, "can't change others", nil)
		return
	}

	if currentUserID == id {
		err := persistance.GetRepoContainer().UserRepository.UpdateValidate(ctx, id, payload.OldPassword)
		if err != nil {
			hxxp.RespondJson(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}
	} else {
		hxxp.RespondJson(w, http.StatusForbidden, "invalid user", nil)
		return
	}

	_id, err := persistance.GetRepoContainer().UserRepository.Update(ctx, id, payload.Password)
	if err != nil {
		hxxp.RespondJson(w, http.StatusInternalServerError, "not updated", nil)
		return
	}

	hxxp.RespondJson(w, http.StatusOK, "success", _id)
}

func usersDelete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := chi.URLParam(r, "id")
	if id == "" {
		hxxp.RespondJson(w, http.StatusBadRequest, "invalid id", nil)
		return
	}

	rows, err := persistance.GetRepoContainer().UserRepository.Delete(ctx, id)
	if err != nil || rows == 0 {
		hxxp.RespondJson(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	hxxp.RespondJson(w, http.StatusOK, "Success", nil)
}
