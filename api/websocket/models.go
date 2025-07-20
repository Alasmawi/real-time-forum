package websocket

import (
	"encoding/json"
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

// WebSocket Event Type Constants
const (
	// EventSendMessage is the event name for new chat messages sent
	EventSendMessage = "send_message"
	// EventReceiveMessage is a response to send_message
	EventReceiveMessage = "receive_message"
	// EventChangeRoom is event when switching rooms
	EventChangeRoom = "change_room"
	// EventUserStatusUpdate is the new event for online/offline status changes
	EventUserStatusUpdate = "user_status_update"
	// EventSendTyping is when user starts/stops typing
	EventSendTyping = "send_typing"
	// EventNewTyping is the response to send_typing
	EventNewTyping = "new_typing"
	// EventError is for error messages
	EventError = "error"
)

// User Status Constants
const (
	StatusOffline = 0
	StatusOnline  = 1
	StatusBusy    = 2
	StatusDND     = 3
)

// SendMessageEvent is the payload sent in the
// send_message event
type SendMessageEvent struct {
	Message      string `json:"message"`
	ReceiverID   int    `json:"receiver_id"`   // Target user
	SessionToken string `json:"session_token"` // Session token for validation
}

// ReceiveMessageEvent is returned when responding to send_message
type ReceiveMessageEvent struct {
	Message    string `json:"message"`
	SenderID   int    `json:"sender_id"`   // Who sent the message
	ReceiverID int    `json:"receiver_id"` // Who received the message
	SentAt     string `json:"sent_at"`     // Server adds timestamp
}

// SendTypingEvent is the payload sent in the send_typing event
type SendTypingEvent struct {
	ReceiverID   int    `json:"receiver_id"`   // Target user
	IsTyping     bool   `json:"is_typing"`     // true = start typing, false = stop typing
	SessionToken string `json:"session_token"` // Session token for validation
}

// NewTypingEvent is returned when responding to send_typing
type NewTypingEvent struct {
	SenderID   int  `json:"sender_id"`
	ReceiverID int  `json:"receiver_id"`
	IsTyping   bool `json:"is_typing"` // true = start typing, false = stop typing
}

// ClientList is a map of connected clients
type ClientList map[*Client]bool

// UserStatusInfo represents a user's status information
type UserStatusInfo struct {
	ID              int     `json:"id"`
	Username        string  `json:"username"`
	Status          int     `json:"status"`
	LastMessageTime *string `json:"last_message_time"`
}

// UserStatusUpdate represents the WebSocket event for status changes
type UserStatusUpdate struct {
	OnlineUsers    []UserStatusInfo `json:"online_users"`     // Full user objects for online users
	OfflineUserIDs []int            `json:"offline_user_ids"` // Just IDs for offline users
}
