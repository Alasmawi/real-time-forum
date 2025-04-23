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

	// Use the JSON function to send the response
	// w.Header().Set("Content-Type", "application/json") // Set the Content-Type header
	// postjson, err := json.Marshal(posts)              // Marshal posts to JSON
	// if err != nil {
	// 	http.Error(w, "Failed to encode posts", http.StatusInternalServerError)
	// 	return
	// }
	// w.Write(postjson) // Write the JSON data to the response body
}
