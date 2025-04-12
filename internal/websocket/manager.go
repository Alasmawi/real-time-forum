package websocket

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	// Upgrades incomming HTTP requests into persitent websocket connections.
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

// checkOrigin will check origin and return true if its allowed
func checkOrigin(r *http.Request) bool {

	// Grab the request origin
	origin := r.Header.Get("Origin")

	switch origin {
	case "http://localhost:8080":
		return true
	default:
		return false
	}
}

// Holds references to all registered clients and boadcasts messages to all clients.
type WebsocketManager struct {
	clients ClientList

	// SyncMutex locks state before editing clients (channels can also be used to block).
	sync.RWMutex

	handlers map[string]EventHandler
}

// Initalizes all the values inside manager.
func NewWebsocketManager() *WebsocketManager {
	m := &WebsocketManager{
		clients:  make(ClientList),
		handlers: make(map[string]EventHandler),
	}
	m.setupEventHandlers()
	return m
}

// setupEventHandlers configures and adds all handlers
func (m *WebsocketManager) setupEventHandlers() {
	m.handlers[EventSendMessage] = func(e Event, c *Client) error {
		fmt.Println(e)
		return nil
	}
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
func (m *WebsocketManager) ServeWebSocket(w http.ResponseWriter, r *http.Request) {

	log.Println("New connection")
	// Begins by upgrading the HTTP request
	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// Creates new client.
	client := NewClient(conn, m)
	// Adds newly created client to manager.
	m.addClient(client)
	// Starts the read / write processes.
	go client.readMessages()
	go client.writeMessages()
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
