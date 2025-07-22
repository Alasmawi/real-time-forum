package api

import (
	"fmt"
	"net/http"
	"strconv"

	"reboot01.com/js/realtime-forum/internal/response"
)

// UserProfileResponse represents the user profile data for the frontend
type UserProfileResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

// UserMessagePriorityResponse represents user message priority data
type UserMessagePriorityResponse struct {
	ID              int     `json:"id"`
	Username        string  `json:"username"`
	LastMessageTime *string `json:"last_message_time"`
}

// fetchUserProfile handles GET /protected/v1/user/me
// Returns current user's profile information
func (app *Application) fetchUserProfile(w http.ResponseWriter, r *http.Request) {
	// Get current user from context (set by authentication middleware)
	user := contextGetAuthenticatedUser(r)

	// Create response with user profile data
	userProfile := UserProfileResponse{
		ID:       user.ID,
		Username: user.Username,
	}

	// Return user profile as JSON
	err := response.JSON(w, http.StatusOK, userProfile)
	if err != nil {
		app.serverError(w, r, err)
	}
}

// fetchUserMessagePriority handles GET /protected/v1/user/{id}/message-priority
// Returns user's message priority data for list ordering
func (app *Application) fetchUserMessagePriority(w http.ResponseWriter, r *http.Request) {
	// Get current user from context
	currentUser := contextGetAuthenticatedUser(r)

	// Get target user ID from URL path parameter
	userIDStr := r.PathValue("id")
	if userIDStr == "" {
		app.badRequest(w, r, fmt.Errorf("user ID is required"))
		return
	}

	targetUserID, err := strconv.Atoi(userIDStr)
	if err != nil {
		app.badRequest(w, r, fmt.Errorf("invalid user ID: %s", userIDStr))
		return
	}

	// Check if target user exists
	_, exists, err := app.DB.GetUserById(targetUserID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	if !exists {
		app.notFound(w, r)
		return
	}

	// Get user with message priority data
	userWithMessage, err := app.DB.GetUserMessagePriority(currentUser.ID, targetUserID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Create response
	userResponse := UserMessagePriorityResponse{
		ID:              userWithMessage.ID,
		Username:        userWithMessage.Username,
		LastMessageTime: userWithMessage.LastMessageTime,
	}

	// Return as JSON
	err = response.JSON(w, http.StatusOK, userResponse)
	if err != nil {
		app.serverError(w, r, err)
	}
}
