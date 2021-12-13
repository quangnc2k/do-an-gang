package hxxp

import (
	"context"
	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/jwt"
	"net/http"
)

const (
	ctxKeyToken  ctxKey = "token"
	ctxKeyUserID ctxKey = "user_id"
	ctxKeyClaims ctxKey = "claims"
)

type ctxKey string

type UserActivity struct {
	IPAddress   string `json:"ip_address"`
	UserID      string `json:"user_id"`
	Method      string `json:"method"`
	Endpoint    string `json:"endpoint"`
	Status      int    `json:"status"`
	Description string `json:"description"`
	Body        []byte `json:"body"`
}

type (
	validatorCallback func(ctx context.Context, w http.ResponseWriter, r *http.Request, id string) error
)

func Authenticator(ctx context.Context, validator validatorCallback) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, claims, err := jwtauth.FromContext(r.Context())
			if err != nil {
				RespondJson(w, http.StatusForbidden, err.Error(), nil)
				return
			}

			if token == nil || jwt.Validate(token) != nil {
				RespondJson(w, http.StatusForbidden, "invalid JWT", nil)
				return
			}

			//jti, ok := claims["jti"].(string)
			//if !ok {
			//	RespondJson(w, http.StatusForbidden, "missing JWT", nil)
			//	return
			//}

			userID, ok := claims["user_id"].(string)
			if !ok {
				RespondJson(w, http.StatusForbidden, "missing user id in jwt", nil)
				return
			}

			if err := validator(ctx, w, r, userID); err != nil {
				RespondJson(w, http.StatusForbidden, "missing JWT", nil)
				return
			}

			ctxx := r.Context()
			ctxx = context.WithValue(ctxx, ctxKeyToken, token)
			ctxx = context.WithValue(ctxx, ctxKeyUserID, userID)
			ctxx = context.WithValue(ctxx, ctxKeyClaims, claims)

			next.ServeHTTP(w, r.WithContext(ctxx))
		})
	}
}

func Claims(ctx context.Context) (map[string]interface{}, bool) {
	v, ok := ctx.Value(ctxKeyClaims).(map[string]interface{})
	return v, ok
}

func UserID(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(ctxKeyUserID).(string)
	return v, ok
}

func Token(ctx context.Context) (jwt.Token, bool) {
	v, ok := ctx.Value(ctxKeyToken).(jwt.Token)
	return v, ok
}
