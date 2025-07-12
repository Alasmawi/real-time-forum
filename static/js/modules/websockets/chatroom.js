//Javascript that is used to Connect to Websocket and Handle New messages

import { ChangeChatRoomEvent, sendEvent, routeEvent, Event } from './events.js';
import { sendMessage } from './messages.js';

//declaring conn as a module variable
let conn;

// selectedchat is by default General.
var selectedchat = "general";

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
* Once the website loads, we want to apply listeners and connect to websocket
**/
function chatFunctionality() {
    // Apply our listener functions to the submit event on both forms
    // we do it this way to avoid redirects
    document.getElementById("chatroom-selection").onsubmit = changeChatRoom;
    document.getElementById("chatroom-message").onsubmit = async function (e) {
        e.preventDefault();
        await sendMessage();
        return false;
    };

    // Check if the browser supports WebSockets
    if (window["WebSocket"]) {
        console.log("supports websockets");
    } else {
        alert("Browser does not support WebSockets");
    }

    try {
        // Instantializes the WebSocket connection and passes the session token to the server
        conn = new WebSocket("ws://" + document.location.host + "/protected/ws");

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

        conn.onclose = function (e) {
            console.log("WebSocket connection closed. Error code: " + e.code + ", reason: " + e.reason + ", wasClean: " + e.wasClean);
        };

        conn.onerror = function (error) {
            console.error("WebSocket error:", error);
        };
    } catch (error) {
        console.error("Error initializing WebSocket:", error);
    }
};

export { chatFunctionality };
