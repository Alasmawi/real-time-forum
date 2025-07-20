// Private Chat Controller - Orchestrates private chat functionality
import { PrivateChatModel } from '../models/private-chat-model.js';
import { PrivateChatView } from '../views/private-chat-view.js';
import PaginationManager from '../modules/pagination/pagination-manager.js';
import { CharacterCounter } from '../modules/utils/character-counter.js';

export class PrivateChatController {
    constructor() {
        this.model = new PrivateChatModel();
        this.zIndexCounter = 1000;
        this.typingTimeouts = new Map(); // userId -> timeout
        this.paginationManager = new PaginationManager();
        this.currentUserId = null;
        this.init();
    }

    // Initialize controller
    async init() {
        this.createModalContainer();
        await this.fetchCurrentUser();
    }

    // Fetch current user ID
    async fetchCurrentUser() {
        try {
            const response = await fetch('/protected/v1/user/me', {
                credentials: 'include'
            });
            if (response.ok) {
                const userData = await response.json();
                this.currentUserId = userData.id;
            } else {
                console.error(`Failed to fetch current user: ${response.status}`);
            }
        } catch (error) {
            console.error('Failed to fetch current user:', error);
        }
    }

    // Create modal container if it doesn't exist
    createModalContainer() {
        if (!document.getElementById('chat-modals-container')) {
            const container = document.createElement('div');
            container.id = 'chat-modals-container';
            document.body.appendChild(container);
        }
    }

    // Open chat with a user
    openChat(user, userElement) {
        // Check if this specific chat is already open
        const existingModal = document.getElementById(`chat-modal-${user.id}`);
        if (existingModal && this.model.hasChat(user.id)) {
            this.focusChat(user.id);
            return;
        }

        // Close any existing chat modals first (for single chat mode)
        this.closeAllChats();

        // Clear notifications for this user
        if (window.NotificationsController) {
            window.NotificationsController.clearNotification(user.id);
        }

        // Add chat to model
        this.model.addChat(user.id, user);

        // Create and show modal
        this.createChatModal(user, userElement);
    }

    // Close all existing chat modals
    closeAllChats() {
        const existingModals = document.querySelectorAll('.private-chat-modal');
        existingModals.forEach(modal => {
            const userId = modal.id.replace('chat-modal-', '');
            if (userId) {
                this.closeChat(parseInt(userId));
            }
        });
    }

    // Create chat modal
    createChatModal(user, userElement) {
        const position = PrivateChatView.calculateModalPosition(userElement);
        const modalHtml = PrivateChatView.generateChatModal(user, position);
        
        // Create modal element
        const container = document.getElementById('chat-modals-container');
        container.insertAdjacentHTML('beforeend', modalHtml);
        
        const modal = document.getElementById(`chat-modal-${user.id}`);
        modal.style.zIndex = ++this.zIndexCounter;

        // Attach event listeners
        this.attachModalEvents(modal, user);

        // Make modal draggable
        this.makeDraggable(modal);

        // Add click-outside-to-close functionality
        this.attachClickOutsideHandler(modal, user.id);

        // Load chat history with pagination
        this.loadChatHistory(user.id);


        return modal;
    }

    // Attach event listeners to modal
    attachModalEvents(modal, user) {
        const closeBtn = modal.querySelector('.chat-modal-close');
        const input = modal.querySelector(`#chat-input-${user.id}`);
        const sendBtn = modal.querySelector(`#chat-send-${user.id}`);

        // Close modal
        closeBtn.addEventListener('click', () => {
            this.closeChat(user.id);
        });

        // Focus modal when clicked
        modal.addEventListener('click', () => {
            this.focusChat(user.id);
        });

        // Send message on Enter
        input.addEventListener('keypress', (e) => {
            if (e.key === 'Enter' && !e.shiftKey) {
                e.preventDefault();
                this.sendMessage(user.id);
            }
        });

        // Send message on button click
        sendBtn.addEventListener('click', () => {
            this.sendMessage(user.id);
        });

        // Typing detection - throttle 1s, debounce 3s
        let typingTimer;
        let throttleTimer;
        let isThrottled = false;
        
        input.addEventListener('input', () => {
            // Clear debounce timer
            clearTimeout(typingTimer);
            
            // Throttle: only send typing event once per second
            if (!isThrottled) {
                console.log('Sending typing started for user:', user.id);
                this.sendTypingStatus(user.id, true);
                isThrottled = true;
                
                throttleTimer = setTimeout(() => {
                    isThrottled = false;
                }, 1000);
            }
            
            // Debounce: stop typing after 3 seconds of inactivity
            typingTimer = setTimeout(() => {
                console.log('Sending typing stopped for user:', user.id);
                this.sendTypingStatus(user.id, false);
            }, 3000);
        });

        // Stop typing when input loses focus
        input.addEventListener('blur', () => {
            clearTimeout(typingTimer);
            clearTimeout(throttleTimer);
            isThrottled = false;
            this.sendTypingStatus(user.id, false);
        });

        // Initialize character counter
        const counterElement = modal.querySelector(`#chat-counter-${user.id}`);
        if (counterElement) {
            CharacterCounter.initialize(input, counterElement, {
                maxLength: 250,
                warningThreshold: 200,
                dangerThreshold: 230
            });
        }
    }

