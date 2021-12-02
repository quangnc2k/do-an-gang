package services

import (
	"context"
	"errors"
	"git.cyradar.com/atd/atd/pkg/httputil"
	"git.cyradar.com/atd/atd/pkg/httputil/render"
	"github.com/quangnc2k/do-an-gang/internal/config"
	"github.com/quangnc2k/service-demo/pkg/hxxp"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

// the ctx is pkg-scoped context - not request-scoped context
func authenticatePasswordBased(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := shared.CtxLog(ctx)
	payload := &authenticatePasswordBasedPayload{}
	if err := render.Bind(ctx, r, payload); err != nil {
		httputil.Respond(ctx, w, r, nil, httputil.ErrRequestValidation)
		return
	}

	datastoreUser, ok := ctxDatastoreUser(ctx)
	if !ok {
		httputil.Respond(ctx, w, r, nil, errServiceMissingDatastoreUser)
		return
	}

	tokenStore, ok := ctxDatastoreToken(ctx)
	if !ok {
		hxxp.ResponseJson()
		return
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		httputil.Respond(ctx, w, r, nil, httputil.ErrRequestMissingIP)
		return
	}

	jwt := config.Env.JWTSecret
	if jwt == "" {
		httputil.Respond(ctx, w, r, nil, errServiceMissingJWT)
		return
	}

	user, err := datastoreUser.Login(r.Context(), payload.Email, payload.Password)
	if err != nil {
		httputil.Respond(ctx, w, r, nil, httputil.ErrAuthInvalidCredentials)
		return
	}

	now := time.Now()
	jti := uuid.New().String()
	claims := map[string]interface{}{
		"aud":     "dashboard",
		"jti":     jti,
		"scopes":  user.Scopes,
		"sub":     user.Email,
		"user_id": user.ID,
	}
	jwtauth.SetIssuedAt(claims, now)
	jwtauth.SetExpiryIn(claims, time.Hour*24*7)
	for _, scope := range user.Scopes {
		if strings.ToLower(scope) == "monitor:*" {
			jwtauth.SetExpiryIn(claims, time.Hour*24*365)
			break
		}
	}

	token, tokenString, err := jwt.Encode(claims)
	if err != nil {
		httputil.Respond(ctx, w, r, nil, httputil.ErrAuthJWTNotIssued)
		return
	}

	log.Debugw("issuing JWT", "token", token)
	t := &datastore.Token{
		ID:        jti,
		UserID:    user.ID,
		ExpiredAt: time.Unix(claims["exp"].(int64), 0),
		Revoked:   false,
		UserAgent: r.UserAgent(),
		IPAddress: ip,
		CreatedAt: now,
	}

	err = tokenStore.Save(r.Context(), t)
	if err != nil {
		httputil.Respond(ctx, w, r, nil, httputil.ErrAuthJWTNotSaved)
		return
	}

	httputil.Respond(ctx, w, r, map[string]string{
		"token": tokenString,
	}, nil)
}

type authenticatePasswordBasedPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (payload *authenticatePasswordBasedPayload) Bind(ctx context.Context, r *http.Request) error {
	email := payload.Email
	password := payload.Password

	if email == "" || password == "" {
		return errors.New("missing login credentials")
	}

	return nil
}
