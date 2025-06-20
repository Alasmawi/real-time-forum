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
  constructor (message, from) {
    this.message = message
    this.from = from
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
  changeRoom: 'change_room'
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
    default:
      console.warn('unsupported action')
  }
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
export { Event, SendMessageEvent, NewMessageEvent, ChangeChatRoomEvent, eventTypes, routeEvent, sendEvent };
