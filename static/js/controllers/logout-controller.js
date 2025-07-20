// Logout Controller - Handles logout functionality
import { UserModel } from "../models/user-model.js";
import ErrorHandler from "../modules/utils/errors/error-handler.js";

export class LogoutController {
    
    // Handle logout using UserModel
    static async handleLogout() {
        try {
            // Use UserModel for logout
            const response = await UserModel.logout();
            
            // Let ErrorHandler process the response
            const errorHandled = await ErrorHandler.handleResponse(response);
            if (errorHandled) {
                return; // ErrorHandler took care of it
            }
            
            window.history.pushState(null, null, '/');
            window.dispatchEvent(new PopStateEvent('popstate'));
            
        } catch (error) {
            console.error('Logout error:', error);
        }
    }

    // Initialize logout button functionality
    static initializeLogoutButton() {
        const logoutBtn = document.getElementById('logout-btn');
        if (logoutBtn) {
            logoutBtn.addEventListener('click', this.handleLogout.bind(this));
        }
    }
}