    // Make modal draggable
    makeDraggable(modal) {
        const header = modal.querySelector('.chat-modal-header');
        let isDragging = false;
        let initialX, initialY;
        let lastTime = 0;
        let animationFrame;

        header.addEventListener('mousedown', (e) => {
            isDragging = true;
            initialX = e.clientX - modal.offsetLeft;
            initialY = e.clientY - modal.offsetTop;
            modal.style.cursor = 'grabbing';
            modal.style.transition = 'none';
            lastTime = performance.now();
            this.focusChat(modal.id.replace('chat-modal-', ''));
            e.preventDefault();
        });

        document.addEventListener('mousemove', (e) => {
            if (isDragging) {
                if (animationFrame) {
                    cancelAnimationFrame(animationFrame);
                }
                
                const currentTime = performance.now();
                const deltaTime = currentTime - lastTime;
                
                // Only update if enough time has passed (throttle to ~60fps equivalent)
                if (deltaTime >= 16.67) {
                    animationFrame = requestAnimationFrame(() => {
                        const x = e.clientX - initialX;
                        const y = e.clientY - initialY;
                        
                        modal.style.left = x + 'px';
                        modal.style.top = y + 'px';
                    });
                    lastTime = currentTime;
                }
            }
        });

        document.addEventListener('mouseup', () => {
            if (isDragging) {
                isDragging = false;
                modal.style.cursor = 'default';
                modal.style.transition = '';
                if (animationFrame) {
                    cancelAnimationFrame(animationFrame);
                }
            }
        });
    }

    // Add click-outside-to-close functionality
    attachClickOutsideHandler(modal, userId) {
        const handleClickOutside = (event) => {
            // Check if clicking on the same user item (don't close)
            const clickedUserItem = event.target.closest('.user-item');
            if (clickedUserItem) {
                const clickedUserId = parseInt(clickedUserItem.getAttribute('data-user-id'));
                if (clickedUserId === parseInt(userId)) {
                    return; // Don't close if clicking the same user
                }
            }
            
            if (!modal.contains(event.target)) {
                this.closeChat(userId);
                document.removeEventListener('click', handleClickOutside);
            }
        };

        // Add slight delay to prevent immediate closure
        setTimeout(() => {
            document.addEventListener('click', handleClickOutside);
        }, 100);
    }

    // Focus chat modal
    focusChat(userId) {
        const modal = document.getElementById(`chat-modal-${userId}`);
        if (modal) {
            modal.style.zIndex = ++this.zIndexCounter;
            const input = modal.querySelector(`#chat-input-${userId}`);
            if (input) input.focus();
            
            // Mark as read
            this.model.markAsRead(userId);
        }
    }


    // Close chat modal
    closeChat(userId) {
        const modal = document.getElementById(`chat-modal-${userId}`);
        if (modal) {
            modal.remove();
            this.model.removeChat(userId);
            
            // Clear typing timeout
            const timeout = this.typingTimeouts.get(userId);
            if (timeout) {
                clearTimeout(timeout);
                this.typingTimeouts.delete(userId);
            }
        }
    }

    // Check if a chat with the given user is currently open
    isCurrentChat(userId) {
        return this.model.hasChat(userId) && document.getElementById(`chat-modal-${userId}`) !== null;
    }

    // Send message
    async sendMessage(userId) {
        const input = document.getElementById(`chat-input-${userId}`);
        const message = input.value.trim();
        
        if (!message) return;

        // Clear input immediately
        input.value = '';

        // Stop typing
        this.sendTypingStatus(userId, false);

        // Send via WebSocket - wait for server response to add to UI
        const success = await this.model.sendMessage(userId, message);
        
        if (!success) {
            // Show error if send failed
            console.error('Failed to send message');
            // Could add error UI here if needed
        }
    }

