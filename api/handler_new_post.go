package api

import (
	"net/http"

	"reboot01.com/js/realtime-forum/internal/request"
	"reboot01.com/js/realtime-forum/internal/response"
	"reboot01.com/js/realtime-forum/internal/validator"
)

func (app *Application) newPostHandler(w http.ResponseWriter, r *http.Request) {
	// Define the input structure to decode JSON
	var input struct {
		Content    string `json:"content"`
		Categories []int  `json:"categories"`
	}

	// Decode the JSON request body
	err := request.DecodeJSON(w, r, &input)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	// Get the authenticated user from context
	user := contextGetAuthenticatedUser(r)
	if user == nil {
		app.authenticationRequired(w, r)
		return
	}

	// Validate the input
	var v validator.Validator
	v.CheckField(input.Content != "", "content", "Content must not be empty")
	// v.CheckField(len(input.Categories) > 0, "categories", "At least one category must be selected")

	if v.HasErrors() {
		app.failedValidation(w, r, v)
		return
	}

	// Insert the post into the database
	postID, err := app.DB.InsertPost(input.Content, user.ID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Insert categories into the database
	for _, categoryID := range input.Categories {
		err := app.DB.InsertPostCategory(postID, categoryID)
		if err != nil {
			app.serverError(w, r, err)
			return
		}
	}

	// Respond with the created post
	response.JSON(w, http.StatusCreated, map[string]any{
		"post_id":    postID,
		"content":    input.Content,
		"user_id":    user.ID,
		"categories": input.Categories,
	})
}
