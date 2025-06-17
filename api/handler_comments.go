package api

import (
	"net/http"

	"reboot01.com/js/realtime-forum/internal/request"
	"reboot01.com/js/realtime-forum/internal/response"
)

type PostCommentsRequest struct {
	PostID int `json:"post_id"`
}

func (app *Application) fetchPostComments(w http.ResponseWriter, r *http.Request) {
	var input PostCommentsRequest

	// Decode JSON request body
	err := request.DecodeJSON(w, r, &input)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	// Validate that post ID is provided
	if input.PostID <= 0 {
		app.errorMessage(w, r, http.StatusBadRequest, "Post ID must be a positive integer", nil)
		return
	}

	// Retrieve comments for the specific post from the database
	comments, err := app.DB.GetCommentsForPost(input.PostID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Send the comments as JSON response
	err = response.JSON(w, http.StatusOK, comments)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
}

func (app *Application) newCommentHandler(w http.ResponseWriter, r *http.Request) {
	// Define the input structure to decode JSON
	var input struct {
		PostID  int    `json:"post_id"`
		Content string `json:"content"`
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
	if input.Content == "" {
		app.errorMessage(w, r, http.StatusBadRequest, "Comment content must not be empty", nil)
		return
	}

	if input.PostID <= 0 {
		app.errorMessage(w, r, http.StatusBadRequest, "Post ID must be a positive integer", nil)
		return
	}

	// Insert the comment into the database
	commentID, err := app.DB.InsertComment(input.PostID, user.ID, input.Content)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Respond with the created comment
	response.JSON(w, http.StatusCreated, map[string]any{
		"comment_id": commentID,
		"content":    input.Content,
		"post_id":    input.PostID,
		"user_id":    user.ID,
	})
}
