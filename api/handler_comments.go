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
