// Settings Menu View - Handles settings menu HTML generation
import { SettingsModel } from '../models/settings-model.js';

export class SettingsMenuView {
    
    // Generate settings button HTML
    static generateSettingsButton() {
        const config = SettingsModel.getMenuConfig();
        return `
            <div class="controls">
                <button id="settings-btn" class="settings-btn">
                    <img src="${config.gearIcon}" alt="Settings" class="settings-icon">
                </button>
                <button id="notifications-btn" class="settings-btn">
                    <img src="images/svg-small-bell.svg" alt="Notifications" class="settings-icon">
                </button>
            </div>
        `;
    }
    
    // Generate complete settings menu HTML
    static generateSettingsMenu() {
        const config = SettingsModel.getMenuConfig();
        return `
            <div id="settings-menu" class="settings-menu">
                <button id="theme-toggle-btn" class="menu-item theme-btn">
                    <img src="${config.themeIcon}" alt="Theme" class="theme-icon">
                    <span>${config.themeText}</span>
                </button>
                <div class="menu-separator"></div>
                <button id="logout-btn" class="menu-item logout-btn">
                    <img src="${config.logoutIcon}" alt="Logout" class="logout-icon-svg">
                    <span>${config.logoutText}</span>
                </button>
            </div>
        `;
    }
    
    // Update theme icon in existing menu
    static updateThemeIcon() {
        const themeIcon = document.querySelector('.theme-icon');
        if (themeIcon) {
            themeIcon.src = SettingsModel.getThemeIcon();
        }
    }
    
    // Show settings menu
    static showMenu() {
        const menu = document.getElementById('settings-menu');
        if (menu) {
            menu.classList.add('expanded');
        }
    }
    
    // Hide settings menu
    static hideMenu() {
        const menu = document.getElementById('settings-menu');
        if (menu) {
            menu.classList.remove('expanded');
        }
    }
    
    // Toggle settings menu visibility
    static toggleMenu() {
        if (SettingsModel.isMenuVisible()) {
            this.hideMenu();
        } else {
            this.showMenu();
        }
    }
}