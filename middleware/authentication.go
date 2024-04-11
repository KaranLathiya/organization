package middleware

import (
	"context"
	"net/http"

	error_handling "organization/error"
)

var UserCtxKey = &contextKey{"user"}

type contextKey struct {
	authorizationKey string
}

func Authentication(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationKey := r.Header.Get("Auth-user")
		if authorizationKey == "" {
			error_handling.ErrorMessageResponse(w, error_handling.HeaderDataMissing)
			return
		}
		ctx := context.WithValue(r.Context(), UserCtxKey, authorizationKey)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
	
}
