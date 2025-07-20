// Private Messaging WebSocket - Simple WebSocket connection integrated with private chat MVC
import { routeEvent, setConnection, Event } from './events.js';

// WebSocket connection variable
let conn;
let reconnectInterval = null;

// Update user profile WebSocket status
function updateUserProfileStatus(status, className = '') {
    const userStatus = document.querySelector('.user-profile .user-status');
    if (userStatus) {
        userStatus.textContent = status;
        userStatus.className = `user-status ${className}`;
    }
}

// Connect/reconnect WebSocket
function connectWebSocket() {
    // Check if the browser supports WebSockets
    if (!window["WebSocket"]) {
        alert("Browser does not support WebSockets");
        updateUserProfileStatus("Offline", "offline");
        return;
    }

    // Update status to connecting
    updateUserProfileStatus("Connecting...", "connecting");

    try {
        // Create WebSocket connection to the same endpoint
        conn = new WebSocket("ws://" + document.location.host + "/protected/ws");

        conn.onopen = function () {
            console.log("Private messaging WebSocket connected");
            updateUserProfileStatus("Online", "online");
            // Clear any existing reconnect interval
            if (reconnectInterval) {
                clearInterval(reconnectInterval);
                reconnectInterval = null;
            }
        };

        // Handle incoming messages
        conn.onmessage = function (e) {
            // Parse websocket message as JSON
            const eventData = JSON.parse(e.data);
            // Assign JSON data to new Event Object
            const event = Object.assign(new Event, eventData);
            // Let router manage message - this will call the private chat controller methods
            routeEvent(event);
        };

        conn.onclose = function (e) {
            console.log("Private messaging WebSocket closed. Code: " + e.code + ", reason: " + e.reason);
            updateUserProfileStatus("Disconnected", "offline");
            // Start auto-reconnect if not already running
            startReconnect();
        };

        conn.onerror = function (error) {
            console.error("Private messaging WebSocket error:", error);
            updateUserProfileStatus("Connection Failed", "failed");
            // Start auto-reconnect on error too
            startReconnect();
        };

        // Set connection in events module so sendEvent can use it
        setConnection(conn);

    } catch (error) {
        console.error("Error initializing private messaging WebSocket:", error);
        updateUserProfileStatus("Connection Failed", "failed");
    }
}

// Initialize WebSocket connection for private messaging
function initializePrivateMessaging() {
    connectWebSocket();
}

// Auto-reconnect function
function startReconnect() {
    // Don't start multiple reconnect intervals
    if (reconnectInterval) {
        return;
    }
    
    console.log("Starting WebSocket auto-reconnect (every 10 seconds)");
    reconnectInterval = setInterval(() => {
        console.log("Attempting WebSocket reconnection...");
        connectWebSocket();
    }, 10000); // 10 seconds
}

export { initializePrivateMessaging };