// Avatar Generator Utility
export class AvatarGenerator {
    
    // Generate avatar letter from username
    static generateAvatarLetter(username) {
        return username.charAt(0).toUpperCase();
    }

    // Generate avatar HTML with optional online indicator
    static generateAvatarHTML(username, options = {}) {
        const {
            size = 'large',           // small, medium, large, extra-large
            color = null,             // Will be auto-generated if not provided
            showOnlineIndicator = false,
            className = ''            // Additional CSS classes
        } = options;
        
        // Auto-generate color if not provided
        const avatarColor = color || this.getColorFromUsername(username);

        const letter = this.generateAvatarLetter(username);
        const classes = `avatar ${size} ${avatarColor} ${className}`.trim();
        const onlineIndicator = showOnlineIndicator ? '<div class="online-indicator"></div>' : '';

        return `
            <div class="${classes}">
                ${letter}
                ${onlineIndicator}
            </div>
        `;
    }

    // Generate avatar for user object
    static generateUserAvatar(user, options = {}) {
        const {
            size = 'large',
            showOnlineIndicator = false,
            className = ''
        } = options;

        // Determine color based on username hash for consistency
        const color = this.getColorFromUsername(user.username);
        const isOnline = user.online || user.status === 'online';

        return this.generateAvatarHTML(user.username, {
            size,
            color,
            showOnlineIndicator: showOnlineIndicator && isOnline,
            className
        });
    }

    // Get consistent color based on username
    static getColorFromUsername(username) {
        const colors = ['primary', 'secondary', 'success', 'info', 'warning', 'danger'];
        let hash = 0;
        
        for (let i = 0; i < username.length; i++) {
            hash = username.charCodeAt(i) + ((hash << 5) - hash);
        }
        
        return colors[Math.abs(hash) % colors.length];
    }

    // Update existing avatar element
    static updateAvatar(element, username, options = {}) {
        if (!element) return;
        
        const avatarHTML = this.generateAvatarHTML(username, options);
        element.outerHTML = avatarHTML;
    }

    // Create avatar element
    static createAvatarElement(username, options = {}) {
        const wrapper = document.createElement('div');
        wrapper.innerHTML = this.generateAvatarHTML(username, options);
        return wrapper.firstElementChild;
    }
}

// Export for global use
window.AvatarGenerator = AvatarGenerator;
