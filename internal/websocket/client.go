package websocket

import "github.com/gorilla/websocket"

// A map used to help manage clients
type ClientList map[*Client]bool

// Represents a websocket client, basically a frontend visitor
type Client struct {
	// Websocket connection
	connection *websocket.Conn

	// Manager used to manage the client
	manager *WebsocketManager
}

// Initializes a new client with all required values
func NewClient(conn *websocket.Conn, manager *WebsocketManager) *Client {
	return &Client{
		connection: conn,
		manager:    manager,
	}
}

/** Locks in addClient and removeClient to ensure thread safety, 
as the functions may be invoked concurrently **/

// Add client to clientList
func (m *WebsocketManager) addClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	// Add Client
	m.clients[client] = true
}

// Removes client and cleans up
func (m *WebsocketManager) removeClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	// Checks if client exists, then deletes it
	if _, ok := m.clients[client]; ok {
		client.connection.Close()
		delete(m.clients, client)
	}
}
