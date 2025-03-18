package api

import (
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from Snippetbox"))
}

func (app *Application) routes() http.Handler {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	// Custom handler to serve index.html for root path
	// mux.HandleFunc("/", app.rerout)
	// mux.HandleFunc("GET /", app.serveIndex)
	mux.HandleFunc("GET /status", app.status)
	mux.HandleFunc("POST /users", app.createUser)
	mux.HandleFunc("POST /authentication-tokens", app.createAuthenticationToken)

	mux.Handle("GET /protected", app.requireAuthenticatedUser(http.HandlerFunc(app.protected)))

	return app.logAccess(app.recoverPanic(app.authenticate(mux)))
}
