// App Initialization Utilities
import { UserModel } from '../../models/user-model.js'
import { SettingsController } from '../../controllers/settings-controller.js'
import { SettingsMenuView } from '../../views/settings-menu.js'
import { AvatarGenerator } from './UI-UX/avatar-generator.js'
import { initializePrivateMessaging } from '../websockets/private-messaging.js'
import { UserListController } from '../../controllers/user-list-controller.js'
import { PrivateChatController } from '../../controllers/private-chat-controller.js'
import { NotificationsController } from '../../controllers/notifications-controller.js'

export class AppInit {
  // Initialize navigation for authenticated views
  static initializeNavigation () {
    const navCard = document.getElementById('nav-card')
    if (navCard) {
      navCard.innerHTML = `
                <div class="nav-links">
                    <a href="/home" class="nav-link" data-link>Home</a>
                </div>
            `
      navCard.style.display = 'block'

      // Update active navigation link
      this.updateActiveNavigation()
    }
  }

  // Update active navigation link based on current path
  static updateActiveNavigation () {
    const currentPath = window.location.pathname
    const navLinks = document.querySelectorAll('.nav-link')

    navLinks.forEach(link => {
      link.classList.remove('active')
      if (link.getAttribute('href') === currentPath) {
        link.classList.add('active')
      }
    })
  }

  // Initialize user list with pagination
  static async initializeUserList () {
    const userListElement = document.getElementById('user-list')
    if (userListElement) {
      userListElement.style.display = 'block'

      // Initialize user list controller
      if (!window.userListController) {
        window.userListController = new UserListController();
      }

      // Initialize private chat controller
      if (!window.privateChatController) {
        window.privateChatController = new PrivateChatController();
        // Make it available globally for WebSocket events
        window.privateChatManager = window.privateChatController;
      }

      // Initialize pagination manager if not already done
      if (!window.paginationManager) {
        const PaginationManager = (await import('../pagination/pagination-manager.js')).default;
        window.paginationManager = new PaginationManager();
      }

      // Set up pagination for users
      try {
        await window.paginationManager.initializePagination(
          'users',
          '/protected/v1/user-list',
          'user-list',
          (responseData, replace = false) => {
            if (replace) {
              // Initial load or refresh - update controller with full response including counts
              if (responseData && responseData.items && responseData.items.length > 0) {
                // Create the expected payload structure for the controller
                const payload = {
                  users: responseData.items,
                  pagination: responseData.pagination || { total_count: 0 }
                };
                window.userListController.updateFromServerPayload(payload);
              }
            } else {
              // Infinite scroll - append users (keeping existing counts)
              if (responseData && responseData.items) {
                responseData.items.forEach(user => window.userListController.addUser(user));
              }
            }
          }
        );
      } catch (error) {
        // Fallback to simple fetch
        const users = await UserModel.fetchUserList();
        if (users && users.users && users.users.length > 0) {
          window.userListController.updateFromServerPayload(users);
        }
      }
    }
  }

  // Initialize WebSocket connection for private messaging
  static async initializeWebSocket () {
    try {
      // Initialize private messaging WebSocket
      initializePrivateMessaging()
    } catch (error) {
      console.error('Error initializing WebSocket:', error)
    }
  }

  // Initialize user profile card
  static async initializeUserProfile (user) {
    const userProfileCard = document.getElementById('user-profile-card')
    if (userProfileCard && user) {
      const avatarHTML = AvatarGenerator.generateUserAvatar(user, {
        size: 'large',
        showOnlineIndicator: true
      })

      userProfileCard.innerHTML = `
                <div class="user-profile">
                    <div class="user-profile-header">
                        ${avatarHTML}
                        <div class="user-info">
                            <h4 class="user-name">${user.username}</h4>
                            <p class="user-status">Offline</p>
                        </div>
                        ${SettingsMenuView.generateSettingsButton()}
                    </div>
                </div>
            `
      userProfileCard.style.display = 'block'

      // Insert settings menu as separate card after user profile
      const leftColumn = document.querySelector('.left-column')
      let settingsCard = document.querySelector('.settings-card')
      
      if (!settingsCard) {
        settingsCard = document.createElement('div')
        settingsCard.className = 'settings-card'
        settingsCard.innerHTML = SettingsMenuView.generateSettingsMenu()
        
        // Insert between user profile and nav card
        const navCard = document.getElementById('nav-card')
        leftColumn.insertBefore(settingsCard, navCard)
      }

      // Setup settings menu functionality using SettingsController (always needed for new button)
      SettingsController.initializeSettingsMenu()
      
      // Initialize notifications controller
      window.NotificationsController = NotificationsController;
      
      // Set initial theme icon
      SettingsMenuView.updateThemeIcon()
    }
  }

  // Initialize complete authenticated app state
  static async initializeAuthenticatedApp () {
    // Get current user data for profile using UserModel
    const user = await UserModel.fetchCurrentUser()
    if (user) {
      AppInit.initializeUserProfile(user)
    }

    AppInit.initializeNavigation()
    await AppInit.initializeUserList()
    await AppInit.initializeWebSocket()
  }
}
