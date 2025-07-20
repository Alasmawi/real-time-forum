// Notifications Controller - Handles notifications modal interactions
import { NotificationsModel } from '../models/notifications-model.js';
import { NotificationsView } from '../views/notifications-view.js';

export class NotificationsController {
    static model = new NotificationsModel();
    static isModalOpen = false;

    // Toggle notifications modal
    static toggleModal() {
        if (this.isModalOpen) {
            this.closeModal();
        } else {
            this.openModal();
        }
    }

    // Open notifications modal
    static openModal() {
        const notifications = this.model.getAllNotifications();
        const position = this.calculateModalPosition();
        const modalHTML = NotificationsView.generateNotificationsModal(notifications, position);
        
        // Create modal container if it doesn't exist
        let modalContainer = document.getElementById('notifications-modal-container');
        if (!modalContainer) {
            modalContainer = document.createElement('div');
            modalContainer.id = 'notifications-modal-container';
            document.body.appendChild(modalContainer);
        }
        
        modalContainer.innerHTML = modalHTML;
        
        // Set position using CSS custom properties if positioned
        if (position) {
            const content = modalContainer.querySelector('.notifications-content');
            content.style.setProperty('--modal-left', `${position.left}px`);
            content.style.setProperty('--modal-top', `${position.top}px`);
        }
        
        this.isModalOpen = true;
        
        // Setup modal event listeners
        this.setupModalEventListeners();
    }

    // Calculate modal position relative to notification button
    static calculateModalPosition() {
        const notificationsBtn = document.getElementById('notifications-btn');
        if (!notificationsBtn) return null;

        const rect = notificationsBtn.getBoundingClientRect();
        const modalWidth = 350;
        const modalHeight = 400;
        const padding = 10;

        // Check if there's enough space to position next to button
        const viewportWidth = window.innerWidth;
        const viewportHeight = window.innerHeight;

        // Try to position to the right of the button
        let left = rect.right + padding;
        let top = rect.top;

        // If modal would go off right edge, position to the left
        if (left + modalWidth > viewportWidth) {
            left = rect.left - modalWidth - padding;
        }

        // If modal would go off left edge, center it
        if (left < 0) {
            return null; // Fall back to centered modal
        }

        // Adjust vertical position if needed
        if (top + modalHeight > viewportHeight) {
            top = viewportHeight - modalHeight - padding;
        }

        // If still not enough space, fall back to centered
        if (top < 0 || viewportWidth < 800) {
            return null;
        }

        return { left, top };
    }

    // Close notifications modal
    static closeModal() {
        const modalContainer = document.getElementById('notifications-modal-container');
        if (modalContainer) {
            modalContainer.remove();
        }
        this.isModalOpen = false;
    }

    // Setup modal event listeners
    static setupModalEventListeners() {
        // Close button
        const closeBtn = document.getElementById('notifications-close');
        if (closeBtn) {
            closeBtn.addEventListener('click', () => this.closeModal());
        }

        // Click outside modal to close
        const modal = document.getElementById('notifications-modal');
        const content = modal?.querySelector('.notifications-content');
        
        if (modal) {
            modal.addEventListener('click', (e) => {
                // Close if clicking on modal background (not content)
                if (e.target === modal) {
                    this.closeModal();
                }
            });
        }

        // For positioned modals, also handle document clicks
        if (modal?.classList.contains('positioned')) {
            // Add slight delay to prevent immediate closure
            setTimeout(() => {
                const handleClickOutside = (e) => {
                    if (content && !content.contains(e.target)) {
                        this.closeModal();
                        document.removeEventListener('click', handleClickOutside);
                    }
                };
                document.addEventListener('click', handleClickOutside);
            }, 100);
        }

        // Notification card clicks - open private chat
        const notificationCards = document.querySelectorAll('.notification-card');
        notificationCards.forEach(card => {
            card.addEventListener('click', () => {
                const userId = parseInt(card.getAttribute('data-user-id'));
                this.handleNotificationClick(userId);
            });
        });
    }

    // Handle notification card click - open chat and clear notification
    static handleNotificationClick(userId) {
        console.log('Notification clicked for user ID:', userId);
        
        // Clear the notification
        this.model.clearNotification(userId);
        
        // Open private chat if controller exists
        if (window.privateChatController && window.userListController) {
            const user = window.userListController.model.getUserById(userId);
            console.log('Found user for notification click:', user);
            
            if (user) {
                const userElement = document.querySelector(`[data-user-id="${userId}"]`);
                console.log('Found user element:', userElement);
                console.log('Opening chat for user:', user);
                window.privateChatController.openChat(user, userElement);
            } else {
                console.error('User not found in cache for ID:', userId);
            }
        } else {
            console.error('Controllers not available:', {
                privateChatController: !!window.privateChatController,
                userListController: !!window.userListController
            });
        }
        
        // Close modal
        this.closeModal();
    }

    // Add notification from websocket message
    static addNotification(userId, username, messageCount = 1) {
        this.model.addNotification(userId, username, messageCount);
        this.updateNotificationBadge();
    }

    // Clear notification when chat is opened
    static clearNotification(userId) {
        this.model.clearNotification(userId);
        this.updateNotificationBadge();
    }

    // Update notification badge on bell icon
    static updateNotificationBadge() {
        const totalCount = this.model.getTotalNotificationCount();
        const notificationsBtn = document.getElementById('notifications-btn');
        
        if (notificationsBtn) {
            // Remove existing badge
            const existingBadge = notificationsBtn.querySelector('.notification-badge');
            if (existingBadge) {
                existingBadge.remove();
            }
            
            // Add new badge if there are notifications
            if (totalCount > 0) {
                const badge = document.createElement('span');
                badge.className = 'notification-badge';
                badge.textContent = totalCount > 99 ? '99+' : totalCount;
                notificationsBtn.appendChild(badge);
            }
        }
    }

    // Get notification model for external access
    static getModel() {
        return this.model;
    }
}