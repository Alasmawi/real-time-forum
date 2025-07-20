// Reply Box View - Handles reply box HTML generation
import { AvatarGenerator } from "../../modules/utils/UI-UX/avatar-generator.js";
import { formatTime } from "../../modules/utils/UI-UX/time-formatter.js";

export class ReplyBoxView {
    
    // Generate HTML for reply box
    static generateReplyBoxHTML(postId, currentUser) {
        return `
            <div class="reply-box" id="reply-box-${postId}" style="display: none;">
                <div class="reply-header">
                    <div class="reply-avatar">
                        ${AvatarGenerator.generateAvatarHTML(currentUser.username, { size: 'medium' })}
                    </div>
                    <div class="reply-username">${currentUser.username}</div>
                </div>
                
                <form class="reply-form" id="reply-form-${postId}">
                    <!-- General error display (same styling as new post card) -->
                    <div id="error-messages" class="general-error"></div>
                    
                    <div class="form-group">
                        <textarea 
                            class="reply-textarea form-input" 
                            id="reply-content"
                            name="reply-content" 
                            placeholder="Write your reply..." 
                            maxlength="500"
                        ></textarea>
                        <span id="reply-content-error" class="form-error"></span>
                    </div>
                    
                    <div class="reply-footer">
                        <div class="character-counter" id="reply-char-count-${postId}">0/500 characters</div>
                        <div class="reply-actions">
                            <button type="submit" class="reply-submit-btn">
                                Reply
                            </button>
                        </div>
                    </div>
                </form>
            </div>
        `;
    }
    
    // Generate HTML for a new reply item (to be added to comments list)
    static generateNewReplyHTML(replyData) {
        // Create the reply element
        const replyDiv = document.createElement('div');
        replyDiv.className = 'comment-item';
        replyDiv.setAttribute('data-comment-id', replyData.comment_id || replyData.id || 'new');
        
        // Set the HTML structure without content
        replyDiv.innerHTML = `
            <div class="comment-header">
                <div class="comment-avatar">
                    ${AvatarGenerator.generateAvatarHTML(replyData.username, { size: 'medium' })}
                </div>
                <div class="comment-meta">
                    <div class="comment-username">${replyData.username}</div>
                    <div class="comment-time">${formatTime(replyData.created_at)}</div>
                </div>
            </div>
            
            <div class="comment-content">
                <p class="comment-text"></p>
            </div>
        `;
        
        // Safely set the reply content using innerText
        const commentTextElement = replyDiv.querySelector('.comment-text');
        if (commentTextElement) {
            commentTextElement.innerText = replyData.content;
        }
        
        return replyDiv.outerHTML;
    }
}
