package api

import (
	"net/http"
	"strconv"

	"reboot01.com/js/realtime-forum/internal/response"
	"reboot01.com/js/realtime-forum/internal/validator"
)

func (app *Application) fetchPostComments(w http.ResponseWriter, r *http.Request) {
	var input struct {
		PostID    int                 `json:"post_id"`
		Page      int                 `json:"page"`
		Limit     int                 `json:"limit"`
		Validator validator.Validator `json:"-"`
	}

	// Parse query parameters
	postIDStr := r.URL.Query().Get("post_id")
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	// Validate that post ID is provided
	input.Validator.CheckField(validator.NotBlank(postIDStr), "post_id", "Post ID is required")

	if postIDStr != "" {
		postID, err := strconv.Atoi(postIDStr)
		if err != nil {
			input.Validator.AddFieldError("post_id", "Invalid post_id parameter")
		} else {
			input.Validator.CheckField(validator.MinInt(postID, 1), "post_id", "Post ID must be a positive integer")
			input.PostID = postID
		}
	}

	// Default pagination values
	input.Page = 1
	input.Limit = 20

	// Parse and validate page parameter
	if pageStr != "" {
		parsedPage, err := strconv.Atoi(pageStr)
		if err != nil {
			input.Validator.AddFieldError("page", "Invalid page parameter")
		} else {
			input.Validator.CheckField(validator.MinInt(parsedPage, 1), "page", "Page must be at least 1")
			input.Validator.CheckField(validator.MaxInt(parsedPage, 10000), "page", "Page cannot exceed 10000")
			input.Page = parsedPage
		}
	}

	// Parse and validate limit parameter
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil {
			input.Validator.AddFieldError("limit", "Invalid limit parameter")
		} else {
			input.Validator.CheckField(validator.MinInt(parsedLimit, 1), "limit", "Limit must be at least 1")
			input.Validator.CheckField(validator.MaxInt(parsedLimit, 100), "limit", "Limit cannot exceed 100")
			input.Limit = parsedLimit
		}
	}

	if input.Validator.HasErrors() {
		app.failedValidation(w, r, input.Validator)
		return
	}

	// Calculate offset
	offset := (input.Page - 1) * input.Limit

	// Retrieves paginated comments for the specific post from the database
	comments, err := app.DB.GetPaginatedCommentsForPost(input.PostID, offset, input.Limit)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Gets total count for pagination metadata
	totalCount, err := app.DB.GetCommentsCountForPost(input.PostID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Calculate pagination metadata
	totalPages := (totalCount + input.Limit - 1) / input.Limit
	hasNext := input.Page < totalPages
	hasPrevious := input.Page > 1

	responseData := map[string]any{
		"comments": comments,
		"pagination": map[string]any{
			"current_page": input.Page,
			"total_pages":  totalPages,
			"total_count":  totalCount,
			"limit":        input.Limit,
			"has_next":     hasNext,
			"has_previous": hasPrevious,
		},
	}

	err = response.JSON(w, http.StatusOK, responseData)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
}
