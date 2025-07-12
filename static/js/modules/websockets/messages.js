import { SendMessageEvent, sendEvent } from './events.js';

/**
* sendMessage will send a new message onto the Chat
**/
async function sendMessage() {
    const newmessage = document.getElementById("message");
    if (newmessage != null && newmessage.value.trim() !== '') {
        let outgoingEvent = new SendMessageEvent(newmessage.value, "percy");
        sendEvent("send_message", outgoingEvent);
        newmessage.value = ''; // Clear the message input
    }
    return false;
}

/**
* appendChatMessage takes in new messages and adds them to the chat
**/
function appendChatMessage(messageEvent) {
    const date = new Date(messageEvent.sent);
    // format message
    const formattedMsg = `${date.toLocaleString()}: ${messageEvent.message}`;
    // Append Message
    let textarea = document.getElementById("chatmessages");
    textarea.innerHTML = textarea.innerHTML + "\n" + formattedMsg;
    textarea.scrollTop = textarea.scrollHeight;
}

// Export functions
export { sendMessage, appendChatMessage };
