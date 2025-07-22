package websocket

import (
	"time"
)

// startPeriodicUserListBroadcast starts a goroutine that broadcasts user status changes every 5 seconds
func (m *WebsocketManager) startPeriodicUserListBroadcast() {
	m.userListTicker = time.NewTicker(5 * time.Second)
	go func() {
		for {
			select {
			case <-m.userListTicker.C:
				m.broadcastUserList()
			case <-m.stopTicker:
				m.userListTicker.Stop()
				return
			}
		}
	}()
}

// stopPeriodicUserListBroadcast stops the periodic user list broadcast
func (m *WebsocketManager) stopPeriodicUserListBroadcast() {
	if m.userListTicker != nil {
		m.stopTicker <- true
	}
}

// broadcastUserList sends user status changes to all connected clients
func (m *WebsocketManager) broadcastUserList() {
	if len(m.clients) == 0 {
		return
	}

	// Detect changes in online users
	statusUpdate := m.detectUserStatusChanges()

	// Only broadcast if there are changes
	if len(statusUpdate.OnlineUsers) == 0 && len(statusUpdate.OfflineUserIDs) == 0 {
		return
	}

	// Broadcast to all clients
	m.RLock()
	clients := make([]*Client, 0, len(m.clients))
	for client := range m.clients {
		clients = append(clients, client)
	}
	m.RUnlock()

	for _, client := range clients {
		go m.sendUserStatusUpdate(client, statusUpdate)
	}
}
