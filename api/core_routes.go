package api

import (
	"net/http"

	ws "reboot01.com/js/realtime-forum/internal/websocket"
)

func (app *Application) routes() http.Handler {
	rootMux := http.NewServeMux()
	protectedMux := http.NewServeMux()

	rootMux.HandleFunc("/js/", http.StripPrefix("/js/", app.neuteredFileHandler("./static/js/")).ServeHTTP)
	rootMux.HandleFunc("/css/", http.StripPrefix("/css/", app.neuteredFileHandler("./static/css/")).ServeHTTP)

	rootMux.HandleFunc("GET /v1/status", app.status)
	rootMux.HandleFunc("GET /v1/checkauth", app.checkAuthenticated)
	rootMux.HandleFunc("POST /v1/login", app.createAuthenticationToken)
	rootMux.HandleFunc("POST /v1/register", app.createUser)
	// Reminder: implement custom JSON encoder for websocket messages
	websocketManager := ws.NewWebsocketManager()
	protectedMux.HandleFunc("GET /ws", websocketManager.ServeWebSocket)

	protectedMux.HandleFunc("GET /v1/posts", app.fetchPosts)
	protectedMux.HandleFunc("POST /v1/comments", app.fetchPostComments)
	protectedMux.HandleFunc("GET /v1/categories", app.fetchCategories)
	protectedMux.HandleFunc("POST /v1/newpost", app.newPostHandler)
	protectedMux.HandleFunc("POST /v1/logout", app.logout)

	// Mounts protected routes with a middleware under /protected/ prefix
	protectedRouteHandler := http.StripPrefix("/protected", app.authorize(protectedMux))
	rootMux.Handle("/protected/", protectedRouteHandler)

	rootMux.HandleFunc("/", app.serveIndex)

	return app.logAccess(app.recoverPanic(app.authenticate(rootMux)))
}
