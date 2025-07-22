package websocket

import (
	"errors"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	db "reboot01.com/js/realtime-forum/internal/database"
)

var (
	// Upgrades incoming HTTP requests into persistent websocket connections.
	websocketUpgrader = websocket.Upgrader{
		// CSRF check
		CheckOrigin:     checkOrigin,
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

var (
	ErrEventNotSupported = errors.New("this event type is not supported")
)

// checkOrigin will check origin and return true if it's allowed
func checkOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")

	switch origin {
	case "http://localhost:8080":
		return true
	default:
		return false
	}
}

// Holds references to all registered clients and broadcasts messages to all clients.
type WebsocketManager struct {
	clients ClientList
	DB      *db.DB
	// SyncMutex locks state before editing clients (channels can also be used to block).
	sync.RWMutex

	handlers map[string]EventHandler
	// userListTicker for periodic user list updates
	userListTicker *time.Ticker
	stopTicker     chan bool
	// previousOnlineUsers for change detection
	previousOnlineUsers map[int]UserStatusInfo
}

// Initializes all the values inside manager.
func NewWebsocketManager() *WebsocketManager {
	m := &WebsocketManager{
		clients:             make(ClientList),
		handlers:            make(map[string]EventHandler),
		stopTicker:          make(chan bool),
		previousOnlineUsers: make(map[int]UserStatusInfo),
	}
	m.setupEventHandlers()
	m.startPeriodicUserListBroadcast()
	return m
}

// setupEventHandlers configures and adds all handlers
func (m *WebsocketManager) setupEventHandlers() {
	m.handlers[EventSendMessage] = SendMessageHandler
	m.handlers[EventSendTyping] = SendTypingHandler
	// Note: EventRequestUserList removed - server now broadcasts periodically
}

// routeEvent is used to make sure the correct event goes into the correct handler
func (m *WebsocketManager) routeEvent(event Event, c *Client) error {
	// Check if Handler is present in Map
	if handler, ok := m.handlers[event.Type]; ok {
		// Execute the handler and return any err
		if err := handler(event, c); err != nil {
			return err
		}
		return nil
	} else {
		return ErrEventNotSupported
	}
}

// HTTP Handler that the has the Manager that allows connections.
func (m *WebsocketManager) HttpToWebsocket(w http.ResponseWriter, r *http.Request, userID int, username, sessionID string) error {
	// Begins by upgrading the HTTP request
	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return err
	}

	// Creates new client with user info.
	client := NewClient(conn, m, userID, username, sessionID)
	// Adds newly created client to manager.
	m.addClient(client)
	// Send initial status update to new client
	go m.sendInitialStatusUpdate(client)
	// Starts the read / write processes.
	go client.readMessages()
	go client.writeMessages()

	return nil
}

// Locks in addClient and removeClient ensure thread safety,
// as the functions may be invoked concurrently.

// Adds a client to the clientList.
func (m *WebsocketManager) addClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	// Add Client
	m.clients[client] = true
}

// Removes client and cleans up.
func (m *WebsocketManager) removeClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	// Checks if client exists, then deletes it.
	if _, ok := m.clients[client]; ok {
		client.connection.Close()
		delete(m.clients, client)
	}
}