    // Add message to UI
    addMessageToUI(userId, message) {
        const messagesContainer = document.getElementById(`chat-messages-${userId}`);
        if (!messagesContainer) return;

        // Remove loading message if present
        const loading = messagesContainer.querySelector('.chat-loading');
        if (loading) loading.remove();

        const messageHtml = PrivateChatView.generateMessage(message, this.isCurrentUser(message.sender));
        messagesContainer.insertAdjacentHTML('beforeend', messageHtml);
        messagesContainer.scrollTop = messagesContainer.scrollHeight;
    }

    // Load chat history using PaginationManager
    async loadChatHistory(userId) {
        try {
            const messagesContainer = document.getElementById(`chat-messages-${userId}`);
            if (!messagesContainer) return;

            // Set page size for messages
            this.paginationManager.setPageSize('messages', 30);

            // Initialize pagination with render callback
            await this.paginationManager.initializePagination(
                `messages-${userId}`, // Unique type per chat
                `/protected/v1/message-history`,
                `chat-messages-${userId}`,
                (messages, replace = false) => {
                    this.renderMessagesFromPagination(userId, messages, replace);
                },
                { user_id: userId }
            );

            // Setup reverse scroll for loading older messages
            this.setupReverseScroll(userId);

        } catch (error) {
            console.error('Error loading chat history:', error);
            // Fallback to welcome message
            const messagesContainer = document.getElementById(`chat-messages-${userId}`);
            if (messagesContainer) {
                messagesContainer.innerHTML = PrivateChatView.generateWelcomeMessage();
            }
        }
    }

    // Setup reverse scroll for loading older messages
    setupReverseScroll(userId) {
        // Use ScrollManager with reverse scrolling (isReverse: true)
        this.paginationManager.scrollManager.setupInfiniteScroll(
            `messages-reverse-${userId}`,
            `chat-messages-${userId}`,
            async () => {
                await this.loadOlderMessages(userId);
            },
            {
                isReverse: true, // Boolean parameter for reverse scrolling
                useThrottle: true,
                useDebounce: false,
                throttleDelay: 200
            }
        );
    }

    // Load older messages
    async loadOlderMessages(userId) {
        const messagesContainer = document.getElementById(`chat-messages-${userId}`);
        if (!messagesContainer) return;

        // Show loading indicator
        this.showTopLoadingIndicator(userId);

        try {
            // Load more messages using pagination manager
            const moreMessages = await this.paginationManager.loadMoreData(`messages-${userId}`);
            // The renderMessagesFromPagination will handle the UI updates
        } catch (error) {
            console.error('Error loading older messages:', error);
        } finally {
            this.hideTopLoadingIndicator(userId);
        }
    }

    // Show loading indicator at top
    showTopLoadingIndicator(userId) {
        const messagesContainer = document.getElementById(`chat-messages-${userId}`);
        if (!messagesContainer) return;

        const existing = messagesContainer.querySelector('.top-loading-indicator');
        if (existing) return;

        const indicator = document.createElement('div');
        indicator.className = 'top-loading-indicator';
        indicator.innerHTML = `
            <div class="loading-spinner"></div>
            <span class="loading-text">Loading older messages...</span>
        `;
        messagesContainer.insertBefore(indicator, messagesContainer.firstChild);
    }

    // Hide loading indicator at top
    hideTopLoadingIndicator(userId) {
        const messagesContainer = document.getElementById(`chat-messages-${userId}`);
        if (!messagesContainer) return;

        const indicator = messagesContainer.querySelector('.top-loading-indicator');
        if (indicator) indicator.remove();
    }

    // Render messages from pagination manager
    renderMessagesFromPagination(userId, messages, replace = false) {
        const messagesContainer = document.getElementById(`chat-messages-${userId}`);
        if (!messagesContainer) return;

        const wasAtBottom = this.isScrolledToBottom(messagesContainer);

        if (replace) {
            // Initial load: clear and add all messages
            messagesContainer.innerHTML = '';
            this.model.clearMessages(userId);
            
            if (Array.isArray(messages)) {
                messages.forEach(message => {
                    const messageData = {
                        content: message.message,
                        sender: message.sender_id,
                        timestamp: message.created_at,
                        senderName: message.sender
                    };
                    this.model.addMessage(userId, messageData);
                });
            }
            
            this.renderMessages(userId);
            
        } else if (Array.isArray(messages) && messages.length > 0) {
            // Loading older messages: prepend to top
            const oldScrollHeight = messagesContainer.scrollHeight;
            const oldScrollTop = messagesContainer.scrollTop;
            
            // Remove loading indicator
            this.hideTopLoadingIndicator(userId);
            
            // Add older messages to model and UI (reverse order since they're older)
            const reversedMessages = [...messages].reverse();
            reversedMessages.forEach(message => {
                const messageData = {
                    content: message.message,
                    sender: message.sender_id,
                    timestamp: message.created_at,
                    senderName: message.sender
                };
                
                this.model.addMessageToTop(userId, messageData);
                
                const messageHtml = PrivateChatView.generateMessage(
                    messageData, 
                    this.isCurrentUser(message.sender_id)
                );
                messagesContainer.insertAdjacentHTML('afterbegin', messageHtml);
            });
            
            // Maintain scroll position
            setTimeout(() => {
                const newScrollHeight = messagesContainer.scrollHeight;
                messagesContainer.scrollTop = oldScrollTop + (newScrollHeight - oldScrollHeight);
            }, 10);
        }

        // If user was at bottom before, keep them there
        if (wasAtBottom && replace) {
            setTimeout(() => {
                messagesContainer.scrollTop = messagesContainer.scrollHeight;
            }, 10);
        }
    }

