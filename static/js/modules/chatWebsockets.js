//Javascript that is used to Connect to Websocket and Handle New messages

//declaring conn as a global variable
let conn;

// selectedchat is by default General.
var selectedchat = "general";

/**
 * Updates the value of selectedchat
 * and notifies the server when chatroom is changed
 * */
function changeChatRoom() {
    // Change Header to reflect the Changed chatroom
    var newchat = document.getElementById("chatroom");
    if (newchat != null && newchat.value != selectedchat) {
        console.log(newchat);
        selectedchat = newchat.value;
    }
    return false;
}

/**
 * Sends a new message onto the Websocket
 **/
function sendMessage() {
    var newmessage = document.getElementById("message");
    if (newmessage != null) {
        conn.send(newmessage.value);
    }
    return false;
}

/**
  * Once the website loads, we want to apply listeners and connect to websocket
  * */
function chatFunctionality() {
    // Apply our listener functions to the submit event on both forms
    // we do it this way to avoid redirects
    document.getElementById("chatroom-selection").onsubmit = changeChatRoom;
    document.getElementById("chatroom-message").onsubmit = sendMessage;

    // Check if the browser supports WebSocket
    if (window["WebSocket"]) {
        console.log("supports websockets");

        try {
            // Initialize the WebSocket connection
            conn = new WebSocket("ws://" + document.location.host + "/ws");
    
            conn.onopen = () => {
                console.log("WebSocket connection established.");
            };
    
            conn.onerror = (error) => {
                console.error("WebSocket connection error:", error);
            };
    
            conn.onclose = () => {
                console.log("WebSocket connection closed.");
            };
        } catch (error) {
            console.error("Error initializing WebSocket:", error);
        }
    } else {
        alert("Not supporting websockets");
    }
};

export { chatFunctionality };