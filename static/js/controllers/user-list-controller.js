// User List Controller - Orchestrates model and view, handles user interactions
import { UserListModel } from '../models/user-list-model.js';
import { UserListView } from '../views/user-list-view.js';
import ErrorHandler from '../modules/utils/errors/error-handler.js';

export class UserListController {
    constructor() {
        this.model = new UserListModel();
        this.userListElement = null;
        this.init();
    }

    // Initialize controller
    init() {
        this.userListElement = document.getElementById('user-list');
        // Don't render immediately - wait for data to be loaded
    }

    // Add user to the model (for pagination - no re-render)
    addUser(user) {
        // Apply cached status before adding user
        user = this.applyStatusFromCache(user);
        this.model.addUser(user);
        
        // Append to DOM without full re-render
        this.appendUserToDOM(user);
        console.log('Added user via pagination:', user.username);
    }
    
    // Append single user to existing DOM structure
    appendUserToDOM(user) {
        const isOnline = user.status === 1;
        const targetSection = isOnline ? '#online-users' : '#offline-users';
        let section = this.userListElement.querySelector(targetSection);
        
        // Create section if it doesn't exist
        if (!section) {
            section = this.createUserSection(isOnline);
        }
        
        // Create user element
        const userElement = this.createUserElement(user, isOnline);
        
        // Append to end of section
        section.appendChild(userElement);
        
        // Attach event listener
        this.attachEventListenerToUser(userElement);
    }

    // Update user list from server payload (replaces entire list)
    updateFromServerPayload(payload) {
        console.log('UserListController.updateFromServerPayload called with:', payload);
        this.model.clearUsers();
        
        // Handle new payload structure with metadata
        if (payload.users) {
            // New structure: { users: [...], online_count: X, offline_count: Y }
            console.log('Processing', payload.users.length, 'users');
            console.log('First few users:', payload.users.slice(0, 3));
            console.log('Users with online=true:', payload.users.filter(u => u.online));
            payload.users.forEach(user => {
                // Apply cached status before adding user
                user = this.applyStatusFromCache(user);
                this.model.addUser(user);
            });
            console.log('Setting total users from pagination:', payload.pagination.total_count);
            this.model.setTotalUsers(payload.pagination.total_count);
        } else {
            // Legacy structure: just array of users
            console.log('Processing legacy payload with', payload.length, 'users');
            payload.forEach(user => {
                // Apply cached status before adding user
                user = this.applyStatusFromCache(user);
                this.model.addUser(user);
            });
        }
        
        console.log('About to call render() from updateFromServerPayload');
        this.render();
    }

    // Render the user list using model data and view
    render() {
        console.log('UserListController.render() called');
        if (!this.userListElement) {
            console.log('No userListElement found in render()');
            return;
        }

        const { onlineUsers, offlineUsers } = this.model.getUsersByStatus();
        const counts = this.model.getCounts();
        console.log('Render data:', { 
            onlineUsers: onlineUsers.length, 
            offlineUsers: offlineUsers.length, 
            counts: counts,
            'counts.online': counts.online,
            'counts.offline': counts.offline
        });
        
        const html = UserListView.generateUserListHTML(onlineUsers, offlineUsers, counts);
        console.log('Generated HTML length:', html.length);
        this.userListElement.innerHTML = html;
        console.log('HTML inserted into DOM');
        
        // Re-attach event listeners after rendering
        this.attachUserClickListeners();
    }

