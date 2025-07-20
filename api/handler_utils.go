package api

import (
	"net/http"
	"strings"

	// "strconv"
	// "time"

	"reboot01.com/js/realtime-forum/internal/response"
	// "github.com/pascaldekloe/jwt"
)

func (app *Application) status(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"Status": "OK",
	}

	err := response.JSON(w, http.StatusOK, data)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *Application) serveIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/index.html")
}

func (app *Application) catchAll(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	// Proper 404 for public and protected API routes
	switch {
	case strings.HasPrefix(path, "/v1/"):
		app.notFound(w, r)
	case strings.HasPrefix(path, "/protected/v1/"):
		app.notFound(w, r)
	case r.Method != http.MethodGet:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	default:
		// SPA fallback for page routes
		app.serveIndex(w, r)
	}
}
