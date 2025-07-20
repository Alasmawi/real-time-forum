package api

import (
	"net/http"
	"slices"
	"strings"
)

// RouteRegistry tracks routes and their allowed methods while building muxes
type RouteRegistry struct {
	app          *Application
	routes       map[string][]string
	publicMux    *http.ServeMux
	guestMux     *http.ServeMux
	protectedMux *http.ServeMux
}

func NewRouteRegistry() *RouteRegistry {
	return &RouteRegistry{
		routes:       make(map[string][]string),
		publicMux:    http.NewServeMux(),
		guestMux:     http.NewServeMux(),
		protectedMux: http.NewServeMux(),
	}
}

// routeToMux determines which mux to use and returns the correct mux and stripped path.
// Easily expandable for new prefixes and muxes.
func (rr *RouteRegistry) routeToMux(path string) (*http.ServeMux, string) {
	switch {
	case strings.HasPrefix(path, "/guest/"):
		strippedPath := strings.TrimPrefix(path, "/guest")
		return rr.guestMux, strippedPath
	case strings.HasPrefix(path, "/protected/"):
		strippedPath := strings.TrimPrefix(path, "/protected")
		return rr.protectedMux, strippedPath
	default:
		// All other routes go to publicMux
		return rr.publicMux, path
	}
}

// GET registers a GET route on the appropriate mux.
func (rr *RouteRegistry) GetMethod(path string, handler http.HandlerFunc) *RouteRegistry {
	rr.routes[path] = append(rr.routes[path], "GET")
	mux, finalPath := rr.routeToMux(path)
	mux.HandleFunc("GET "+finalPath, handler)
	return rr
}

// POST registers a POST route on the appropriate mux.
func (rr *RouteRegistry) PostMethod(path string, handler http.HandlerFunc) *RouteRegistry {
	rr.routes[path] = append(rr.routes[path], "POST")
	mux, finalPath := rr.routeToMux(path)
	mux.HandleFunc("POST "+finalPath, handler)
	return rr
}

// HandleFunc registers a route without method prefix
func (rr *RouteRegistry) HandleFunc(pattern string, handler http.HandlerFunc) *RouteRegistry {
	mux, finalPattern := rr.routeToMux(pattern)
	mux.HandleFunc(finalPattern, handler)
	return rr
}

// GetMuxes returns the built muxes
func (rr *RouteRegistry) GetMuxes() (*http.ServeMux, *http.ServeMux, *http.ServeMux) {
	return rr.publicMux, rr.guestMux, rr.protectedMux
}

// ValidateMethodsMiddleware returns a middleware that validates HTTP methods
func (rr *RouteRegistry) ValidateMethod() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if this path has registered methods
			if allowedMethods, exists := rr.routes[r.URL.Path]; exists {
				// Check if current method is allowed
				if slices.Contains(allowedMethods, r.Method) {
					next.ServeHTTP(w, r)
					return
				}

				// Method not allowed
				w.Header().Set("Allow", strings.Join(allowedMethods, ", "))
				rr.app.methodNotAllowed(w, r)
				return
			}

			// Path not in registry, continue to next handler (catchall, etc.)
			next.ServeHTTP(w, r)
		})
	}
}
