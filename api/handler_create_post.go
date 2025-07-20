package api

import (
	"net/http"

	"reboot01.com/js/realtime-forum/internal/request"
	"reboot01.com/js/realtime-forum/internal/response"
	t "reboot01.com/js/realtime-forum/internal/time"
	"reboot01.com/js/realtime-forum/internal/validator"
)

func (app *Application) createPost(w http.ResponseWriter, r *http.Request) {
	// Define the input structure to decode JSON
	var input struct {
		Content    string              `json:"content"`
		Categories []int               `json:"categories"`
		Validator  validator.Validator `json:"-"`
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

	input.Validator.CheckField(validator.NotBlank(input.Content), "post-content", "Content must not be empty")
	input.Validator.CheckField(validator.MaxRunes(input.Content, 500), "post-content", "Content must not exceed 500 characters")

	// Validate categories (optional but if provided, must be valid)
	input.Validator.CheckField(len(input.Categories) <= 10, "categories", "Cannot select more than 10 categories")

	if input.Validator.HasErrors() {
		app.failedValidation(w, r, input.Validator)
		return
	}

	currentDateTime := t.CurrentTime()

	// Insert the post into the database
	postID, err := app.DB.InsertPost(input.Content, currentDateTime, user.ID)
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
	responseData := map[string]any{
		"post_id":    postID,
		"content":    input.Content,
		"username":   user.Username,
		"categories": input.Categories,
		"created_at": currentDateTime,
	}

	err = response.JSON(w, http.StatusCreated, responseData)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
}
