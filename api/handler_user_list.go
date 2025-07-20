package api

import (
	"net/http"
	"strconv"

	"reboot01.com/js/realtime-forum/internal/response"
	"reboot01.com/js/realtime-forum/internal/validator"
)

// UserListResponse represents the structure for user list API response
type UserListResponse struct {
	ID              int     `json:"id"`
	Username        string  `json:"username"`
	Status          int     `json:"status"` // 0 = offline, 1 = online
	LastMessageTime *string `json:"last_message_time"`
}

// fetchUserList handles GET /protected/v1/user-list
// Returns paginated user list sorted by activity
func (app *Application) fetchUserList(w http.ResponseWriter, r *http.Request) {
	user := contextGetAuthenticatedUser(r)

	var input struct {
		Page      int                 `json:"page"`
		Limit     int                 `json:"limit"`
		Validator validator.Validator `json:"-"`
	}

	// Parse pagination parameters (same pattern as posts handler)
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	// Default values
	input.Page = 1
	input.Limit = 50

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

	// Get paginated users from database
	users, err := app.DB.GetPaginatedUsersForList(user.ID, offset, input.Limit)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Convert to response format (status will be updated by frontend from WebSocket cache)
	var userList []UserListResponse
	for _, dbUser := range users {
		userList = append(userList, UserListResponse{
			ID:              dbUser.ID,
			Username:        dbUser.Username,
			Status:          0, // Default to offline, frontend will update from WebSocket cache
			LastMessageTime: dbUser.LastMessageTime,
		})
	}

	// Get total count excluding current user (consistent with pagination query)
	totalCount, err := app.DB.GetTotalUserCountExcludingUser(user.ID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Get online count from WebSocket manager (excluding current user)
	onlineCount := 0
	offlineCount := totalCount
	if app.WSManager != nil {
		onlineUserIDs := app.WSManager.GetOnlineUserIDs()
		
		// Exclude current user from online count
		for _, userID := range onlineUserIDs {
			if userID != user.ID {
				onlineCount++
			}
		}
		
		offlineCount = totalCount - onlineCount
		
		// Ensure counts are not negative
		if onlineCount < 0 {
			onlineCount = 0
		}
		if offlineCount < 0 {
			offlineCount = 0
		}
	}

	// Calculate pagination metadata
	totalPages := (totalCount + input.Limit - 1) / input.Limit
	hasNext := input.Page < totalPages
	hasPrevious := input.Page > 1

	// Prepare response with pagination metadata (same pattern as posts)
	responseData := map[string]interface{}{
		"users": userList,
		"pagination": map[string]interface{}{
			"current_page": input.Page,
			"total_pages":  totalPages,
			"total_count":  totalCount,
			"limit":        input.Limit,
			"has_next":     hasNext,
			"has_previous": hasPrevious,
		},
		"online_count":  onlineCount,
		"offline_count": offlineCount,
	}

	// Return the paginated user list as JSON
	err = response.JSON(w, http.StatusOK, responseData)
	if err != nil {
		app.serverError(w, r, err)
	}
}
