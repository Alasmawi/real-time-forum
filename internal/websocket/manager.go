package websocket

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	// Upgrades incomming HTTP requests into a persitent websocket connection
	websocketUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

// Holds references to all registered clients and boadcasts messages to all clients
type WebsocketManager struct {
	clients ClientList

	// SyncMutex locks state before editing clients (channels can also be used to block)
	sync.RWMutex
}

// Initalizes all the values inside manager
func NewWebsocketManager() *WebsocketManager {
	return &WebsocketManager{
		clients: make(ClientList),
	}
}

// HTTP Handler that the has the Manager that allows connections
func (m *WebsocketManager) ServeWebSocket(w http.ResponseWriter, r *http.Request) {

	log.Println("New connection")
	// Begins by upgrading the HTTP request
	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// Creates new client
	client := NewClient(conn, m)
	// Adds newly created client to manager
	m.addClient(client)
	// Starts the read / write processes
	go client.readMessages()
	// go client.writeMessages()
}

/** Starts the reading client messages and handles them appropriatly.
Ran as a goroutine **/
func (c *Client) readMessages() {
	defer func() {
		// Gracefully closes connection once function is done
		c.manager.removeClient(c)
	}()
	// Infinite loop to keep reading messages
	for {
		// Read the next message in queue in the connection
		messageType, payload, err := c.connection.ReadMessage()

		if err != nil {
			/** Recieves an error if connection is lost
			only logs unexpected errors, not simple disconnection **/
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v", err)
			}
			break // Breaks the loop to close conn & cleanup
		}
		log.Println("MessageType: ", messageType)
		log.Println("Payload: ", string(payload))
	}
}
