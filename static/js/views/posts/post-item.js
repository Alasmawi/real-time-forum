import { formatTime } from "../../modules/utils/UI-UX/time-formatter.js";
import { AvatarGenerator } from "../../modules/utils/UI-UX/avatar-generator.js";

// Generate HTML for a single post item
export function generatePostHTML(post) {
    // Create the post element
    const postDiv = document.createElement('div');
    postDiv.className = 'post-item';
    postDiv.setAttribute('data-post-id', post.id);
    
    // Set the HTML structure without content
    postDiv.innerHTML = `
        <div class="post-header">
            <div class="post-avatar">
                ${AvatarGenerator.generateAvatarHTML(post.username, { size: 'medium' })}
            </div>
            <div class="post-username">${post.username}</div>
        </div>
        
        <div class="post-time">${formatTime(post.created_at)}</div>
        
        <div class="post-content">
            <p class="post-text"></p>
        </div>
            
        <div class="post-categories">
            ${post.categories && post.categories.length > 0 
                ? post.categories.map(cat => `<span class="post-category">${cat.name}</span>`).join('')
                : ''
            }
        </div>
        
        <div class="post-stats">
            <span class="comment-count" onclick="toggleComments(${post.id})">
                <img src="images/svg-small-comments.svg" alt="Comments" class="comment-icon">
                ${post.comment_count || 0} comments
            </span>
        </div>
        
        <div class="post-comments" id="comments-${post.id}">
            <div class="reply-box" id="reply-box-${post.id}" style="display: none;">
                <div class="reply-header">
                    <div class="reply-avatar">
                        ${AvatarGenerator.generateAvatarHTML(post.username, { size: 'small' })}
                    </div>
                    <div class="reply-username">Reply to ${post.username}</div>
                </div>
                <form class="reply-form" id="reply-form-${post.id}">
                    <textarea 
                        class="reply-textarea" 
                        placeholder="Write your reply..." 
                        maxlength="500"
                        required
                    ></textarea>
                    <div class="reply-actions">
                        <button type="submit" class="reply-submit-btn">Reply</button>
                        <button type="button" class="reply-cancel-btn" onclick="hideReplyBox(${post.id})">Cancel</button>
                    </div>
                </form>
            </div>
            <div class="comments-list" id="comments-list-${post.id}">
            </div>
        </div>
    `;
    
    // Safely set the post content using innerText
    const postTextElement = postDiv.querySelector('.post-text');
    if (postTextElement) {
        postTextElement.innerText = post.content;
    }
    
    return postDiv.outerHTML;
}
