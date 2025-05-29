package api

import (
	"net/http"

	ws "reboot01.com/js/forum/internal/websocket"
)

func (app *Application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./static/")})
	mux.Handle("/static", http.NotFoundHandler())
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	websocketManager := ws.NewWebsocketManager()
	mux.HandleFunc("GET /ws", websocketManager.ServeWebSocket)

	mux.HandleFunc("GET /v1/status", app.status)
	mux.HandleFunc("POST /v1/login", app.createAuthenticationToken)
	mux.HandleFunc("POST /v1/register", app.createUser)
	mux.HandleFunc("POST /v1/notification", app.createUser)

	mux.HandleFunc("/", app.serverIndex)

	// mux.HandleFunc("GET /{$}", app.serverIndex)
	// mux.HandleFunc("GET /register", app.serverIndex)
	// mux.HandleFunc("GET /chat", app.serverIndex)

	mux.Handle(http.MethodGet+" /protected", app.requireAuthenticatedUser(http.HandlerFunc(app.protected)))

	return app.logAccess(app.recoverPanic(app.authenticate(mux)))
}