    // Check if scrolled to bottom
    isScrolledToBottom(container) {
        return container.scrollTop + container.clientHeight >= container.scrollHeight - 10;
    }

    // Check if sender is current user
    isCurrentUser(senderId) {
        return this.currentUserId && (senderId === this.currentUserId || senderId === 'current_user');
    }

    // Render all messages for a chat
    renderMessages(userId) {
        const messagesContainer = document.getElementById(`chat-messages-${userId}`);
        if (!messagesContainer) return;

        const messages = this.model.getMessages(userId);
        messagesContainer.innerHTML = '';

        messages.forEach(message => {
            this.addMessageToUI(userId, message);
        });

        // Auto scroll to bottom
        messagesContainer.scrollTop = messagesContainer.scrollHeight;
    }

    // Send typing status
    async sendTypingStatus(userId, isTyping) {
        console.log('sendTypingStatus called:', { userId, isTyping });
        
        // Update local state
        this.model.setTypingStatus(userId, isTyping);
        
        // Send to server
        try {
            const result = await this.model.sendTypingStatus(userId, isTyping);
            console.log('Typing status sent:', { userId, isTyping, result });
        } catch (error) {
            console.error('Error sending typing status:', error);
        }
    }

    // Handle incoming message from WebSocket
    receiveMessage(senderId, receiverId, message) {
        // Add to model
        const messageData = {
            content: message.content,
            sender: senderId,
            timestamp: message.timestamp || new Date().toISOString()
        };

        // Determine which user the chat is with (the "other" person)
        const otherUserId = this.isCurrentUser(senderId) ? receiverId : senderId;
        
        if (this.model.hasChat(otherUserId)) {
            // Chat is open with this user
            this.model.addMessage(otherUserId, messageData);
            this.addMessageToUI(otherUserId, messageData);
            
            // Update user list priority
            if (window.userListController) {
                window.userListController.updateMessagePriority(otherUserId, messageData.timestamp);
            }
        } else {
            // No chat open - could auto-open or show notification
            console.log(`New message from user ${senderId}: ${message.content}`);
            
            // Still update user list priority
            if (window.userListController) {
                window.userListController.updateMessagePriority(otherUserId, messageData.timestamp);
            }
        }
    }

    // Handle typing status from WebSocket
    handleTypingStatus(senderId, isTyping) {
        console.log('handleTypingStatus called:', { senderId, isTyping });
        
        // Update typing indicator in chat modal if open
        const messagesContainer = document.getElementById(`chat-messages-${senderId}`);
        console.log('Looking for chat messages container:', `chat-messages-${senderId}`, 'found:', !!messagesContainer);
        
        if (messagesContainer) {
            // Get user data for username
            const chat = this.model.getChat(senderId);
            const username = chat ? chat.user.username : `User ${senderId}`;
            
            console.log('Updating typing indicator in chat modal for:', username, 'isTyping:', isTyping);
            PrivateChatView.updateTypingIndicator(messagesContainer, username, isTyping);
        }

        // Update typing indicator on user list card
        if (window.userListController) {
            console.log('Updating typing indicator in user list for user:', senderId, 'isTyping:', isTyping);
            if (isTyping) {
                window.userListController.showTypingIndicator(senderId);
            } else {
                window.userListController.hideTypingIndicator(senderId);
            }
        }

        // Clear existing timeout
        const existingTimeout = this.typingTimeouts.get(senderId);
        if (existingTimeout) {
            clearTimeout(existingTimeout);
            this.typingTimeouts.delete(senderId);
        }
    }

    // Update user status in open chat modal
    updateUserStatus(userId, isOnline) {
        const modal = document.getElementById(`chat-modal-${userId}`);
        if (modal) {
            const statusDot = modal.querySelector('.status-dot');
            if (statusDot) {
                statusDot.className = `status-dot ${isOnline ? 'online' : 'offline'}`;
            }
        }
    }
}