    // Attach click event listeners to user items
    attachUserClickListeners() {
        // Remove existing listeners first to avoid duplicates
        const existingItems = this.userListElement.querySelectorAll('.user-item[data-listeners-attached="true"]');
        existingItems.forEach(item => {
            item.removeAttribute('data-listeners-attached');
        });
        
        // Attach listeners to all user items (both online and offline)
        const userItems = this.userListElement.querySelectorAll('.user-item:not([data-listeners-attached])');
        userItems.forEach(item => {
            item.addEventListener('click', () => {
                const userId = parseInt(item.getAttribute('data-user-id'));
                
                // Try to get user from model first (paginated users)
                let user = this.model.getUserById(userId);
                
                // If not in model, try to get from cache (injected users)
                if (!user) {
                    const cachedUser = this.model.statusCache.get(userId);
                    if (cachedUser) {
                        user = {
                            id: userId,
                            username: cachedUser.username,
                            online: cachedUser.status === 1
                        };
                    }
                }
                
                // If still no user, get from DOM
                if (!user) {
                    const usernameElement = item.querySelector('.user-name');
                    const isOnline = item.classList.contains('online');
                    if (usernameElement) {
                        user = {
                            id: userId,
                            username: usernameElement.textContent,
                            online: isOnline
                        };
                    }
                }
                
                if (user) {
                    this.handleUserClick(user, item);
                }
            });
            // Mark as having listeners attached
            item.setAttribute('data-listeners-attached', 'true');
        });
    }

    // Handle user click - open private chat
    handleUserClick(user, userElement) {
        if (window.privateChatController && userElement) {
            // Create user object with current status for private chat
            const userForChat = {
                id: user.id,
                username: user.username,
                online: user.status === 1
            };
            window.privateChatController.openChat(userForChat, userElement);
        }
    }

    // Show typing indicator for a user
    showTypingIndicator(userId) {
        const indicator = document.getElementById(`typing-indicator-${userId}`);
        
        if (indicator) {
            indicator.style.display = 'flex';
        }
    }

    // Hide typing indicator for a user
    hideTypingIndicator(userId) {
        const indicator = document.getElementById(`typing-indicator-${userId}`);
        
        if (indicator) {
            indicator.style.display = 'none';
        }
    }

    // Handle WebSocket user status updates
    async handleUserStatusUpdate(statusUpdate) {
        console.log('UserListController.handleUserStatusUpdate called with:', statusUpdate);
        
        // Handle online users
        if (statusUpdate.online_users) {
            const onlineUserPromises = statusUpdate.online_users.map(async user => {
                if (!this.model.hasUser(user.id)) {
                    // New user - fetch message priority data before adding
                    console.log(`Fetching message priority for new user ${user.id} (${user.username})`);
                    try {
                        const messagePriorityData = await this.fetchUserMessagePriority(user.id);
                        const userWithPriority = {
                            ...user,
                            last_message_time: messagePriorityData.last_message_time
                        };
                        
                        console.log(`Adding new user ${user.id} with message priority to cache`);
                        console.log(`User object from WebSocket:`, user);
                        console.log(`User with priority:`, userWithPriority);
                        const insertIndex = this.model.findInsertionIndex(userWithPriority);
                        this.model.insertUserAtIndex(userWithPriority, insertIndex);
                        this.insertUserIntoDOMAtIndex(userWithPriority, insertIndex, true);
                        
                        // Debug: Check what's in cache after adding user
                        console.log(`Cache after adding user ${user.id}:`, this.model.users.map(u => ({id: u.id, username: u.username, status: u.status})));
                        const counts = this.model.getCounts();
                        console.log(`Counts after adding user ${user.id}:`, counts);
                    } catch (error) {
                        console.error(`Failed to fetch message priority for user ${user.id}:`, error);
                        // Fallback: add user without message priority data
                        const insertIndex = this.model.findInsertionIndex(user);
                        this.model.insertUserAtIndex(user, insertIndex);
                        this.insertUserIntoDOMAtIndex(user, insertIndex, true);
                    }
                } else {
                    // Existing user - update status
                    console.log(`Updating existing user ${user.id} (${user.username}) to online`);
                    this.model.updateUserStatus(user.id, user.username, user.status);
                    this.updateUserStatusInDOM(user.id, true);
                    
                    // Update chat modal status if open
                    if (window.privateChatController) {
                        window.privateChatController.updateUserStatus(user.id, true);
                    }
                }
            });
            
            // Wait for all async operations to complete
            await Promise.all(onlineUserPromises);
        }
        
        // Handle offline user IDs (only update if user exists in cache)
        if (statusUpdate.offline_user_ids) {
            statusUpdate.offline_user_ids.forEach(userId => {
                if (this.model.hasUser(userId)) {
                    // Existing user - update status to offline
                    console.log(`Updating existing user ${userId} to offline`);
                    this.model.setUserOffline(userId);
                    this.updateUserStatusInDOM(userId, false);
                    
                    // Update chat modal status if open
                    if (window.privateChatController) {
                        window.privateChatController.updateUserStatus(userId, false);
                    }
                } else {
                    // User not in cache - just update status cache for future reference
                    console.log(`User ${userId} went offline but not in cache - updating status cache only`);
                    this.model.setUserOffline(userId);
                }
            });
        }
        
        // Update section counts in DOM (calculated from cache)
        this.updateSectionCounts();
        
        // Re-render to ensure users are in correct sections (online/offline)
        // This is needed for initial status updates where users were loaded as offline
        // but need to be moved to online section
        console.log('Re-rendering after status update to ensure correct user placement');
        this.render();
    }


