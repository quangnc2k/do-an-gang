package services

import (
	"encoding/json"
	"github.com/quangnc2k/do-an-gang/pkg/hxxp"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
	"github.com/quangnc2k/do-an-gang/internal/config"
	"github.com/quangnc2k/do-an-gang/internal/persistance"
)

func AuthenticatePasswordBased(w http.ResponseWriter, r *http.Request) {
	payload := authenticatePasswordBasedPayload{}

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil || payload.Password == "" || payload.Email == "" {
		hxxp.RespondJson(w, http.StatusUnprocessableEntity, "invalid input", nil)
	}

	jwt := jwtauth.New("HS256", []byte(config.Env.JWTSecret), nil)

	user, err := persistance.GetRepoContainer().UserRepository.Login(r.Context(), payload.Email, payload.Password)
	if err != nil {
		hxxp.RespondJson(w, http.StatusForbidden, "invalid credential", nil)
		return
	}

	now := time.Now()
	jti := uuid.New().String()
	claims := map[string]interface{}{
		"aud":     "dashboard",
		"jti":     jti,
		"sub":     user.Email,
		"user_id": user.ID,
	}
	jwtauth.SetIssuedAt(claims, now)
	jwtauth.SetExpiryIn(claims, time.Hour*24*7)

	_, tokenString, err := jwt.Encode(claims)
	if err != nil {
		hxxp.RespondJson(w, http.StatusInternalServerError, "jwt not issued", nil)
		return
	}

	hxxp.RespondJson(w, http.StatusOK, "", map[string]string{
		"token": tokenString,
	})
}

type authenticatePasswordBasedPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
