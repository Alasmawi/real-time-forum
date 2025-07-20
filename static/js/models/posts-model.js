// Posts Model - Handles all posts data operations
import ErrorHandler from '../modules/utils/errors/error-handler.js';

export class PostsModel {
    
    // Fetch posts from API
    static async fetchPosts(params = {}) {
        try {
            const url = new URL('/guest/v1/posts', window.location.origin);
            
            // Add query parameters if provided
            Object.keys(params).forEach(key => {
                if (params[key]) {
                    url.searchParams.append(key, params[key]);
                }
            });
            
            const response = await fetch(url, {
                credentials: 'include'
            });
            
            // Use ErrorHandler to process response
            const errorHandled = await ErrorHandler.handleResponse(response);
            if (errorHandled) {
                return { posts: [], pagination: null }; // Return empty structure if error was handled
            }
            
            const data = await response.json();
            return data; // Return the full paginated response structure
        } catch (error) {
            console.error('PostsModel.fetchPosts error:', error);
            return { posts: [], pagination: null }; // Return empty structure on network error
        }
    }
    
    // Fetch single post by ID
    static async fetchPostById(postId) {
        try {
            const response = await fetch(`/guest/v1/posts/${postId}`, {
                credentials: 'include'
            });
            
            // Use ErrorHandler to process response
            const errorHandled = await ErrorHandler.handleResponse(response);
            if (errorHandled) {
                return null; // Return null if error was handled
            }
            
            return await response.json();
        } catch (error) {
            console.error('PostsModel.fetchPostById error:', error);
            return null; // Return null on network error
        }
    }
    
    // Create new post
    static async createPost(postData) {
        try {
            const response = await fetch('/protected/v1/newpost', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                credentials: 'include',
                body: JSON.stringify(postData)
            });
            
            // Use ErrorHandler to process response
            const errorHandled = await ErrorHandler.handleResponse(response);
            if (errorHandled) {
                return null; // Return null if error was handled (401, 404, 500, etc.)
            }
            
            // If not handled by ErrorHandler, return the response for controller to process
            // This includes successful responses (201) and validation errors (422)
            return response;
        } catch (error) {
            console.error('PostsModel.createPost error:', error);
            throw error;
        }
    }
}
