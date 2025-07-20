// Events module - handles WebSocket event routing for private messaging

// WebSocket connection variable (set by the WebSocket manager)
let conn;

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
  constructor (message, receiver_id, session_token) {
    this.message = message
    this.receiver_id = receiver_id
    this.session_token = session_token
  }
}
/**
 * ReceiveMessageEvent is messages coming from backend
 **/
class ReceiveMessageEvent {
  constructor (message, sender_id, receiver_id, sent_at) {
    this.message = message
    this.sender_id = sender_id
    this.receiver_id = receiver_id
    this.sent_at = sent_at
  }
}


/**
 * SendTypingEvent is used to send typing status to other users
 **/
class SendTypingEvent {
  constructor (receiver_id, is_typing, session_token) {
    this.receiver_id = receiver_id
    this.is_typing = is_typing
    this.session_token = session_token
  }
}

/**
 * NewTypingEvent is received from backend for typing status updates
 **/
class NewTypingEvent {
  constructor (sender_id, receiver_id, is_typing) {
    this.sender_id = sender_id
    this.receiver_id = receiver_id
    this.is_typing = is_typing
  }
}

/**
 * Event type constants for better maintainability
 */
const eventTypes = {
  receiveMessage: 'receive_message',
  sendMessage: 'send_message',
  sendTyping: 'send_typing',
  newTyping: 'new_typing',
  userStatusUpdate: 'user_status_update',
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
  
  console.log('WebSocket routeEvent called with event:', event);
  console.log('Event type:', event.type);
  console.log('Event payload:', event.payload);
  console.log('Event payload type:', typeof event.payload);
  
  // Parse payload if it's a string (backend double-encodes JSON)
  let payload = event.payload;
  if (typeof payload === 'string') {
    try {
      payload = JSON.parse(payload);
      console.log('Parsed payload from string:', payload);
    } catch (e) {
      console.error('Failed to parse payload JSON:', e);
      console.error('Raw payload:', event.payload);
    }
  }
  
  switch (event.type) {
    case eventTypes.receiveMessage:
      // Format payload
      const messageEvent = Object.assign(new ReceiveMessageEvent(), payload)
      
      // Route to private chat controller
      if (window.privateChatController) {
        window.privateChatController.receiveMessage(messageEvent.sender_id, messageEvent.receiver_id, {
          content: messageEvent.message,
          timestamp: messageEvent.sent_at
        });
      }
      
      // Add notification if message is not for currently open chat
      if (window.privateChatController && window.NotificationsController) {
        const isCurrentChat = window.privateChatController.isCurrentChat && 
                              window.privateChatController.isCurrentChat(messageEvent.sender_id);
        if (!isCurrentChat) {
          // Get sender username from user list
          const senderUser = window.userListController ? 
                           window.userListController.model.getUserById(messageEvent.sender_id) : null;
          const senderUsername = senderUser ? senderUser.username : `User ${messageEvent.sender_id}`;
          
          window.NotificationsController.addNotification(messageEvent.sender_id, senderUsername, 1);
        }
      }
      break
    case eventTypes.userStatusUpdate:
      console.log('Routing to userStatusUpdate handler');
      handleUserStatusUpdate(payload)
      break
    case eventTypes.newTyping:
      handleNewTyping(payload)
      break
    case eventTypes.error:
      console.log('Error event received, calling handleError with payload:', payload)
      handleError(payload)
      break
    default:
      console.warn('unsupported action:', event.type)
  }
}

/**
 * Handle user status updates (new efficient system)
 */
function handleUserStatusUpdate(payload) {
  console.log('WebSocket handleUserStatusUpdate called with payload:', payload)
  console.log('Payload type:', typeof payload)
  console.log('Payload keys:', Object.keys(payload))
  console.log('Online users:', payload.online_users)
  console.log('Offline user IDs:', payload.offline_user_ids)
  console.log('Online users count:', payload.online_users?.length || 0)
  console.log('Offline user IDs count:', payload.offline_user_ids?.length || 0)
  console.log('userListController available?', !!window.userListController)
  
  // Update user list UI if controller is available
  if (window.userListController) {
    console.log('Calling userListController.handleUserStatusUpdate with:', payload)
    window.userListController.handleUserStatusUpdate(payload);
  } else {
    console.error('userListController not available for status update')
    console.error('Available window objects:', Object.keys(window).filter(key => key.includes('Controller')))
  }
}

/**
 * Handle new typing events
 */
function handleNewTyping(payload) {
  // Payload is already an object, no need to parse
  const typingEvent = payload
  
  // Route typing status to private chat controller if available
  if (window.privateChatController) {
    window.privateChatController.handleTypingStatus(typingEvent.sender_id, typingEvent.is_typing);
  }
}

/**
 * Handle error events
 */
function handleError(payload) {
  // Payload is already an object, no need to parse
  const error = payload
  console.error('WebSocket error:', error)
  console.log('Error code:', error.code)
  console.log('Error message:', error.message)
  console.log('Error receiver_id:', error.receiver_id)
  
  switch (error.code) {
    case 'RECEIVER_OFFLINE':
      showInlineError(error.receiver_id, 'The user you are trying to message is not online.')
      break
    case 'SELF_MESSAGE':
      showInlineError(error.receiver_id, 'You cannot send messages to yourself.')
      break
    case 'VALIDATION_ERROR':
      showInlineError(error.receiver_id, error.message)
      break
    case 'SELF_TYPING':
      // Silent ignore for typing to self
      break
    default:
      // For errors without receiver_id context, show generic alert
      alert(`Error: ${error.message}`)
      break
  }
}

/**
 * Show inline error message under the message input box
 */
function showInlineError(receiverId, message) {
  console.log('showInlineError called with receiverId:', receiverId, 'message:', message)
  
  if (!receiverId) {
    // Fallback to alert if no receiver context
    alert(message)
    return
  }

  const chatModal = document.getElementById(`chat-modal-${receiverId}`)
  if (!chatModal) {
    console.warn(`Chat modal not found for user ${receiverId}`)
    return
  }

  // Ensure error container exists with same structure as forms
  let errorContainer = chatModal.querySelector('#error-messages')
  if (!errorContainer) {
    const chatModalInput = chatModal.querySelector('.chat-modal-input')
    
    if (!chatModalInput) {
      console.warn(`Chat input elements not found for user ${receiverId}`)
      return
    }

    // Create error container using same structure as FormErrorHelpers
    errorContainer = document.createElement('div')
    errorContainer.id = 'error-messages'
    errorContainer.className = 'general-error'
    
    chatModalInput.appendChild(errorContainer)
  }

  // Use FormErrorHelpers smart error handling
  if (window.FormErrorHelpers) {
    window.FormErrorHelpers.handleValidationErrors({
      Errors: [message]
    })
  } else {
    // Fallback if FormErrorHelpers not available
    errorContainer.innerHTML = message
  }

  // Auto-hide after 5 seconds
  setTimeout(() => {
    if (window.FormErrorHelpers) {
      window.FormErrorHelpers.clearAllErrors()
    } else if (errorContainer) {
      errorContainer.innerHTML = ''
    }
  }, 5000)
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
 * Set the WebSocket connection (called by WebSocket manager)
 */
function setConnection(connection) {
    conn = connection;
}

// Export all classes and functions
export { Event, SendMessageEvent, ReceiveMessageEvent, SendTypingEvent, NewTypingEvent, eventTypes, routeEvent, sendEvent, setConnection };
