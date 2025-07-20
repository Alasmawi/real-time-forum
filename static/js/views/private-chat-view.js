// Private Chat View - Handles chat modal HTML generation and UI
import { formatTime } from '../modules/utils/UI-UX/time-formatter.js';

export class PrivateChatView {
    
    // Generate the chat modal HTML
    static generateChatModal(user, position) {
        const modalId = `chat-modal-${user.id}`;
        
        return `
            <div id="${modalId}" class="private-chat-modal" 
                 style="left: ${position.left}px; top: ${position.top}px;">
                <div class="chat-modal-header">
                    <div class="chat-modal-user">
                        <div class="avatar medium ${this.getAvatarColor(user)}">
                            ${this.generateAvatar(user)}
                            <div class="status-dot ${user.online ? 'online' : 'offline'}"></div>
                        </div>
                        <span class="chat-modal-username">${user.username}</span>
                    </div>
                    <div class="chat-modal-controls">
                        <button class="chat-modal-close" title="Close">×</button>
                    </div>
                </div>
                <div class="chat-modal-messages" id="chat-messages-${user.id}">
                    <div class="chat-loading">Loading conversation...</div>
                </div>
                <div class="chat-modal-input">
                    <div class="chat-input-row">
                        <input type="text" 
                               class="chat-input" 
                               placeholder="Type a message..." 
                               id="chat-input-${user.id}"
                               maxlength="250">
                        <button class="chat-send-btn" id="chat-send-${user.id}">Send</button>
                    </div>
                    <div class="chat-counter" id="chat-counter-${user.id}">0/250 characters</div>
                    <div id="error-messages" class="general-error"></div>
                </div>
            </div>
        `;
    }

    // Generate avatar for user
    static generateAvatar(user) {
        // Use existing avatar generator if available
        if (window.AvatarGenerator) {
            return window.AvatarGenerator.generateAvatarLetter(user.username);
        }
        return user.username.charAt(0).toUpperCase();
    }

    // Get avatar color for user
    static getAvatarColor(user) {
        // Use existing avatar generator if available
        if (window.AvatarGenerator) {
            return window.AvatarGenerator.getColorFromUsername(user.username);
        }
        return 'primary';
    }

    // Generate message HTML
    static generateMessage(message, isCurrentUser) {
        const time = formatTime(message.timestamp);

        const messageType = isCurrentUser ? 'sent' : 'received';
        
        return `
            <div class="message ${messageType}">
                <div class="message-bubble ${messageType}">
                    ${this.escapeHtml(message.content)}
                    <div class="message-time">${time}</div>
                    ${message.status ? `<div class="message-status">${message.status}</div>` : ''}
                </div>
            </div>
        `;
    }

    // Generate typing indicator
    static generateTypingIndicator(username) {
        return `
            <div class="typing-indicator" id="typing-${username}">
                <div class="typing-text">
                    <span class="typing-dots">
                        <span></span>
                        <span></span>
                        <span></span>
                    </span>
                    ${username} is typing...
                </div>
            </div>
        `;
    }

    // Generate welcome message
    static generateWelcomeMessage() {
        return `
            <div class="chat-welcome">
                <p>Start your conversation!</p>
            </div>
        `;
    }

    // Calculate modal position relative to user element
    static calculateModalPosition(userElement) {
        if (!userElement) {
            return { left: 50, top: 50 }; // Default fallback position
        }

        const rect = userElement.getBoundingClientRect();
        const modalWidth = 350;
        const modalHeight = 400;
        const padding = 10;

        // Check if there's enough space to position next to user element
        const viewportWidth = window.innerWidth;
        const viewportHeight = window.innerHeight;

        // Try to position to the right of the user element
        let left = rect.right + padding;
        let top = rect.top;

        // If modal would go off right edge, position to the left
        if (left + modalWidth > viewportWidth) {
            left = rect.left - modalWidth - padding;
        }

        // If modal would go off left edge, center it horizontally
        if (left < 0) {
            left = (viewportWidth - modalWidth) / 2;
        }

        // Adjust vertical position if needed
        if (top + modalHeight > viewportHeight) {
            top = viewportHeight - modalHeight - padding;
        }

        // Ensure modal doesn't go above viewport
        if (top < padding) {
            top = padding;
        }

        // If viewport is too small, center the modal
        if (viewportWidth < modalWidth + 100 || viewportHeight < modalHeight + 100) {
            left = (viewportWidth - modalWidth) / 2;
            top = (viewportHeight - modalHeight) / 2;
        }

        return { left: Math.max(0, left), top: Math.max(0, top) };
    }

    // Escape HTML to prevent XSS
    static escapeHtml(text) {
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }

    // Update modal position (for dragging)
    static updateModalPosition(modal, x, y) {
        // Keep modal within viewport
        const maxX = window.innerWidth - modal.offsetWidth;
        const maxY = window.innerHeight - modal.offsetHeight;
        
        x = Math.max(0, Math.min(x, maxX));
        y = Math.max(0, Math.min(y, maxY));
        
        modal.style.left = x + 'px';
        modal.style.top = y + 'px';
    }

    // Show/hide typing indicator
    static updateTypingIndicator(messagesContainer, username, isTyping) {
        const existingIndicator = messagesContainer.querySelector('.typing-indicator');
        
        if (isTyping && !existingIndicator) {
            const indicator = document.createElement('div');
            indicator.innerHTML = this.generateTypingIndicator(username);
            messagesContainer.appendChild(indicator.firstElementChild);
            messagesContainer.scrollTop = messagesContainer.scrollHeight;
        } else if (!isTyping && existingIndicator) {
            existingIndicator.remove();
        }
    }
}