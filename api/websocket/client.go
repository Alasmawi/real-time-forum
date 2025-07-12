package websocket

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

var (
	// pongWait is how long the server will await a pong response from a client
	pongWait = 10 * time.Second
	// pingInterval has to be less than pongWait, otherwise it will send a new Ping before getting response
	pingInterval = (pongWait * 9) / 10 //90% of pongWait
)

type ClientList map[*Client]bool

type Client struct {
	connection *websocket.Conn
	manager *WebsocketManager

	// Egress is used to avoid concurrent writes on the WebSocket
	egress chan Event
	userID        int    `json:"user_id"`
	username      string `json:"username"`
	sessionToken  string `json:"session_token"`
	lastValidated time.Time
}

// Initializes a new c with all required values.
func NewClient(conn *websocket.Conn, manager *WebsocketManager, userID int, username, sessionToken string) *Client {
	return &Client{
		connection:    conn,
		manager:       manager,
		egress:        make(chan Event),
		userID:        userID,
		username:      username,
		sessionToken:  sessionToken,
		lastValidated: time.Now(), // Set initial validation time
	}
}

// pongHandler is used to handle PongMessages for the Client
func (c *Client) pongHandler(pongMsg string) error {
	// Current time + Pong Wait time
	log.Println("pong")
	return c.connection.SetReadDeadline(time.Now().Add(pongWait))
}

// Starts reading client messages and handles them appropriately.
// Runs as a goroutine.
func (c *Client) readMessages() {
	defer func() {
		// Gracefully closes the connection once the function is done.
		c.manager.removeClient(c)
	}()

	// Set Max Size of Messages in Bytes
	c.connection.SetReadLimit(512)

	// Configure Wait time for Pong response, use Current time + pongWait
	// This has to be done here to set the first initial timer.
	if err := c.connection.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Println(err)
		return
	}
	// Configure how to handle Pong responses
	c.connection.SetPongHandler(c.pongHandler)

	// Infinite loop to keep reading messages.
	for {
		// Reads the next message in queue in the connection.
		_, payload, err := c.connection.ReadMessage()

		if err != nil {
			// Receives an error if the connection is lost.
			// Only logs unexpected errors, not simple disconnection.
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v", err)
			}
			break
		}
		// Marshal incoming data into a Event struct
		var request Event
		if err := json.Unmarshal(payload, &request); err != nil {
			log.Printf("error marshalling message: %v", err)
			break // Breaking the connection here might be harsh xD
		}
		// Route the Event
		if err := c.manager.routeEvent(request, c); err != nil {
			log.Println("Error handeling Message: ", err)
		}
	}
}

// closeWithReason closes the WebSocket connection with proper close code and reason
func (c *Client) closeWithReason(code int, reason string) {
	closeData := websocket.FormatCloseMessage(code, reason)
	if err := c.connection.WriteMessage(websocket.CloseMessage, closeData); err != nil {
		log.Printf("Error sending close frame: %v", err)
	}
	c.connection.Close()
}

// sendErrorEvent sends an error event to the client
func (c *Client) sendErrorEvent(code, message string) {
	errorEvent := map[string]interface{}{
		"code":    code,
		"message": message,
		"type":    "error",
	}

	data, err := json.Marshal(errorEvent)
	if err != nil {
		log.Printf("Failed to marshal error event: %v", err)
		return
	}

	event := Event{
		Type:    "error",
		Payload: data,
	}

	select {
	case c.egress <- event:
	default:
		log.Printf("Client egress channel full, dropping error message")
	}
}

// writeMessages listens for new messages to output to the Client.
func (c *Client) writeMessages() {
	// Create a ticker that triggers a ping at given interval
	ticker := time.NewTicker(pingInterval)

	defer func() {
		ticker.Stop()

		// Gracefully closes the connection.
		c.manager.removeClient(c)
	}()

	for {
		select {
		case message, ok := <-c.egress:
			// Ok will be false Incase the egress channel is closed
			if !ok {
				// Manager has closed this connection channel, so communicate that to frontend with proper close code
				c.closeWithReason(websocket.CloseNormalClosure, "Normal closure")
				// Return to close the goroutine
				return
			}

			data, err := json.Marshal(message)
			if err != nil {
				log.Println(err)
				return // closes the connection, should we really
			}

			// Write a Regular text message to the connection
			if err := c.connection.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Println(err)
			}
			log.Println("sent message")

		case <-ticker.C:
			log.Println("ping")
			// Send the Ping
			if err := c.connection.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Println("writemsg: ", err)
				return // return to break this goroutine triggeing cleanup
			}
		}
	}
}
