package api

import (
	"net/http"

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
