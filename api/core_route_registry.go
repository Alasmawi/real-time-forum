package api

import (
	"net/http"
	"strings"
)

// RouteRegistry tracks routes and their allowed methods while building muxes
type RouteRegistry struct {
	routes       map[string][]string
	rootMux      *http.ServeMux
	protectedMux *http.ServeMux
}

// NewRouteRegistry creates a new route registry
func NewRouteRegistry() *RouteRegistry {
	return &RouteRegistry{
		routes:       make(map[string][]string),
		rootMux:      http.NewServeMux(),
		protectedMux: http.NewServeMux(),
	}
}

// GET registers a GET route on the appropriate mux
func (rr *RouteRegistry) GetMethod(path string, handler http.HandlerFunc) *RouteRegistry {
	rr.routes[path] = append(rr.routes[path], "GET")

	switch {
	case strings.HasPrefix(path, "/protected/"):
		// Strip /protected for the actual mux registration
		strippedPath := strings.TrimPrefix(path, "/protected")
		rr.protectedMux.HandleFunc("GET "+strippedPath, handler)
	default:
		rr.rootMux.HandleFunc("GET "+path, handler)
	}
	return rr
}

// POST registers a POST route on the appropriate mux
func (rr *RouteRegistry) PostMethod(path string, handler http.HandlerFunc) *RouteRegistry {
	rr.routes[path] = append(rr.routes[path], "POST")

	switch {
	case strings.HasPrefix(path, "/protected/"):
		// Strip /protected for the actual mux registration
		strippedPath := strings.TrimPrefix(path, "/protected")
		rr.protectedMux.HandleFunc("POST "+strippedPath, handler)
	default:
		rr.rootMux.HandleFunc("POST "+path, handler)
	}
	return rr
}

// HandleFunc registers a route without method prefix (like your current catchall)
func (rr *RouteRegistry) HandleFunc(pattern string, handler http.HandlerFunc) *RouteRegistry {
	switch {
	case strings.HasPrefix(pattern, "/protected/"):
		strippedPattern := strings.TrimPrefix(pattern, "/protected")
		rr.protectedMux.HandleFunc(strippedPattern, handler)
	default:
		rr.rootMux.HandleFunc(pattern, handler)
	}
	return rr
}

// Handle registers a route with a Handler (like your current Handle calls)
func (rr *RouteRegistry) Handle(pattern string, handler http.Handler) *RouteRegistry {
	rr.rootMux.Handle(pattern, handler)
	return rr
}

// GetMuxes returns the built muxes
func (rr *RouteRegistry) GetMuxes() (*http.ServeMux, *http.ServeMux) {
	return rr.rootMux, rr.protectedMux
}

// ValidateMethodsMiddleware returns a middleware that validates HTTP methods
func (rr *RouteRegistry) ValidateMethod() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if this path has registered methods
			if allowedMethods, exists := rr.routes[r.URL.Path]; exists {
				// Check if current method is allowed
				for _, method := range allowedMethods {
					if r.Method == method {
						next.ServeHTTP(w, r)
						return
					}
				}

				// Method not allowed
				w.Header().Set("Allow", strings.Join(allowedMethods, ", "))
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}

			// Path not in registry, continue to next handler (catchall, etc.)
			next.ServeHTTP(w, r)
		})
	}
}
