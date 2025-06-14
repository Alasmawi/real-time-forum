package api

import (
	"context"
	"net/http"

	"reboot01.com/js/realtime-forum/internal/database"
)

// declaring custom contextKey type
type contextKey string

// creating constant with the type contextKey
const (
	authenticatedUserContextKey = contextKey("authenticatedUser")
)

// sets the value in request context, using the constant above
func contextSetAuthenticatedUser(r *http.Request, user *database.User) *http.Request {
	ctx := context.WithValue(r.Context(), authenticatedUserContextKey, user)
	return r.WithContext(ctx)
}

func contextGetAuthenticatedUser(r *http.Request) *database.User {
	user, ok := r.Context().Value(authenticatedUserContextKey).(*database.User)
	if !ok {
		return nil
	}

	return user
}
