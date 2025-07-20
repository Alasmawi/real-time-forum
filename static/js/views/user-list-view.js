// User List View - Handles only HTML generation
export class UserListView {
    
    // Generate HTML for a single user item
    static generateUserItem(user) {
        const isOnline = user.status === 1;
        
        // Use avatar generator for consistent avatars
        const avatarLetter = window.AvatarGenerator ? 
            window.AvatarGenerator.generateAvatarLetter(user.username) : 
            user.username.charAt(0).toUpperCase();
        
        const avatarColor = window.AvatarGenerator ? 
            window.AvatarGenerator.getColorFromUsername(user.username) : 
            'primary';
            
        return `
            <div class="user-item ${isOnline ? 'online' : 'offline'}" data-user-id="${user.id}">
                <div class="avatar medium ${avatarColor}">
                    ${avatarLetter}
                </div>
                <div class="user-name">${user.username}</div>
                <div class="typing-indicator" id="typing-indicator-${user.id}" style="display: none;">
                    <span class="typing-dots">
                        <span></span>
                        <span></span>
                        <span></span>
                    </span>
                </div>
                <div class="status-dot ${isOnline ? 'online' : 'offline'}"></div>
            </div>
        `;
    }

    // Generate complete user list HTML
    static generateUserListHTML(onlineUsers, offlineUsers, counts = null) {
        // Use server counts if available, otherwise fall back to array lengths
        const onlineCount = counts ? counts.online : onlineUsers.length;
        const offlineCount = counts ? counts.offline : offlineUsers.length;
        
        return `
            <div class="user-list-content">
                <div class="user-section" id="online-users">
                    <div class="user-section-title">Online - ${onlineCount}</div>
                    ${onlineUsers.map(user => this.generateUserItem(user)).join('')}
                </div>
                
                ${offlineUsers.length > 0 || (counts && counts.offline > 0) ? `
                    <div class="user-section" id="offline-users">
                        <div class="user-section-title">Offline - ${offlineCount}</div>
                        ${offlineUsers.map(user => this.generateUserItem(user)).join('')}
                    </div>
                ` : ''}
                
                ${onlineUsers.length === 0 && offlineUsers.length === 0 ? `
                    <div class="user-list-empty">
                        No users found
                    </div>
                ` : ''}
            </div>
        `;
    }
}