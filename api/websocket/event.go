package websocket

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

// Event is the Messages sent over the websocket
// Used to differ between different actions
type Event struct {
	// Type is the message type sent
	Type string `json:"type"`
	// Payload is the data Based on the Type
	Payload json.RawMessage `json:"payload"`
}

// EventHandler is a function signature that is used to affect messages on the socket and triggered
// depending on the type
type EventHandler func(event Event, c *Client) error

const (
	// EventSendMessage is the event name for new chat messages sent
	EventSendMessage = "send_message"
	// EventNewMessage is a response to send_message
	EventNewMessage = "new_message"
	// EventChangeRoom is event when switching rooms
	EventChangeRoom = "change_room"
	// EventRequestUserList is the event to request user list
	EventRequestUserList = "request_user_list"
	// EventUserListUpdate is the response with user list
	EventUserListUpdate = "user_list_update"
	// EventUserStatusChange is when users come online/offline
	EventUserStatusChange = "user_status_change"
)

// SendMessageEvent is the payload sent in the
// send_message event
type SendMessageEvent struct {
	Message      string `json:"message"`
	SenderID     int    `json:"sender_id"`     // Must match client's userID
	ReceiverID   int    `json:"receiver_id"`   // Target user
	SessionToken string `json:"session_token"` // Session token for validation
}

// NewMessageEvent is returned when responding to send_message
type NewMessageEvent struct {
	SendMessageEvent           // Embedded struct
	From             string    `json:"from"` // Username (server adds this)
	Sent             time.Time `json:"sent"` // Server adds timestamp
}

// SendMessageHandler will send out a message to all other participants in the chat
func SendMessageHandler(event Event, c *Client) error {
	// Marshal Payload into wanted format
	var chatevent SendMessageEvent
	if err := json.Unmarshal(event.Payload, &chatevent); err != nil {
		return fmt.Errorf("bad payload in request: %v", err)
	}

	// CRITICAL: Validate session token against database
	user, found, err := c.manager.DB.GetUserBySession(chatevent.SessionToken)
	if err != nil || !found {
		c.closeWithReason(websocket.ClosePolicyViolation, "Invalid session token")
		return nil
	}

	// Prevent users from sending messages to themselves
	if user.ID == chatevent.ReceiverID {
		c.sendErrorEvent("SELF_MESSAGE", "Cannot send messages to yourself")
		return nil
	}

	// Check if receiver is online
	receiverClient := c.manager.getClientByUserID(chatevent.ReceiverID)
	if receiverClient == nil {
		c.sendErrorEvent("RECEIVER_OFFLINE", "Recipient is not online")
		return nil
	}

	// Generate single timestamp for consistency between database and client
	messageTime := time.Now()

	// Save message to database
	_, err = c.manager.DB.InsertMessage(user.ID, chatevent.ReceiverID, chatevent.Message, messageTime)
	if err != nil {
		return fmt.Errorf("failed to save message: %v", err)
	}

	// Prepare outgoing message
	var broadMessage NewMessageEvent
	broadMessage.SendMessageEvent = chatevent // Copy all fields from embedded struct
	broadMessage.From = user.Username         // Server sets the username from database
	broadMessage.Sent = messageTime           // Use same timestamp as database

	data, err := json.Marshal(broadMessage)
	if err != nil {
		return fmt.Errorf("failed to marshal broadcast message: %v", err)
	}

	// Create event
	var outgoingEvent Event
	outgoingEvent.Payload = data
	outgoingEvent.Type = EventNewMessage

	// Send to receiver
	receiverClient.egress <- outgoingEvent

	// Send confirmation back to sender
	c.egress <- outgoingEvent

	return nil
}

// UserListItem represents a user in the user list
type UserListItem struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

// UserListResponse is the response to request_user_list
type UserListResponse struct {
	OnlineUsers  []UserListItem `json:"online-users"`
	OfflineUsers []UserListItem `json:"offline-users"`
}

// UserStatusChangeEvent is sent when users come online/offline
type UserStatusChangeEvent struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Status   string `json:"status"` // "online" or "offline"
}

// RequestUserListEvent is the payload for requesting user list
type RequestUserListEvent struct {
	SessionToken string `json:"session_token"` // Session token for validation
}

func RequestUserListHandler(event Event, c *Client) error {

	var requestEvent RequestUserListEvent
	if err := json.Unmarshal(event.Payload, &requestEvent); err != nil {
		return fmt.Errorf("bad payload in request: %v", err)
	}

	user, found, err := c.manager.DB.GetUserBySession(requestEvent.SessionToken)
	if err != nil || !found {
		c.closeWithReason(websocket.ClosePolicyViolation, "Invalid session token")
		return nil
	}

	// Get users sorted by last message activity
	users, err := c.manager.DB.GetUsersOrderedByActivity(user.ID)
	if err != nil {
		return fmt.Errorf("failed to get users: %v", err)
	}

	// Get online user IDs
	onlineUserIDs := c.manager.GetOnlineUserIDs()
	onlineMap := make(map[int]bool)
	for _, id := range onlineUserIDs {
		onlineMap[id] = true
	}

	// Split users into online and offline
	var response UserListResponse
	for _, dbUser := range users {
		userItem := UserListItem{
			ID:       dbUser.ID,
			Username: dbUser.Username,
		}

		if onlineMap[dbUser.ID] {
			response.OnlineUsers = append(response.OnlineUsers, userItem)
		} else {
			response.OfflineUsers = append(response.OfflineUsers, userItem)
		}
	}

	// Send response
	data, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("failed to marshal user list response: %v", err)
	}

	outgoingEvent := Event{
		Type:    EventUserListUpdate,
		Payload: data,
	}

	c.egress <- outgoingEvent
	return nil
}
