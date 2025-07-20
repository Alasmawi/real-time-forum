// Notifications View - Handles notifications modal HTML generation
export class NotificationsView {
    
    // Generate notifications modal HTML
    static generateNotificationsModal(notifications = [], position = null) {
        const positionClass = position ? 'positioned' : 'centered';
            
        return `
            <div id="notifications-modal" class="notifications-modal ${positionClass}">
                <div class="notifications-content">
                    <div class="notifications-header">
                        <span class="notifications-title">Notifications</span>
                        <button id="notifications-close" class="notifications-close">×</button>
                    </div>
                    <div class="notifications-list">
                        ${notifications.length > 0 ? 
                            notifications.map(notification => this.generateNotificationCard(notification)).join('') :
                            '<div class="no-notifications">No new notifications</div>'
                        }
                    </div>
                </div>
            </div>
        `;
    }

    // Generate single notification card
    static generateNotificationCard(notification) {
        const avatarLetter = window.AvatarGenerator ? 
            window.AvatarGenerator.generateAvatarLetter(notification.username) : 
            notification.username.charAt(0).toUpperCase();
        
        const avatarColor = window.AvatarGenerator ? 
            window.AvatarGenerator.getColorFromUsername(notification.username) : 
            'primary';

        return `
            <div class="notification-card" data-user-id="${notification.userId}">
                <div class="avatar small ${avatarColor}">
                    ${avatarLetter}
                </div>
                <div class="notification-content">
                    <span class="notification-text">
                        <strong>${notification.username}</strong> ${notification.messageCount > 1 ? 
                            `sent ${notification.messageCount} messages` : 
                            'messaged you!'
                        }
                    </span>
                </div>
            </div>
        `;
    }
}