    // Apply cached status when loading new users via pagination
    applyStatusFromCache(user) {
        const cachedStatus = this.model.getUserStatusFromCache(user.id);
        if (cachedStatus !== null && cachedStatus !== undefined) {
            // Use cached status - WebSocket data takes precedence
            user.status = cachedStatus; // Use status directly
        } else {
            // No cached status, keep API response (defaults to offline)
            user.status = user.status || 0; // 0 = offline
        }
        return user;
    }

    // Insert new user into DOM at specific index
    insertUserIntoDOMAtIndex(user, index, isOnline) {
        const userElement = this.createUserElement(user, isOnline);
        const { onlineUsers, offlineUsers } = this.model.getUsersByStatus();
        
        // Find the target section
        const targetSection = isOnline ? '#online-users' : '#offline-users';
        let section = this.userListElement.querySelector(targetSection);
        
        // Create section if it doesn't exist
        if (!section) {
            section = this.createUserSection(isOnline);
        }
        
        // Find correct DOM position within the section
        const sectionUsers = isOnline ? onlineUsers : offlineUsers;
        const userIndexInSection = sectionUsers.findIndex(u => u.id === user.id);
        
        // Find the DOM element after which to insert
        const existingUserElements = section.querySelectorAll('.user-item');
        if (userIndexInSection === 0 || existingUserElements.length === 0) {
            // Insert after section title
            const sectionTitle = section.querySelector('.user-section-title');
            sectionTitle.insertAdjacentElement('afterend', userElement);
        } else if (userIndexInSection < existingUserElements.length) {
            // Insert before the element at userIndexInSection
            existingUserElements[userIndexInSection].insertAdjacentElement('beforebegin', userElement);
        } else {
            // Insert at end of section
            section.appendChild(userElement);
        }
        
        // Attach event listener
        this.attachEventListenerToUser(userElement);
    }
    
    // Update user status in DOM (move between sections)
    updateUserStatusInDOM(userId, isOnline) {
        const userElement = this.userListElement.querySelector(`[data-user-id="${userId}"]`);
        if (!userElement) return;
        
        const user = this.model.getUserById(userId);
        if (!user) return;
        
        // Remove from current position
        userElement.remove();
        
        // Find new position in cache
        const { onlineUsers, offlineUsers } = this.model.getUsersByStatus();
        const targetUsers = isOnline ? onlineUsers : offlineUsers;
        const userIndex = targetUsers.findIndex(u => u.id === userId);
        
        // Insert at new position
        this.insertUserIntoDOMAtIndex(user, userIndex, isOnline);
    }
    
    // Create user element HTML
    createUserElement(user, isOnline) {
        const userElement = document.createElement('div');
        userElement.className = `user-item ${isOnline ? 'online' : 'offline'}`;
        userElement.setAttribute('data-user-id', user.id);
        
        // Use avatar generator for consistent avatars
        const avatarLetter = window.AvatarGenerator ? 
            window.AvatarGenerator.generateAvatarLetter(user.username) : 
            user.username.charAt(0).toUpperCase();
        
        const avatarColor = window.AvatarGenerator ? 
            window.AvatarGenerator.getColorFromUsername(user.username) : 
            'primary';
            
        userElement.innerHTML = `
            <div class="avatar medium ${avatarColor}">
                ${avatarLetter}
            </div>
            <div class="user-info">
                <div class="user-name">${user.username}</div>
                <div class="user-status ${isOnline ? 'online' : 'offline'}">
                    <div class="status-dot ${isOnline ? 'online' : 'offline'}"></div>
                    <div class="typing-indicator" id="typing-indicator-${user.id}" style="display: none;">
                        <span class="typing-dots">
                            <span></span>
                            <span></span>
                            <span></span>
                        </span>
                    </div>
                </div>
            </div>
        `;
        
        return userElement;
    }
    
