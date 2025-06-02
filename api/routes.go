package api

import (
	"net/http"

	ws "reboot01.com/js/forum/internal/websocket"
)

func (app *Application) routes() http.Handler {
	mux := http.NewServeMux()

	// fileServer := http.FileServer(neuteredFileSystem{http.Dir("./static/")})
	// mux.Handle("/static", http.NotFoundHandler())
	// mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	jsFileServer := http.FileServer(neuteredFileSystem{http.Dir("./static/js/")})
	// mux.Handle("/js", http.NotFoundHandler())
	mux.Handle("/js/", http.StripPrefix("/js", jsFileServer))

	csFileServer := http.FileServer(neuteredFileSystem{http.Dir("./static/css/")})
	// mux.Handle("/css", http.NotFoundHandler())
	mux.Handle("/css/", http.StripPrefix("/css", csFileServer))

	websocketManager := ws.NewWebsocketManager()
	// Reminder: implement custom JSON encoder for websocket messages
	mux.HandleFunc("GET /ws", websocketManager.ServeWebSocket)

	mux.HandleFunc("GET /v1/status", app.status)
	mux.HandleFunc("POST /v1/login", app.createAuthenticationToken)
	mux.HandleFunc("POST /v1/register", app.createUser)
	mux.HandleFunc("POST /v1/notification", app.createUser)
	mux.HandleFunc("GET /v1/posts", app.fetchAllPosts)
	mux.HandleFunc("POST /v1/newpost", app.newPostHandler)

	mux.HandleFunc("/", app.serverIndex)

	mux.Handle(http.MethodGet+" /protected", app.requireAuthenticatedUser(http.HandlerFunc(app.protected)))

	return app.logAccess(app.recoverPanic(app.authenticate(mux)))
}
