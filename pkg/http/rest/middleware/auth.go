package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/tnynlabs/wyrm/pkg/http/rest"
	"github.com/tnynlabs/wyrm/pkg/utils"

	"github.com/tnynlabs/wyrm/pkg/users"
)

// Auth checks the request for credentials (either in cookie "auth_key" or header "Authorization").
// If authentication is successful the user instance will be added to the request context which
// could be accessed from handlers (e.g. r.Context().Value(UserCtxKey{})).
// If authentication is unsuccessful an appropriate error will be returned with status code 401.
func Auth(userService users.Service) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// retreive auth key from cookie or header
			authKey := keyFromCookie(r)
			if authKey == "" {
				authKey = keyFromHeader(r)
			}

			user, err := userService.GetByKey(authKey)
			if err != nil {
				serviceErr := utils.ToServiceErr(err)
				switch serviceErr.Code {
				case users.InvalidInputCode:
					rest.SendError(w, r, *serviceErr, http.StatusUnauthorized)
				default:
					rest.SendUnexpectedErr(w, r)
				}
				return
			}

			ctx := context.WithValue(r.Context(), rest.UserCtxKey{}, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// keyFromCookie tries to retreive the key string from a cookie named "auth_key".
func keyFromCookie(r *http.Request) string {
	cookie, err := r.Cookie("auth_key")
	if err != nil {
		return ""
	}
	return cookie.Value
}

// keyFromHeader tries to retreive the key string from the "Authorization" reqeust header.
// Example: "Authorization: BEARER <KEY>".
func keyFromHeader(r *http.Request) string {
	bearer := r.Header.Get("Authorization")
	if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
		return bearer[7:]
	}
	return ""
}
