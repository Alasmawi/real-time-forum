// Reply Controller - Handles reply-related operations following MVC pattern
import { ReplyModel } from '../models/reply-model.js';
import { ReplyBoxView } from '../views/posts/reply-box.js';
import { CharacterCounter } from '../modules/utils/character-counter.js';
import { FormHandler } from '../modules/utils/form-handler.js';
import { ReplyBoxManager } from '../modules/utils/dom/reply-box-manager.js';
import { DOMUtilities } from '../modules/utils/dom/dom-utilities.js';
import { FormErrorHelpers } from '../modules/utils/errors/form-error-helpers.js';

export class RepliesController {
    
    // Show reply box for a specific post
    static async showReplyBox(postId) {
        ReplyBoxManager.show(postId);
        // Set up event listeners for the newly shown reply form
        this.setupReplyEventListeners();
    }
    
    // Hide reply box for a specific post
    static hideReplyBox(postId) {
        ReplyBoxManager.hide(postId);
    }
    
    // Setup event listeners for reply forms
    static setupReplyEventListeners() {
        const forms = document.querySelectorAll('.reply-form');
        console.log(`Setting up reply event listeners for ${forms.length} forms`);
        forms.forEach(form => {
            const newForm = DOMUtilities.removeEventListeners(form);
            const postId = parseInt(newForm.id.replace('reply-form-', ''));
            
            // Initialize character counter for this reply form
            const textarea = newForm.querySelector('#reply-content');
            const charCount = document.getElementById(`reply-char-count-${postId}`);
            
            if (textarea && charCount) {
                CharacterCounter.initialize(textarea, charCount, {
                    maxLength: 500,
                    warningThreshold: 400,
                    dangerThreshold: 450
                });
            }
            
            // Set up form submission manually to ensure proper context
            newForm.addEventListener('submit', async (e) => {
                e.preventDefault();
                
                const submitBtn = newForm.querySelector('button[type="submit"]');
                
                try {
                    // Set loading state
                    submitBtn.disabled = true;
                    submitBtn.textContent = 'Posting...';
                    
                    // Clear previous errors
                    FormErrorHelpers.smartClearErrors({});
                    
                    // Submit reply
                    const result = await this.submitReply(newForm, postId);
                    
                    if (result.success) {
                        // Clear form on success
                        const textarea = newForm.querySelector('.reply-textarea');
                        if (textarea) {
                            textarea.value = '';
                            // Reset character counter
                            CharacterCounter.reset(charCount, { maxLength: 500 });
                        }
                    } else {
                        // Handle errors
                        if (result.validationErrors) {
                            FormErrorHelpers.handleValidationErrors(result.validationErrors);
                        } else {
                            FormErrorHelpers.displayGeneralError(result.error || 'An error occurred');
                        }
                    }
                    
                } catch (error) {
                    console.error('Form submission error:', error);
                    FormErrorHelpers.displayGeneralError('An error occurred. Please try again.');
                } finally {
                    // Reset loading state
                    submitBtn.disabled = false;
                    submitBtn.textContent = 'Reply';
                }
            });
        });
    }
    
    // Handle successful reply submission
    static async handleReplySuccess(postId, replyData) {
        console.log('Adding reply to comments list:', replyData);
        
        // Add reply to comments list
        await this.addNewReplyToCommentsList(postId, replyData);
        this.updateCommentCount(postId);
    }
    
    // Submit reply (extracted from main handler)
    static async submitReply(form, postId) {
        const textarea = form.querySelector('.reply-textarea');
        const content = textarea?.value.trim();
        
        // Validate using model
        const validation = ReplyModel.validateReply(content);
        if (!validation.isValid) {
            // Return validation errors in the same format as backend
            return { 
                success: false, 
                validationErrors: {
                    FieldErrors: { "reply-content": validation.errors[0] }
                }
            };
        }
        
        // Additional validation using character counter
        const currentCount = CharacterCounter.getCount(textarea);
        if (currentCount > 500) {
            return { 
                success: false, 
                validationErrors: {
                    FieldErrors: { "reply-content": 'Reply cannot exceed 500 characters' }
                }
            };
        }
        
        // Submit using model
        const response = await ReplyModel.createReply({
            post_id: postId,
            content: content
        });
        
        if (!response) {
            return { success: false, error: 'Failed to post reply. Please try again.' };
        }
        
        if (!response.ok) {
            const errorData = await response.json();
            if (response.status === 422) {
                return { success: false, validationErrors: errorData };
            }
            return { success: false, error: 'An error occurred while posting your reply.' };
        }
        
        const responseData = await response.json();
        
        // Handle successful reply submission
        await this.handleReplySuccess(postId, responseData);
        
        return { success: true, data: responseData };
    }
    
    // Add new reply to comments list
    static async addNewReplyToCommentsList(postId, replyData) {
        try {
            const commentsList = document.getElementById(`comments-list-${postId}`);
            if (!commentsList) return;
            
            // Check if currently showing "no comments" message
            const noCommentsContainer = commentsList.querySelector('.no-comments-container');
            if (noCommentsContainer) {
                // Replace "no comments" message with the new reply
                const replyHTML = ReplyBoxView.generateNewReplyHTML(replyData);
                commentsList.innerHTML = replyHTML;
            } else {
                // Add to existing comments list at the top
                const replyHTML = ReplyBoxView.generateNewReplyHTML(replyData);
                commentsList.insertAdjacentHTML('afterbegin', replyHTML);
            }
        } catch (error) {
            console.error('Error adding reply to comments list:', error);
        }
    }
    
    // Update comment count in post stats
    static updateCommentCount(postId) {
        const commentCountElement = document.querySelector(`[data-post-id="${postId}"] .comment-count`);
        if (!commentCountElement) return;
        
        const currentText = commentCountElement.textContent;
        const match = currentText.match(/\d+/);
        const currentCount = match ? parseInt(match[0]) : 0;
        const newCount = currentCount + 1;
        
        // Update the text while preserving the structure
        const img = commentCountElement.querySelector('img');
        const imgHTML = img ? img.outerHTML : '';
        commentCountElement.innerHTML = `${imgHTML} ${newCount} comments`;
    }
}

// Global functions for onclick events
window.showReplyBox = (postId) => {
    RepliesController.showReplyBox(postId);
};

window.hideReplyBox = (postId) => {
    RepliesController.hideReplyBox(postId);
};
