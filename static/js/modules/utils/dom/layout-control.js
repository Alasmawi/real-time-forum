// DOM Control Utilities - Handle layout visibility for different page types
import { DOMUtilities } from './dom-utilities.js';

export class DOMControl {
    
    /**
     * Hide sidebars for login/register pages to center the form
     */
    static hideSidebars() {
        const leftColumn = document.querySelector('.left-column');
        const userList = document.getElementById('user-list');
        
        DOMUtilities.toggleVisibility(leftColumn, false);
        DOMUtilities.toggleVisibility(userList, false);
    }

    /**
     * Show sidebars for authenticated pages
     */
    static showSidebars() {
        const leftColumn = document.querySelector('.left-column');
        const userList = document.getElementById('user-list');
        
        if (leftColumn) {
            leftColumn.style.display = 'flex';
        }
        DOMUtilities.toggleVisibility(userList, true);
    }

    /**
     * Set page layout mode (auth or guest)
     */
    static setLayoutMode(mode) {
        switch (mode) {
            case 'guest':
                this.hideSidebars();
                break;
            case 'authenticated':
                this.showSidebars();
                break;
            default:
                console.warn('Unknown layout mode:', mode);
        }
    }
}
