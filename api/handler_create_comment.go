package api

import (
	"net/http"

	"reboot01.com/js/realtime-forum/internal/request"
	"reboot01.com/js/realtime-forum/internal/response"
	t "reboot01.com/js/realtime-forum/internal/time"
	"reboot01.com/js/realtime-forum/internal/validator"
)

func (app *Application) createComment(w http.ResponseWriter, r *http.Request) {
	var input struct {
		PostID    int                 `json:"post_id"`
		Content   string              `json:"content"`
		Validator validator.Validator `json:"-"`
	}

	err := request.DecodeJSON(w, r, &input)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	user := contextGetAuthenticatedUser(r)
	if user == nil {
		app.authenticationRequired(w, r)
		return
	}

	input.Validator.CheckField(validator.NotBlank(input.Content), "reply-content", "Content must not be empty")
	input.Validator.CheckField(validator.MaxRunes(input.Content, 500), "reply-content", "Content must not exceed 500 characters")

	if input.Validator.HasErrors() {
		app.failedValidation(w, r, input.Validator)
		return
	}

	currentDateTime := t.CurrentTime()

	// Insert the comment into the database
	commentID, err := app.DB.InsertComment(input.PostID, user.ID, input.Content, currentDateTime)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	responseData := map[string]any{
		"comment_id": commentID,
		"username":   user.Username,
		"content":    input.Content,
		"created_at": currentDateTime,
	}

	err = response.JSON(w, http.StatusCreated, responseData)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
}
