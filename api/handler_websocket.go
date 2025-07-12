package api

import (
	"log"
	"net/http"
)

func (app *Application) ServeWebSocket(w http.ResponseWriter, r *http.Request) {
	log.Println("New WebSocket connection")

	// Get authenticated user from context (middleware already validated session)
	user := contextGetAuthenticatedUser(r)
	if user == nil {
		app.authenticationRequired(w, r)
		return
	}

	// Get session token from cookie
	sessionCookie, err := r.Cookie("session_token")
	if err != nil {
		app.authenticationRequired(w, r)
		return
	}

	// Delegate to WebSocket manager for the actual upgrade
	// Pass session token for validation in WebSocket events
	err = app.WSManager.HttpToWebsocket(w, r, user.ID, user.Username, sessionCookie.Value)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
}
