package api

import (
	"net/http"
	"strconv"

	"reboot01.com/js/realtime-forum/internal/response"
	"reboot01.com/js/realtime-forum/internal/validator"
)

// MessageHistoryResponse represents the structure for message history API response
type MessageHistoryResponse struct {
	SenderID   int    `json:"sender_id"`
	Sender     string `json:"sender"`
	ReceiverID int    `json:"receiver_id"`
	Message    string `json:"message"`
	CreatedAt  string `json:"created_at"`
}

// getMessageHistory handles GET /protected/v1/message-history
// Returns paginated message history between current user and another user
func (app *Application) getMessageHistory(w http.ResponseWriter, r *http.Request) {
	user := contextGetAuthenticatedUser(r)

	var input struct {
		UserID    int                 `json:"user_id"`
		Page      int                 `json:"page"`
		Limit     int                 `json:"limit"`
		Validator validator.Validator `json:"-"`
	}

	// Get target user ID from query parameter
	targetUserIDStr := r.URL.Query().Get("user_id")
	input.Validator.CheckField(validator.NotBlank(targetUserIDStr), "user_id", "user_id parameter is required")

	targetUserID, err := strconv.Atoi(targetUserIDStr)
	if err != nil {
		input.Validator.AddFieldError("user_id", "Invalid user_id parameter")
	} else {
		input.UserID = targetUserID

		if targetUserID != user.ID {
			_, exists, err := app.DB.GetUserById(targetUserID)
			if err != nil {
				app.serverError(w, r, err)
				return
			}
			input.Validator.CheckField(exists, "user_id", "User does not exist")
		} else {
			input.Validator.AddFieldError("user_id", "Cannot get message history with yourself")
		}
	}

	if input.Validator.HasErrors() {
		app.failedValidation(w, r, input.Validator)
		return
	}

	// Parse pagination parameters
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	// Default values
	input.Page = 1
	input.Limit = 20

	// Parse and validate page parameter
	if pageStr != "" {
		parsedPage, err := strconv.Atoi(pageStr)
		if err != nil {
			input.Validator.AddFieldError("page", "Invalid page parameter")
		} else {
			input.Validator.CheckField(validator.MinInt(parsedPage, 1), "page", "Page must be at least 1")
			input.Validator.CheckField(validator.MaxInt(parsedPage, 1000), "page", "Page cannot exceed 1000")
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

	// Get paginated messages from database
	messages, err := app.DB.GetPaginatedMessageHistory(user.ID, input.UserID, offset, input.Limit)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Get total message count
	totalCount, err := app.DB.GetMessageCount(user.ID, input.UserID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Convert to response format
	var messageHistory []MessageHistoryResponse
	for _, dbMessage := range messages {
		messageHistory = append(messageHistory, MessageHistoryResponse{
			SenderID:   dbMessage.SenderID,
			Sender:     dbMessage.Sender,
			ReceiverID: dbMessage.ReceiverID,
			Message:    dbMessage.Message,
			CreatedAt:  dbMessage.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	// Calculate pagination metadata
	totalPages := (totalCount + input.Limit - 1) / input.Limit
	hasNext := input.Page < totalPages
	hasPrevious := input.Page > 1

	// Prepare response with pagination metadata
	responseData := map[string]interface{}{
		"messages": messageHistory,
		"pagination": map[string]interface{}{
			"current_page": input.Page,
			"total_pages":  totalPages,
			"total_count":  totalCount,
			"limit":        input.Limit,
			"has_next":     hasNext,
			"has_previous": hasPrevious,
		},
		"target_user_id": input.UserID,
	}

	// Return the paginated message history as JSON
	err = response.JSON(w, http.StatusOK, responseData)
	if err != nil {
		app.serverError(w, r, err)
	}
}
