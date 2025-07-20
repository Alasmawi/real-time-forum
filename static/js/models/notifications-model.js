// Notifications Model - Data management for notifications using user list structure
export class NotificationsModel {
    constructor() {
        this.notifications = []; // Array of notification objects
        this.unreadCounts = new Map(); // userId -> unread message count
    }

    // Add/update notification for user
    addNotification(userId, username, messageCount = 1) {
        const existingIndex = this.notifications.findIndex(n => n.userId === userId);
        
        if (existingIndex !== -1) {
            // Update existing notification
            this.notifications[existingIndex].messageCount += messageCount;
        } else {
            // Add new notification
            this.notifications.unshift({
                userId,
                username,
                messageCount
            });
        }
        
        // Update unread count
        const currentCount = this.unreadCounts.get(userId) || 0;
        this.unreadCounts.set(userId, currentCount + messageCount);
    }

    // Clear notification for user (when chat is opened)
    clearNotification(userId) {
        this.notifications = this.notifications.filter(n => n.userId !== userId);
        this.unreadCounts.delete(userId);
    }

    // Get all notifications
    getAllNotifications() {
        return this.notifications;
    }

    // Get unread count for user
    getUnreadCount(userId) {
        return this.unreadCounts.get(userId) || 0;
    }

    // Get total notification count
    getTotalNotificationCount() {
        return this.notifications.reduce((total, n) => total + n.messageCount, 0);
    }

    // Check if user has notifications
    hasNotifications(userId) {
        return this.notifications.some(n => n.userId === userId);
    }

    // Clear all notifications
    clearAllNotifications() {
        this.notifications = [];
        this.unreadCounts.clear();
    }
}