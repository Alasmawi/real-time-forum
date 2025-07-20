// Settings Controller - Handles settings menu interactions
import { SettingsMenuView } from '../views/settings-menu.js';
import { LogoutController } from './logout-controller.js';
import { ThemeToggle } from '../modules/utils/UI-UX/theme-toggle.js';
import { NotificationsController } from './notifications-controller.js';

export class SettingsController {
    
    // Initialize settings menu functionality
    static initializeSettingsMenu() {
        this.setupSettingsButton();
        this.setupNotificationsButton();
        this.setupThemeToggle();
        this.setupOutsideClickClose();
        
        // Logout button is handled by LogoutController
        LogoutController.initializeLogoutButton();
    }
    
    // Setup settings button (gear icon) click handler
    static setupSettingsButton() {
        const settingsBtn = document.getElementById('settings-btn');
        if (settingsBtn) {
            settingsBtn.addEventListener('click', (e) => {
                e.stopPropagation();
                SettingsMenuView.toggleMenu();
            });
        }
    }
    
    // Setup notifications button (bell icon) click handler
    static setupNotificationsButton() {
        const notificationsBtn = document.getElementById('notifications-btn');
        if (notificationsBtn) {
            notificationsBtn.addEventListener('click', (e) => {
                e.stopPropagation();
                NotificationsController.toggleModal();
            });
        }
    }
    
    // Setup theme toggle functionality
    static async setupThemeToggle() {
        const themeToggleBtn = document.getElementById('theme-toggle-btn');
        if (themeToggleBtn) {
            try {
                // Add click handler
                themeToggleBtn.addEventListener('click', () => {
                    ThemeToggle.toggle();
                    SettingsMenuView.updateThemeIcon();
                });
            } catch (error) {
                console.error('Error loading theme toggle:', error);
            }
        }
    }
    
    // Setup outside click to close menu
    static setupOutsideClickClose() {
        document.addEventListener('click', (e) => {
            const settingsBtn = document.getElementById('settings-btn');
            const settingsMenu = document.getElementById('settings-menu');
            
            // Close menu if clicking outside both button and menu
            if (settingsBtn && settingsMenu && 
                !settingsBtn.contains(e.target) && 
                !settingsMenu.contains(e.target)) {
                SettingsMenuView.hideMenu();
            }
        });
    }
    
    // Handle settings menu toggle (called by settings button)
    static toggleMenu() {
        SettingsMenuView.toggleMenu();
    }
    
    // Update theme icon after theme change
    static updateThemeIcon() {
        SettingsMenuView.updateThemeIcon();
    }
    
    // Close settings menu
    static closeMenu() {
        SettingsMenuView.hideMenu();
    }
}