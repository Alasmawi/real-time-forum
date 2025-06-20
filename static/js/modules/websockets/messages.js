import { SendMessageEvent, sendEvent } from './events.js';

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

// Export functions
export { sendMessage, appendChatMessage };
