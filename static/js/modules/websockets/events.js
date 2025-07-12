import { appendChatMessage } from './messages.js';

/**
 * Event is used to wrap all messages Send and Recieved
 * on the Websocket
 * The type is used as a RPC
 **/
class Event {
  // Each Event needs a Type
  // The payload is not required
  constructor (type, payload) {
    this.type = type
    this.payload = payload
  }
}

/**
 * SendMessageEvent is used to send messages to other clients
 **/
class SendMessageEvent {
  constructor (message, sender_id, receiver_id, session_token) {
    this.message = message
    this.sender_id = sender_id
    this.receiver_id = receiver_id
    this.session_token = session_token
  }
}
/**
 * NewMessageEvent is messages comming from clients
 **/
class NewMessageEvent {
  constructor (message, from, sent) {
    this.message = message
    this.from = from
    this.sent = sent
  }
}

/**
 * ChangeChatRoomEvent is used to switch chatroom
 **/
class ChangeChatRoomEvent {
  constructor (name) {
    this.name = name
  }
}

/**
 * Event type constants for better maintainability
 */
const eventTypes = {
  newMessage: 'new_message',
  sendMessage: 'send_message',
  requestUserList: 'request_user_list',
  userListUpdate: 'user_list_update',
  userStatusChange: 'user_status_change',
  error: 'error'
}

/**
 * routeEvent is a proxy function that routes
 * events into their correct Handler
 * based on the type field
 **/
function routeEvent (event) {
  if (event.type === undefined) {
    alert("no 'type' field in event")
  }
  switch (event.type) {
    case eventTypes.newMessage:
      // Format payload
      const messageEvent = Object.assign(new NewMessageEvent(), event.payload)
      appendChatMessage(messageEvent)
      break
    case eventTypes.userListUpdate:
      handleUserListUpdate(event.payload)
      break
    case eventTypes.userStatusChange:
      handleUserStatusChange(event.payload)
      break
    case eventTypes.error:
      handleError(event.payload)
      break
    default:
      console.warn('unsupported action:', event.type)
  }
}

/**
 * Handle user list updates
 */
function handleUserListUpdate(payload) {
  const userList = JSON.parse(payload)
  console.log('User list updated:', userList)
  // TODO: Update user list UI
}

/**
 * Handle user status changes (online/offline)
 */
function handleUserStatusChange(payload) {
  const statusChange = JSON.parse(payload)
  console.log('User status change:', statusChange)
  // TODO: Update user list UI
}

/**
 * Handle error events
 */
function handleError(payload) {
  const error = JSON.parse(payload)
  console.error('WebSocket error:', error)
  
  if (error.code === 'RECEIVER_OFFLINE') {
    alert('The user you are trying to message is not online.')
  } else if (error.code === 'SELF_MESSAGE') {
    alert('You cannot send messages to yourself.')
  } else {
    alert(`Error: ${error.message}`)
  }
}

/**
 * Request user list from server
 */
async function requestUserList() {
  const sessionToken = await getSessionToken()
  if (!sessionToken) {
    console.error('No session token found')
    return
  }
  
  sendEvent(eventTypes.requestUserList, { session_token: sessionToken })
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

// Export all classes and functions
export { Event, SendMessageEvent, NewMessageEvent, ChangeChatRoomEvent, eventTypes, routeEvent, sendEvent, requestUserList};
