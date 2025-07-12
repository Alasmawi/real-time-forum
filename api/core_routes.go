package api

import (
	"net/http"
)

func (app *Application) routes() http.Handler {
	// Create route registry
	registry := NewRouteRegistry()

	// Register static file handlers
	registry.HandleFunc("/js/", http.StripPrefix("/js/", app.neuteredFileHandler("./static/js/")).ServeHTTP).
		HandleFunc("/css/", http.StripPrefix("/css/", app.neuteredFileHandler("./static/css/")).ServeHTTP)

	// Register API routes with method validation
	registry.GetMethod("/v1/status", app.status).
		GetMethod("/v1/checkauth", app.checkAuthenticated).
		GetMethod("/v1/404", app.notFound).
		PostMethod("/v1/login", app.createAuthenticationToken).
		PostMethod("/v1/register", app.createUser)

	// Register protected routes with method validation
	// Reminder: implement custom JSON encoder for websocket messages

	registry.GetMethod("/protected/ws", app.ServeWebSocket).
		GetMethod("/protected/v1/posts", app.fetchPosts).
		PostMethod("/protected/v1/comments", app.fetchPostComments).
		GetMethod("/protected/v1/categories", app.fetchCategories).
		PostMethod("/protected/v1/newpost", app.newPostHandler).
		PostMethod("/protected/v1/logout", app.logout)

	rootMux, protectedMux := registry.GetMuxes()

	registry.HandleFunc("/", app.serveIndex)

	// Mounts protected routes with a middleware under /protected/ prefix
	protectedRouteHandler := http.StripPrefix("/protected", app.authorize(protectedMux))
	rootMux.Handle("/protected/", protectedRouteHandler)

	return app.logAccess(app.recoverPanic(app.authenticate(registry.ValidateMethod()(rootMux))))
}
