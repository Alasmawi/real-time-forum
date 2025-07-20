// Settings Model - Handles settings data and state management
export class SettingsModel {
    
    // Get current theme setting
    static getCurrentTheme() {
        return document.documentElement.getAttribute('data-theme') || 'dark';
    }
    
    // Check if current theme is dark mode
    static isDarkMode() {
        return this.getCurrentTheme() === 'dark';
    }
    
    // Get appropriate theme icon based on current theme
    static getThemeIcon() {
        return this.isDarkMode() 
            ? 'images/svg-small-light-theme.svg'  // Show sun in dark mode (to switch TO light)
            : 'images/svg-small-dark-theme.svg';  // Show moon in light mode (to switch TO dark)
    }
    
    // Get settings menu visibility state
    static isMenuVisible() {
        const menu = document.getElementById('settings-menu');
        return menu ? menu.classList.contains('expanded') : false;
    }
    
    // Settings menu configuration
    static getMenuConfig() {
        return {
            gearIcon: 'images/svg-small-gear.svg',
            themeIcon: this.getThemeIcon(),
            themeText: 'Switch Theme',
            logoutIcon: 'images/svg-small-logout.svg',
            logoutText: 'Logout'
        };
    }
}