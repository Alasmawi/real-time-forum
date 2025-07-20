// Reply Model - Handles all reply data operations
import ErrorHandler from '../modules/utils/errors/error-handler.js';

export class ReplyModel {
    
    // Create new reply (comment)
    static async createReply(replyData) {
        try {
            const response = await fetch('/protected/v1/newcomment', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                credentials: 'include',
                body: JSON.stringify(replyData)
            });
            
            // Use ErrorHandler to process response
            const errorHandled = await ErrorHandler.handleResponse(response);
            if (errorHandled) {
                return null; // Return null if error was handled
            }
            
            return response; // Return response object for further processing
        } catch (error) {
            console.error('ReplyModel.createReply error:', error);
            return null; // Return null on network error
        }
    }
    
    // Validate reply data before submission
    static validateReply(content) {
        const errors = [];
        
        if (!content || content.trim().length === 0) {
            errors.push('Reply content cannot be empty');
        }
        
        if (content.length > 500) {
            errors.push('Reply cannot exceed 500 characters');
        }
        
        return {
            isValid: errors.length === 0,
            errors: errors
        };
    }
}
