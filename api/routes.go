package api

import (
	"net/http"
)

func (app *Application) routes() http.Handler {
	mux := http.NewServeMux()
	// fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/", http.FileServer(http.Dir("./static")))

	mux.HandleFunc("GET /status", app.status)
	mux.HandleFunc("POST /users", app.createUser)
	mux.HandleFunc("POST /authentication-tokens", app.createAuthenticationToken)

	mux.Handle("GET /protected", app.requireAuthenticatedUser(http.HandlerFunc(app.protected)))

	return app.logAccess(app.recoverPanic(app.authenticate(mux)))
}
