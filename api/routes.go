package api

import (
	"net/http"

	ws "reboot01.com/js/forum/internal/websocket"
)

func (app *Application) routes() http.Handler {
	rootMux := http.NewServeMux()
	privateMux := http.NewServeMux()

	// Static file servers
	jsFileServer := http.FileServer(neuteredFileSystem{http.Dir("./static/js/")})
	rootMux.Handle("/js/", http.StripPrefix("/js", jsFileServer))

	csFileServer := http.FileServer(neuteredFileSystem{http.Dir("./static/css/")})
	rootMux.Handle("/css/", http.StripPrefix("/css", csFileServer))

	websocketManager := ws.NewWebsocketManager()

	// Public routes on rootMux
	rootMux.HandleFunc("GET /v1/status", app.status)
	rootMux.HandleFunc("GET /v1/checkauth", app.checkAuthentication)
	rootMux.HandleFunc("POST /v1/login", app.createAuthenticationToken)
	rootMux.HandleFunc("POST /v1/register", app.createUser)
	rootMux.HandleFunc("POST /v1/logout", app.logout)
	rootMux.HandleFunc("/", app.serverIndex)

	// Protected routes on privateMux (these will be under /private/ path)
	// Reminder: implement custom JSON encoder for websocket messages
	privateMux.HandleFunc("GET /ws", websocketManager.ServeWebSocket)
	privateMux.HandleFunc("GET /v1/posts", app.fetchPosts)
	privateMux.HandleFunc("POST /v1/newpost", app.newPostHandler)

	// Mount privateMux under /private/ with authentication middleware
	rootMux.Handle("/private/", http.StripPrefix("/private", app.authenticate(privateMux)))

	return app.logAccess(app.recoverPanic(rootMux))
}
