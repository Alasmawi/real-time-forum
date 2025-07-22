package websocket

import (
	"log"

	"reboot01.com/js/realtime-forum/internal/response"
)

// GetOnlineUsersData returns online user data for frontend
func (m *WebsocketManager) GetOnlineUsersData() []map[string]interface{} {
	m.RLock()
	defer m.RUnlock()

	var users []map[string]interface{}
	for client := range m.clients {
		users = append(users, map[string]interface{}{
			"id":       client.userID,
			"username": client.username,
			"status":   "online",
		})
	}
	return users
}

// GetOnlineUserIDs returns a slice of user IDs that are currently online
func (m *WebsocketManager) GetOnlineUserIDs() []int {
	m.RLock()
	defer m.RUnlock()

	var userIDs []int
	for client := range m.clients {
		userIDs = append(userIDs, client.userID)
	}
	return userIDs
}

// getClientByUserID finds a client by their user ID
func (m *WebsocketManager) getClientByUserID(userID int) *Client {
	m.RLock()
	defer m.RUnlock()

	for client := range m.clients {
		if client.userID == userID {
			return client
		}
	}
	return nil
}

// detectUserStatusChanges compares current online users with previous state
func (m *WebsocketManager) detectUserStatusChanges() UserStatusUpdate {
	m.Lock()
	defer m.Unlock()

	// Get current online users
	currentOnlineUsers := make(map[int]UserStatusInfo)
	for client := range m.clients {
		currentOnlineUsers[client.userID] = UserStatusInfo{
			ID:       client.userID,
			Username: client.username,
			Status:   StatusOnline,
		}
	}

	var onlineUsers []UserStatusInfo
	var offlineUsers []int

	// Find newly online users
	for id, user := range currentOnlineUsers {
		if _, wasOnline := m.previousOnlineUsers[id]; !wasOnline {
			onlineUsers = append(onlineUsers, user)
		}
	}

	// Find newly offline users
	for id := range m.previousOnlineUsers {
		if _, stillOnline := currentOnlineUsers[id]; !stillOnline {
			offlineUsers = append(offlineUsers, id)
		}
	}

	// Update previous state
	m.previousOnlineUsers = currentOnlineUsers

	return UserStatusUpdate{
		OnlineUsers:    onlineUsers,
		OfflineUserIDs: offlineUsers,
	}
}

// sendUserStatusUpdate sends status changes to a specific client
func (m *WebsocketManager) sendUserStatusUpdate(client *Client, statusUpdate UserStatusUpdate) {
	// Filter out the current user from the status update
	filteredStatusUpdate := UserStatusUpdate{
		OnlineUsers:    []UserStatusInfo{},
		OfflineUserIDs: []int{},
	}

	// Add online users excluding the current client
	for _, user := range statusUpdate.OnlineUsers {
		if user.ID != client.userID {
			filteredStatusUpdate.OnlineUsers = append(filteredStatusUpdate.OnlineUsers, user)
		}
	}

	// Add offline users excluding the current client
	for _, userID := range statusUpdate.OfflineUserIDs {
		if userID != client.userID {
			filteredStatusUpdate.OfflineUserIDs = append(filteredStatusUpdate.OfflineUserIDs, userID)
		}
	}

	// Only send if there are changes for this client
	if len(filteredStatusUpdate.OnlineUsers) == 0 && len(filteredStatusUpdate.OfflineUserIDs) == 0 {
		return
	}

	// Encode status update
	data, err := response.EncodeJSON(filteredStatusUpdate)
	if err != nil {
		log.Printf("Error encoding status update for client %d: %v", client.userID, err)
		return
	}

	// Create event
	event := Event{
		Type:    EventUserStatusUpdate,
		Payload: data,
	}

	// Send to client
	select {
	case client.egress <- event:
	default:
		// Client's channel is full, skip
		log.Printf("Client %d channel full, skipping status update", client.userID)
	}
}

// sendInitialStatusUpdate sends current online users to a newly connected client
func (m *WebsocketManager) sendInitialStatusUpdate(client *Client) {
	m.RLock()
	defer m.RUnlock()

	// Get current online users (excluding the client that just connected)
	var onlineUsers []UserStatusInfo
	for otherClient := range m.clients {
		if otherClient.userID != client.userID {
			onlineUsers = append(onlineUsers, UserStatusInfo{
				ID:       otherClient.userID,
				Username: otherClient.username,
				Status:   StatusOnline,
			})
		}
	}

	// Log initial status update info
	log.Printf("sendInitialStatusUpdate for client %d: found %d online users", client.userID, len(onlineUsers))

	// Only send if there are online users
	if len(onlineUsers) == 0 {
		log.Printf("No online users to send to client %d, skipping initial status update", client.userID)
		return
	}

	// Create initial status update
	statusUpdate := UserStatusUpdate{
		OnlineUsers:    onlineUsers,
		OfflineUserIDs: []int{}, // Empty since this is initial state
	}

	// Encode status update
	data, err := response.EncodeJSON(statusUpdate)
	if err != nil {
		log.Printf("Error encoding initial status update for client %d: %v", client.userID, err)
		return
	}

	// Create event
	event := Event{
		Type:    EventUserStatusUpdate,
		Payload: data,
	}

	// Send to client
	select {
	case client.egress <- event:
		log.Printf("Sent initial status update to client %d with %d online users", client.userID, len(onlineUsers))
	default:
		// Client's channel is full, skip
		log.Printf("Client %d channel full, skipping initial status update", client.userID)
	}
}
