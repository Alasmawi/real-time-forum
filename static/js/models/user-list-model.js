// User List Model - Data management for user list with message priority ordering
export class UserListModel {
    constructor() {
        this.users = []; // Ordered array - message priority > alphabetic
        this.statusCache = new Map(); // User status cache (id -> {username, status})
        this.lastMessageTimes = new Map(); // userId -> timestamp for ordering
        this.unreadCounts = new Map(); // userId -> unread message count
        this.totalUsers = 0; // Total users from backend (excluding current user)
    }

    // Add user to the list (from pagination - append to end)
    addUser(user) {
        const userObj = {
            id: user.id,
            username: user.username,
            status: user.status || 0 // 0 = offline, 1 = online
        };
        
        // Store message time if provided
        if (user.last_message_time) {
            this.lastMessageTimes.set(user.id, new Date(user.last_message_time));
        }
        
        this.users.push(userObj);
    }

    // Clear all users (used when receiving new complete list from server)
    clearUsers() {
        this.users = [];
        // Keep message times and status cache for ordering
    }

    // Set total users from API
    setTotalUsers(totalUsers) {
        console.log('UserListModel.setTotalUsers called with:', totalUsers);
        this.totalUsers = totalUsers;
    }

    // Calculate counts from cache
    getCounts() {
        const onlineCount = this.users.filter(user => user.status === 1).length;
        const offlineCount = this.totalUsers - onlineCount;
        
        console.log('Calculated counts:', { 
            onlineCount, 
            offlineCount, 
            totalUsers: this.totalUsers,
            usersInCache: this.users.length 
        });
        
        return {
            online: onlineCount,
            offline: Math.max(0, offlineCount) // Ensure non-negative
        };
    }

    // Get user by ID
    getUserById(userId) {
        // First check paginated users
        const user = this.users.find(user => user.id === userId);
        if (user) {
            return user;
        }
        
        // Fallback to status cache if user not in paginated list
        const cachedUser = this.statusCache.get(userId);
        if (cachedUser) {
            return {
                id: userId,
                username: cachedUser.username,
                status: cachedUser.status
            };
        }
        
        return null;
    }
    
    // Check if user exists in cache
    hasUser(userId) {
        return this.users.some(user => user.id === userId);
    }

    // Update user status cache from WebSocket
    updateUserStatus(userId, username, status) {
        const isOnline = status === 1;

        this.statusCache.set(userId, {
            username: username,
            status: status
        });

        // Update user in loaded list if they exist
        const user = this.users.find(u => u.id === userId);
        if (user) {
            user.status = status;
        }
    }

    // Set user as offline in cache
    setUserOffline(userId) {
        const cachedUser = this.statusCache.get(userId);
        if (cachedUser) {
            cachedUser.status = 0; // offline
            this.statusCache.set(userId, cachedUser);
        }

        // Update user in loaded list if they exist
        const user = this.users.find(u => u.id === userId);
        if (user) {
            user.status = 0; // offline
        }
    }

    // Get user status from cache (for new users being loaded)
    getUserStatusFromCache(userId) {
        const cachedUser = this.statusCache.get(userId);
        return cachedUser ? cachedUser.status : null; // null = no cache data
    }

    // Get users separated by online/offline status (preserving cache order)
    getUsersByStatus() {
        const onlineUsers = this.users.filter(user => user.status === 1);
        const offlineUsers = this.users.filter(user => user.status === 0);
        
        // Users are already ordered correctly in cache (message priority > alphabetic)
        return { onlineUsers, offlineUsers };
    }
    
    // Find correct insertion position for new user: message priority > alphabetic
    findInsertionIndex(newUser) {
        const newUserLastMessage = this.lastMessageTimes.get(newUser.id) || null;
        
        for (let i = 0; i < this.users.length; i++) {
            const existingUser = this.users[i];
            const existingLastMessage = this.lastMessageTimes.get(existingUser.id) || null;
            
            // Compare by message priority first
            if (newUserLastMessage && existingLastMessage) {
                // Both have messages - newer message goes first
                if (newUserLastMessage > existingLastMessage) return i;
            } else if (newUserLastMessage && !existingLastMessage) {
                // New user has messages, existing doesn't - new user goes first
                return i;
            } else if (!newUserLastMessage && !existingLastMessage) {
                // Neither have messages - alphabetic order
                if (newUser.username.localeCompare(existingUser.username) < 0) return i;
            }
            // Continue if existing user has messages but new user doesn't
        }
        
        return this.users.length; // Insert at end
    }
    
    // Insert new user at specific position in cache
    insertUserAtIndex(user, index) {
        const userObj = {
            id: user.id,
            username: user.username,
            status: user.status || 0 // 0 = offline, 1 = online
        };
        
        // Store message time if provided
        if (user.last_message_time) {
            this.lastMessageTimes.set(user.id, new Date(user.last_message_time));
        }
        
        this.users.splice(index, 0, userObj);
    }
    
    // Update last message time and reorder user to top
    updateLastMessageTime(userId, timestamp) {
        this.lastMessageTimes.set(userId, new Date(timestamp));
        
        // Find user and move to appropriate position
        const userIndex = this.users.findIndex(u => u.id === userId);
        if (userIndex !== -1) {
            const user = this.users.splice(userIndex, 1)[0]; // Remove from current position
            const newIndex = this.findInsertionIndex(user); // Find new position
            this.users.splice(newIndex, 0, user); // Insert at new position
        }
    }
}