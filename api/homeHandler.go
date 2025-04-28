package api

import (
	"net/http"

	response "reboot01.com/js/forum/internal/response"
)

func (app *Application) homeHandler(w http.ResponseWriter, r *http.Request) {

	posts, err := app.DB.GetAllPosts()
	if err != nil {
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}
	users, err := app.DB.GetAllUsers()
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}
	

	response.JSON(w, http.StatusOK, posts)
	response.JSON(w, http.StatusOK, users)
}
