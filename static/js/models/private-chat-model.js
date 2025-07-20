// Private Chat Model - Manages chat data and WebSocket communication
import { getSessionValue } from '../modules/utils/session.js';

export class PrivateChatModel {
    constructor() {
        this.activeChats = new Map(); // userId -> chatData
        this.typingStates = new Map(); // userId -> isTyping
    }

    // Add a new chat
    addChat(userId, userData) {
        this.activeChats.set(userId, {
            user: userData,
            messages: [],
            unreadCount: 0,
            lastActivity: new Date()
        });
    }

    // Remove a chat
    removeChat(userId) {
        this.activeChats.delete(userId);
        this.typingStates.delete(userId);
    }

    // Get chat data
    getChat(userId) {
        return this.activeChats.get(userId);
    }

    // Check if chat exists
    hasChat(userId) {
        return this.activeChats.has(userId);
    }

    // Add message to chat
    addMessage(userId, message) {
        const chat = this.getChat(userId);
        if (chat) {
            chat.messages.push(message);
            chat.lastActivity = new Date();
            if (message.sender !== 'current_user') {
                chat.unreadCount++;
            }
        }
    }

    // Add message to top of chat (for older messages)
    addMessageToTop(userId, message) {
        const chat = this.getChat(userId);
        if (chat) {
            chat.messages.unshift(message);
            chat.lastActivity = new Date();
        }
    }

    // Clear messages for a chat
    clearMessages(userId) {
        const chat = this.getChat(userId);
        if (chat) {
            chat.messages = [];
        }
    }

    // Get messages for a chat
    getMessages(userId) {
        const chat = this.getChat(userId);
        return chat ? chat.messages : [];
    }

    // Mark chat as read
    markAsRead(userId) {
        const chat = this.getChat(userId);
        if (chat) {
            chat.unreadCount = 0;
        }
    }

    // Set typing status
    setTypingStatus(userId, isTyping) {
        this.typingStates.set(userId, isTyping);
    }

    // Get typing status
    getTypingStatus(userId) {
        return this.typingStates.get(userId) || false;
    }

    // Get all active chats
    getAllChats() {
        return Array.from(this.activeChats.entries()).map(([userId, chatData]) => ({
            userId,
            ...chatData
        }));
    }

    // Send message via WebSocket
    async sendMessage(receiverId, message) {
        try {
            // Get session token using session utilities
            const sessionToken = await getSessionValue();
            if (!sessionToken) {
                throw new Error('No session token found');
            }

            // Import WebSocket functionality
            const { SendMessageEvent, sendEvent } = await import('../modules/websockets/events.js');
            
            // Create private message event
            const messageEvent = new SendMessageEvent(message, receiverId, sessionToken);
            
            // Send via WebSocket
            sendEvent("send_message", messageEvent);
            
            return true;
        } catch (error) {
            console.error('Error sending private message:', error);
            return false;
        }
    }

    // Send typing status via WebSocket
    async sendTypingStatus(receiverId, isTyping) {
        try {
            const sessionToken = await getSessionValue();
            if (!sessionToken) return false;

            const { SendTypingEvent, sendEvent } = await import('../modules/websockets/events.js');
            const typingEvent = new SendTypingEvent(receiverId, isTyping, sessionToken);
            sendEvent("send_typing", typingEvent);
            
            return true;
        } catch (error) {
            console.error('Error sending typing status:', error);
            return false;
        }
    }

}