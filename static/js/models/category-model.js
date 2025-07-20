// Categories Model - Handles all category data operations
import ErrorHandler from '../modules/utils/errors/error-handler.js';

export class CategoriesModel {
    
    // Fetch all categories
    static async fetchCategories() {
        try {
            const response = await fetch('/guest/v1/categories', {
                credentials: 'include'
            });
            
            // Use ErrorHandler to process response
            const errorHandled = await ErrorHandler.handleResponse(response);
            if (errorHandled) {
                return []; // Return empty array if error was handled
            }
            
            return await response.json();
        } catch (error) {
            console.error('CategoriesModel.fetchCategories error:', error);
            return []; // Return empty array on network error
        }
    }
    
    // Fetch category by ID
    static async fetchCategoryById(categoryId) {
        try {
            const response = await fetch(`/guest/v1/categories/${categoryId}`, {
                credentials: 'include'
            });
            
            // Use ErrorHandler to process response
            const errorHandled = await ErrorHandler.handleResponse(response);
            if (errorHandled) {
                return null; // Return null if error was handled
            }
            
            return await response.json();
        } catch (error) {
            console.error('CategoriesModel.fetchCategoryById error:', error);
            return null; // Return null on network error
        }
    }
}
