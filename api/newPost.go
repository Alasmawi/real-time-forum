package api

// import (
// 	"net/http"

// 	"reboot01.com/js/forum/internal/database"
// 	"reboot01.com/js/forum/internal/request"
// 	"reboot01.com/js/forum/internal/response"
// 	"reboot01.com/js/forum/internal/validator"
// )

// func (app *Application) newPostHandler(w http.ResponseWriter, r *http.Request) {
// 	// Define the input structure to decode JSON
// 	var input struct {
// 		ID         int    `json:"id"`
// 		Content    string `json:"content"`
// 		Categories []int  `json:"categories"`
// 	}

// 	// Decode the JSON request body
// 	err := request.DecodeJSON(w, r, &input)
// 	if err != nil {
// 		app.badRequest(w, r, err)
// 		return
// 	}

// 	// Validate the input
// 	var v validator.Validator
// 	v.CheckField(input.Content != "", "content", "Content must not be empty")
// 	// v.CheckField(len(input.Categories) > 0, "categories", "At least one category must be selected")

// 	if !v.Valid() {
// 		app.failedValidation(w, r, v.Errors)
// 		return
// 	}

// 	// Create a Post object
// 	post := database.Post{
// 		ID:      input.ID,
// 		Content: input.Content,
// 	}

// 	// Insert the post into the database
// 	postID, err := app.DB.InsertPost(post.Content, post.ID)
// 	if err != nil {
// 		app.serverError(w, r, err)
// 		return
// 	}

// 	// Insert categories into the database
// 	for _, categoryID := range input.Categories {
// 		err := app.DB.InsertPostCategory(postID, categoryID)
// 		if err != nil {
// 			app.serverError(w, r, err)
// 			return
// 		}
// 	}

// 	// Respond with the created post
// 	response.JSON(w, http.StatusCreated, map[string]any{
// 		"post_id":    postID,
// 		"content":    post.Content,
// 		"categories": input.Categories,
// 	})
// }
