package middleware

import (
	"context"
	"net/http"
	"strings"
	"workout-tracker/store"
	"workout-tracker/tokens"
)

type UserMiddleware struct {
	userStore store.UserStore
}

type contextKey string

const userContextKey = contextKey("user")

func SetUser(r *http.Request, user *store.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func GetUser(r *http.Request) *store.User {
	user, ok := r.Context().Value(userContextKey).(*store.User)
	if !ok {
		panic("user not found in request")
	}
	return user
}

func (um *UserMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			r = SetUser(r, store.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			http.Error(w, "Invalid Authorization header", http.StatusUnauthorized)
			return
		}
		token := headerParts[1]
		user, err := um.userStore.GetUserToken(tokens.ScopeAuth, token)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		if user == nil {
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}
		r = SetUser(r, user)
		next.ServeHTTP(w, r)
		return
	})
}