    // Create user section if it doesn't exist
    createUserSection(isOnline) {
        const userListContent = this.userListElement.querySelector('.user-list-content');
        if (!userListContent) return null;
        
        const sectionId = isOnline ? 'online-users' : 'offline-users';
        const sectionTitle = isOnline ? 'Online' : 'Offline';
        const counts = this.model.getCounts();
        const count = isOnline ? counts.online : counts.offline;
        
        const sectionElement = document.createElement('div');
        sectionElement.className = 'user-section';
        sectionElement.id = sectionId;
        sectionElement.innerHTML = `<div class="user-section-title">${sectionTitle} - ${count}</div>`;
        
        if (isOnline) {
            // Insert online section at the beginning
            userListContent.insertBefore(sectionElement, userListContent.firstChild);
        } else {
            // Append offline section at the end
            userListContent.appendChild(sectionElement);
        }
        
        return sectionElement;
    }
    
    // Update section counts in DOM
    updateSectionCounts() {
        const counts = this.model.getCounts();
        
        const onlineTitle = this.userListElement.querySelector('#online-users .user-section-title');
        if (onlineTitle) {
            onlineTitle.textContent = `Online - ${counts.online}`;
        }
        
        const offlineTitle = this.userListElement.querySelector('#offline-users .user-section-title');
        if (offlineTitle) {
            offlineTitle.textContent = `Offline - ${counts.offline}`;
        }
    }
    
    // Attach event listener to single user element
    attachEventListenerToUser(userElement) {
        userElement.addEventListener('click', () => {
            const userId = parseInt(userElement.getAttribute('data-user-id'));
            const user = this.model.getUserById(userId);
            
            if (!user) {
                // Try to get from cache
                const cachedUser = this.model.statusCache.get(userId);
                if (cachedUser) {
                    const userForChat = {
                        id: userId,
                        username: cachedUser.username,
                        online: cachedUser.status === 1
                    };
                    this.handleUserClick(userForChat, userElement);
                }
            } else {
                this.handleUserClick(user, userElement);
            }
        });
        
        // Mark as having listeners attached
        userElement.setAttribute('data-listeners-attached', 'true');
    }
    
    // Update message priority when user is messaged (moves to top)
    updateMessagePriority(userId, timestamp) {
        console.log(`Updating message priority for user ${userId}`);
        
        // Update model with new message time
        this.model.updateLastMessageTime(userId, timestamp);
        
        // Find user element in DOM
        const userElement = this.userListElement.querySelector(`[data-user-id="${userId}"]`);
        if (!userElement) return;
        
        const user = this.model.getUserById(userId);
        if (!user) return;
        
        // Remove from current position
        userElement.remove();
        
        // Find new position (should be at top of their section)
        const { onlineUsers, offlineUsers } = this.model.getUsersByStatus();
        const targetUsers = user.status === 1 ? onlineUsers : offlineUsers;
        const userIndex = targetUsers.findIndex(u => u.id === userId);
        
        // Insert at new position
        this.insertUserIntoDOMAtIndex(user, userIndex, user.status === 1);
    }
    
    // Fetch user message priority data from API
    async fetchUserMessagePriority(userId) {
        const response = await fetch(`/protected/v1/user/${userId}/message-priority`, {
            method: 'GET',
            credentials: 'include'
        });
        
        const errorHandled = await ErrorHandler.handleResponse(response);
        if (errorHandled) {
            throw new Error(`Failed to fetch message priority: ${response.status}`);
        }
        
        return await response.json();
    }

}