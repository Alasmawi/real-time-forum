import { formatTime } from "../../modules/utils/UI-UX/time-formatter.js";
import { AvatarGenerator } from "../../modules/utils/UI-UX/avatar-generator.js";

// Generate HTML for a single comment
export function generateCommentHTML(comment) {
    // Create the comment element
    const commentDiv = document.createElement('div');
    commentDiv.className = 'comment-item';
    commentDiv.setAttribute('data-comment-id', comment.id);
    
    // Set the HTML structure without content
    commentDiv.innerHTML = `
        <div class="comment-header">
            <div class="comment-avatar">
                ${AvatarGenerator.generateAvatarHTML(comment.username, { size: 'medium' })}
            </div>
            <div class="comment-meta">
                <div class="comment-username">${comment.username}</div>
                <div class="comment-time">${formatTime(comment.created_at)}</div>
            </div>
        </div>
        
        <div class="comment-content">
            <p class="comment-text"></p>
        </div>
    `;
    
    // Safely set the comment content using innerText
    const commentTextElement = commentDiv.querySelector('.comment-text');
    if (commentTextElement) {
        commentTextElement.innerText = comment.content;
    }
    
    return commentDiv.outerHTML;
}

// Generate HTML for multiple comments
export function generateCommentsListHTML(comments) {
    if (!comments || comments.length === 0) {
        return '<p class="no-comments">No comments yet.</p>';
    }

    return comments.map(comment => generateCommentHTML(comment)).join('');
}
