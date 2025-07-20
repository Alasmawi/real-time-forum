// User Model - Handles all user data operations
import ErrorHandler from '../modules/utils/errors/error-handler.js';

export class UserModel {
    
    // Fetch current user profile
    static async fetchCurrentUser() {
        try {
            const response = await fetch('/protected/v1/user/me', {
                credentials: 'include'
            });
            
            // Use ErrorHandler to process response
            const errorHandled = await ErrorHandler.handleResponse(response);
            if (errorHandled) {
                return null; // Return null if error was handled
            }
            
            return await response.json();
        } catch (error) {
            console.error('UserModel.fetchCurrentUser error:', error);
            return null; // Return null on network error
        }
    }
    
    // Fetch user list
    static async fetchUserList() {
        try {
            const response = await fetch('/protected/v1/user-list', {
                credentials: 'include'
            });
            
            // Use ErrorHandler to process response
            const errorHandled = await ErrorHandler.handleResponse(response);
            if (errorHandled) {
                return []; // Return empty array if error was handled
            }
            
            return await response.json();
        } catch (error) {
            console.error('UserModel.fetchUserList error:', error);
            return []; // Return empty array on network error
        }
    }
    
    // Login user
    static async login(credentials) {
        try {
            const response = await fetch('/v1/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                credentials: 'include',
                body: JSON.stringify(credentials)
            });
            
            // Use ErrorHandler to process response
            const errorHandled = await ErrorHandler.handleResponse(response);
            if (errorHandled) {
                return null; // Return null if error was handled (401, 404, 500, etc.)
            }
            
            // If not handled by ErrorHandler, return the response for controller to process
            // This includes successful responses (200) and validation errors (422)
            return response;
        } catch (error) {
            console.error('UserModel.login error:', error);
            throw error;
        }
    }
    
    // Register user
    static async register(userData) {
        try {
            const response = await fetch('/v1/register', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                credentials: 'include',
                body: JSON.stringify(userData)
            });
            
            // Use ErrorHandler to process response
            const errorHandled = await ErrorHandler.handleResponse(response);
            if (errorHandled) {
                return null; // Return null if error was handled (401, 404, 500, etc.)
            }
            
            // If not handled by ErrorHandler, return the response for controller to process
            // This includes successful responses (200) and validation errors (422)
            return response;
        } catch (error) {
            console.error('UserModel.register error:', error);
            throw error;
        }
    }
    
    // Logout user
    static async logout() {
        try {
            const response = await fetch('/protected/v1/logout', {
                method: 'POST',
                credentials: 'include'
            });
            
            return response; // Let controller handle response
        } catch (error) {
            console.error('UserModel.logout error:', error);
            throw error;
        }
    }
}
