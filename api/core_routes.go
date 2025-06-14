package api

import (
	"net/http"

	ws "reboot01.com/js/realtime-forum/internal/websocket"
)

func (app *Application) routes() http.Handler {
	rootMux := http.NewServeMux()
	protectedMux := http.NewServeMux()

	jsFileServer := http.FileServer(neuteredFileSystem{http.Dir("./static/js/")})
	rootMux.Handle("/js/", http.StripPrefix("/js", jsFileServer))

	csFileServer := http.FileServer(neuteredFileSystem{http.Dir("./static/css/")})
	rootMux.Handle("/css/", http.StripPrefix("/css", csFileServer))

	rootMux.HandleFunc("GET /v1/status", app.status)
	rootMux.HandleFunc("GET /v1/checkauth", app.checkAuthentication)
	rootMux.HandleFunc("POST /v1/login", app.createAuthenticationToken)
	rootMux.HandleFunc("POST /v1/register", app.createUser)
	rootMux.HandleFunc("POST /v1/logout", app.logout)
	rootMux.HandleFunc("/", app.serveIndex)

	// Reminder: implement custom JSON encoder for websocket messages
	websocketManager := ws.NewWebsocketManager()
	protectedMux.HandleFunc("GET /ws", websocketManager.ServeWebSocket)

	protectedMux.HandleFunc("GET /v1/posts", app.fetchPosts)
	protectedMux.HandleFunc("POST /v1/newpost", app.newPostHandler)

	// Mounts protected routes with authentication middleware under /private/ prefix
	protectedRouteHandler := http.StripPrefix("/protected", app.authenticate(protectedMux))
	rootMux.Handle("/protected/", protectedRouteHandler)

	return app.logAccess(app.recoverPanic(rootMux))
}
