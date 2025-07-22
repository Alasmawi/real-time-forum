package api

import (
	"net/http"
)

func (app *Application) routes() http.Handler {
	// Create route registry
	registry := NewRouteRegistry()

	// Public routes (no authentication required)
	registry.HandleFunc("/js/", http.StripPrefix("/js/", app.neuteredFileHandler("./static/js/")).ServeHTTP).
		HandleFunc("/css/", http.StripPrefix("/css/", app.neuteredFileHandler("./static/css/")).ServeHTTP).
		HandleFunc("/images/", http.StripPrefix("/images/", app.neuteredFileHandler("./static/images/")).ServeHTTP).
		GetMethod("/v1/status", app.status).
		GetMethod("/v1/404", app.notFound).
		PostMethod("/v1/login", app.createAuthenticationToken).
		PostMethod("/v1/register", app.createUser)

	// Guest routes (optional authentication - guests allowed, user context if authenticated)
	registry.GetMethod("/guest/v1/posts", app.fetchPosts).
		GetMethod("/guest/v1/categories", app.fetchCategories).
		GetMethod("/guest/v1/comments", app.fetchPostComments)

	// Protected routes (authentication required)
	registry.GetMethod("/protected/ws", app.ServeWebSocket).
		GetMethod("/protected/v1/user/me", app.fetchUserProfile).
		GetMethod("/protected/v1/user/{id}/message-priority", app.fetchUserMessagePriority).
		GetMethod("/protected/v1/user-list", app.fetchUserList).
		GetMethod("/protected/v1/message-history", app.getMessageHistory).
		PostMethod("/protected/v1/newcomment", app.createComment).
		PostMethod("/protected/v1/newpost", app.createPost).
		PostMethod("/protected/v1/logout", app.logout)

	publicMux, guestMux, protectedMux := registry.GetMuxes()

	// Add catchAll to publicMux for SPA routes
	publicMux.HandleFunc("/", app.catchAll)

	// Create final mux to mount different route groups with different middleware
	finalMux := http.NewServeMux()

	// Mount public routes (no authentication) - these include static files and public API
	finalMux.Handle("/", publicMux)

	// Mount guest routes (optional authentication - guests allowed)
	finalMux.Handle("/guest/", http.StripPrefix("/guest", app.authenticate(guestMux)))

	// Mount protected routes (authentication required)
	finalMux.Handle("/protected/", http.StripPrefix("/protected", app.authenticate(app.authorize(protectedMux))))

	// Apply method validation to the entire final mux
	return app.logAccess(app.recoverPanic(registry.ValidateMethod()(finalMux)))
}
