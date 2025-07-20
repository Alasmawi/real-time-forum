// Comments Model - Handles all comment data operations
import ErrorHandler from '../modules/utils/errors/error-handler.js';

export class CommentsModel {
    
    // Fetch comments for a specific post with pagination
    static async fetchCommentsByPostId(postId, params = {}) {
        try {
            const url = new URL('/guest/v1/comments', window.location.origin);
            
            // Add post_id as required parameter
            url.searchParams.append('post_id', postId);
            
            // Add optional pagination parameters
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
                return { comments: [], pagination: null }; // Return empty structure if error was handled
            }
            
            return await response.json();
        } catch (error) {
            console.error('CommentsModel.fetchCommentsByPostId error:', error);
            return { comments: [], pagination: null }; // Return empty structure on network error
        }
    }
    
    // Create new comment
    static async createComment(commentData) {
        try {
            const response = await fetch('/protected/v1/newcomment', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                credentials: 'include',
                body: JSON.stringify(commentData)
            });
            
            // Use ErrorHandler to process response
            const errorHandled = await ErrorHandler.handleResponse(response);
            if (errorHandled) {
                return null; // Return null if error was handled
            }
            
            return response; // Return response object for further processing
        } catch (error) {
            console.error('CommentsModel.createComment error:', error);
            return null; // Return null on network error
        }
    }
}
