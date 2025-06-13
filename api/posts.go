package api

import (
	"net/http"

	"reboot01.com/js/realtime-forum/internal/response"
)

func (app *Application) fetchPosts(w http.ResponseWriter, r *http.Request) {
	// Retrieve posts from the database
	posts, err := app.DB.GetAllPosts()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Send the posts as JSON response
	err = response.JSON(w, http.StatusOK, posts)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
}
