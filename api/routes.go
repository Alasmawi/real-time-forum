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

	mux.HandleFunc("GET /status", app.status)
	mux.HandleFunc("/register", app.createUser)
	mux.HandleFunc("GET /login", app.createAuthenticationToken)
	// mux.HandleFunc("GET /error", )
	mux.HandleFunc("GET /ws", websocketManager.ServeWebSocket)

	mux.Handle("GET /protected", app.requireAuthenticatedUser(http.HandlerFunc(app.protected)))

	return app.logAccess(app.recoverPanic(app.authenticate(mux)))
}
