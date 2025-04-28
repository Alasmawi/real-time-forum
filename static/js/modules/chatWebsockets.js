//Javascript that is used to Connect to Websocket and Handle New messages

//declaring conn as a global variable
let conn;

// selectedchat is by default General.
var selectedchat = "general";

/**
* Event is used to wrap all messages Send and Recieved
* on the Websocket
* The type is used as a RPC
**/
class Event {
    // Each Event needs a Type
    // The payload is not required
    constructor(type, payload) {
        this.type = type;
        this.payload = payload;
    }
}

/**
* SendMessageEvent is used to send messages to other clients
**/
class SendMessageEvent {
    constructor(message, from) {
        this.message = message;
        this.from = from;
    }
}
/**
* NewMessageEvent is messages comming from clients
**/
class NewMessageEvent {
    constructor(message, from, sent) {
        this.message = message;
        this.from = from;
        this.sent = sent;
    }
}

/**
* ChangeChatRoomEvent is used to switch chatroom
**/
class ChangeChatRoomEvent {
    constructor(name) {
        this.name = name;
    }
}

/**
* routeEvent is a proxy function that routes
* events into their correct Handler
* based on the type field
**/
function routeEvent(event) {

    if (event.type === undefined) {
        alert("no 'type' field in event");
    }
    switch (event.type) {
        case "new_message":
            // Format payload
            const messageEvent = Object.assign(new NewMessageEvent, event.payload);
            appendChatMessage(messageEvent);
            break;
        default:
            alert("unsupported message type");
            break;
    }
}

/**
 * Updates the value of selectedchat
 * and notifies the server when chatroom is changed
 **/
function changeChatRoom() {
    // Change Header to reflect the Changed chatroom
    var newchat = document.getElementById("chatroom");
    if (newchat != null && newchat.value != selectedchat) {
        selectedchat = newchat.value;

        let header = document.getElementById("chat-header");
        header.innerHTML = "Currently in chat: " + selectedchat;

        let changeEvent = new ChangeChatRoomEvent(selectedchat);
        sendEvent("change_room", changeEvent);
        
        let textarea = document.getElementById("chatmessages");
        textarea.innerHTML = `You changed room into: ${selectedchat}`;
    }
    return false;
}

/**
* sendMessage will send a new message onto the Chat
**/
function sendMessage() {
    var newmessage = document.getElementById("message");
    if (newmessage != null) {
        let outgoingEvent = new SendMessageEvent(newmessage.value, "percy");
        sendEvent("send_message", outgoingEvent)
    }
    return false;
}

/**
* sendEvent
* eventname - the event name to send on
* payload - the data payload
**/
function sendEvent(eventName, payload) {
    // Create a event Object with a event named send_message
    const event = new Event(eventName, payload);
    // Format as JSON and send
    conn.send(JSON.stringify(event));
}

/**
* appendChatMessage takes in new messages and adds them to the chat
**/
function appendChatMessage(messageEvent) {
    var date = new Date(messageEvent.sent);
    // format message
    const formattedMsg = `${date.toLocaleString()}: ${messageEvent.message}`;
    // Append Message
    let textarea = document.getElementById("chatmessages");
    textarea.innerHTML = textarea.innerHTML + "\n" + formattedMsg;
    textarea.scrollTop = textarea.scrollHeight;
}

/**
* Once the website loads, we want to apply listeners and connect to websocket
**/
function chatFunctionality() {
    // Apply our listener functions to the submit event on both forms
    // we do it this way to avoid redirects
    document.getElementById("chatroom-selection").onsubmit = changeChatRoom;
    document.getElementById("chatroom-message").onsubmit = sendMessage;

    // Check if the browser supports WebSocket0........................0000000000000000000000000000.0
    if (window["WebSocket"]) {
        console.log("supports websockets");

        try {
            // Initialize the WebSocket connection
            conn = new WebSocket("ws://" + document.location.host + "/ws");

            conn.onopen = function () {
                console.log("WebSocket connection established.");
            };

            // Add a listener to the onmessage event
            conn.onmessage = function (e) {
                console.log(e);
                // parse websocket message as JSON
                const eventData = JSON.parse(e.data);
                // Assign JSON data to new Event Object
                const event = Object.assign(new Event, eventData);
                // Let router manage message
                routeEvent(event);
            };

            conn.onclose = function () {
                console.log("WebSocket connection closed.");
            };

            conn.onerror = function (error) {
                console.error("WebSocket error:", error);
            };
        } catch (error) {
            console.error("Error initializing WebSocket:", error);
        }
    } else {
        alert("Not supporting websockets");
    }
};

export { chatFunctionality };