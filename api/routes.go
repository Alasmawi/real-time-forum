package api

import (
	"net/http"
)

func (app *Application) routes() http.Handler {
	mux := http.NewServeMux()
	FileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/", FileServer)

	// Custom handler to serve index.html for root path
    // mux.HandleFunc("/", app.rerout)
	mux.HandleFunc("GET /status", app.status)
	mux.HandleFunc("POST /users", app.createUser)
	mux.HandleFunc("POST /authentication-tokens", app.createAuthenticationToken)

	mux.Handle("GET /protected", app.requireAuthenticatedUser(http.HandlerFunc(app.protected)))

	return app.logAccess(app.recoverPanic(app.authenticate(mux)))
}
