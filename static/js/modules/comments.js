import { formatTime } from "../utils/time-formatter.js";

export class CommentsManager {
    constructor() {
        this.attachedForms = new Set();
    }

    generateCommentsHTML(comments) {
        if (!comments || comments.length === 0) {
            return '<p style="color: #666; font-style: italic;">No comments yet.</p>';
        }

        return `
            <h5 style="margin-bottom: 10px;">Comments (${comments.length})</h5>
            ${comments.map(comment => `
                <div class="comment-item" style="background: #f8f9fa; padding: 10px; margin: 5px 0; border-radius: 3px; border-left: 3px solid #007bff;">
                    <div class="comment-header" style="margin-bottom: 5px;">
                        <strong>${comment.username}</strong>
                        <small style="color: #666; margin-left: 10px;">${formatTime(comment.created_at)}</small>
                    </div>
                    <div class="comment-content">
                        <p style="margin: 0;">${comment.content}</p>
                    </div>
                </div>
            `).join('')}
        `;
    }

    generateCommentForm(postId) {
        return `
            <div class="comment-form" style="margin-top: 15px; padding: 15px; background: #f8f9fa; border-radius: 5px;">
                <h6>Add a comment</h6>
                <form class="new-comment-form" data-post-id="${postId}">
                    <textarea id="comment-content-${postId}" name="content" rows="3" cols="50" placeholder="Write your comment..." required style="width: 100%; margin-bottom: 10px; padding: 8px; border: 1px solid #ddd; border-radius: 3px;"></textarea><br>
                    <input type="submit" value="Post Comment" style="background: #007bff; color: white; border: none; padding: 8px 15px; border-radius: 3px; cursor: pointer;">
                    <div class="comment-message" style="margin-top: 10px; color: green;"></div>
                    <div class="comment-error" style="margin-top: 10px; color: red;"></div>
                </form>
            </div>
        `;
    }

    async loadComments(postId) {
        try {
            const response = await fetch('/protected/v1/comments', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    post_id: postId
                })
            });

            if (response.ok) {
                return await response.json();
            } else {
                throw new Error('Failed to load comments');
            }
        } catch (error) {
            console.error('Error loading comments:', error);
            throw error;
        }
    }

    async submitComment(postId, content) {
        const commentData = {
            post_id: postId,
            content: content
        };

        console.log('Submitting comment:', commentData);

        try {
            const response = await fetch('/protected/v1/newcomment', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(commentData),
            });

            if (response.ok) {
                return await response.json();
            } else {
                const errorData = await response.json();
                throw new Error(errorData.Error || 'Failed to post comment');
            }
        } catch (error) {
            console.error('Error posting comment:', error);
            throw error;
        }
    }

    attachCommentFormListeners() {
        document.querySelectorAll('.new-comment-form').forEach(form => {
            // Skip if already attached
            if (this.attachedForms.has(form)) {
                return;
            }
            this.attachedForms.add(form);

            form.addEventListener('submit', async (e) => {
                e.preventDefault();
                
                const postId = parseInt(form.getAttribute('data-post-id'));
                const content = document.getElementById(`comment-content-${postId}`).value;
                const messageDiv = form.querySelector('.comment-message');
                const errorDiv = form.querySelector('.comment-error');
                
                // Clear previous messages
                messageDiv.textContent = '';
                errorDiv.textContent = '';
                
                if (!content) {
                    errorDiv.textContent = 'Comment content is required';
                    return;
                }
                
                try {
                    await this.submitComment(postId, content);
                    messageDiv.textContent = 'Comment posted successfully!';
                    form.reset();
                    
                    // Refresh the comments for this post
                    setTimeout(async () => {
                        try {
                            const comments = await this.loadComments(postId);
                            const commentsDiv = document.getElementById(`comments-${postId}`);
                            commentsDiv.innerHTML = this.generateCommentsHTML(comments) + this.generateCommentForm(postId);
                            this.attachCommentFormListeners();
                        } catch (error) {
                            console.error('Error refreshing comments:', error);
                        }
                    }, 1000);
                    
                } catch (error) {
                    errorDiv.textContent = error.message || 'Network error occurred';
                }
            });
        });
    }

    async displayCommentsForPost(postId, commentsDiv) {
        try {
            commentsDiv.innerHTML = '<p style="color: #666; font-style: italic;">Loading comments...</p>';
            
            const comments = await this.loadComments(postId);
            commentsDiv.innerHTML = this.generateCommentsHTML(comments) + this.generateCommentForm(postId);
            this.attachCommentFormListeners();
        } catch (error) {
            commentsDiv.innerHTML = '<p style="color: #dc3545; font-style: italic;">Failed to load comments.</p>';
        }
    }

    isCommentFormElement(target) {
        return target.closest('.comment-form') ||
               target.tagName === 'TEXTAREA' ||
               target.tagName === 'INPUT';
    }
}
