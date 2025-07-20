// Comment Controller - Handles comment-related operations
import { CommentsModel } from "../models/comments-model.js";
import { RepliesController } from "./reply-controller.js";
import { generateCommentsListHTML } from "../views/posts/comment-item.js";
import { DOMUtilities } from "../modules/utils/dom/dom-utilities.js";
import { ReplyBoxView } from "../views/posts/reply-box.js";
import { UserModel } from "../models/user-model.js";
import PaginationManager from "../modules/pagination/pagination-manager.js";

export class CommentController {
    static paginationManager = new PaginationManager();
    
    // Generate comments HTML using view
    static generateCommentsHTML(comments) {
        if (!comments || comments.length === 0) {
            return `
                <div class="no-comments-container">
                    <div class="no-comments-message">
                        <p>No comments yet.</p>
                        <span class="no-comments-subtitle">Be the first to comment!</span>
                    </div>
                </div>
            `;
        }

        return generateCommentsListHTML(comments);
    }

    // Initialize pagination for comments of a specific post
    static async initializeCommentsPagination(postId) {
        const commentsType = `comments-${postId}`;
        const containerId = `comments-${postId}`;
        
        await this.paginationManager.initializePagination(
            commentsType,
            '/guest/v1/comments',
            containerId,
            (comments, replace = false) => {
                this.renderComments(postId, comments, replace);
            },
            { post_id: postId }
        );
    }
    
    // Render comments to a specific post's comments section
    static renderComments(postId, comments, replace = false) {
        const commentsListDiv = document.getElementById(`comments-list-${postId}`);
        if (!commentsListDiv) return;
        
        const commentsHTML = this.generateCommentsHTML(comments);
        
        if (replace) {
            commentsListDiv.innerHTML = commentsHTML;
        } else {
            commentsListDiv.insertAdjacentHTML('beforeend', commentsHTML);
        }
        
        // Setup reply event listeners after rendering
        RepliesController.setupReplyEventListeners();
    }
    
    // Load comments for a specific post using CommentsModel (legacy method)
    static async loadComments(postId, params = {}) {
        try {
            const response = await CommentsModel.fetchCommentsByPostId(postId, params);
            return response.comments || []; // Return just the comments array
        } catch (error) {
            console.error('CommentController.loadComments error:', error);
            return []; // Return empty array on error
        }
    }

    // Setup comment click listeners for expanding/collapsing comments
    static setupCommentClickListeners() {
        document.querySelectorAll('.post-item').forEach(postElement => {
            // Remove any existing listeners to prevent duplicates
            const newElement = DOMUtilities.removeEventListeners(postElement);
            
            newElement.addEventListener('click', async (e) => {
                // Don't trigger if clicking on category buttons or inside reply box
                if (e.target.classList.contains('category-filter') || 
                    e.target.closest('.reply-box') || 
                    e.target.closest('.post-comments')) {
                    return;
                }
                
                const postId = parseInt(newElement.getAttribute('data-post-id'));
                const commentsDiv = document.getElementById(`comments-${postId}`);
                
                if (!commentsDiv.classList.contains('expanded')) {
                    // Get current user and set up reply box immediately
                    const currentUser = await UserModel.fetchCurrentUser();
                    
                    // Show reply box first, then loading message for comments
                    commentsDiv.innerHTML = `
                        ${currentUser ? ReplyBoxView.generateReplyBoxHTML(postId, currentUser) : ''}
                        <div id="comments-list-${postId}">
                            <p class="comments-loading">Loading comments...</p>
                        </div>
                    `;
                    
                    // Expand the comments section
                    commentsDiv.classList.add('expanded');
                    
                    // Show reply box immediately
                    RepliesController.showReplyBox(postId);
                    
                    // Scroll to ensure reply box is visible
                    setTimeout(() => {
                        const replyBox = document.getElementById(`reply-box-${postId}`);
                        if (replyBox) {
                            replyBox.scrollIntoView({ 
                                behavior: 'smooth', 
                                block: 'center' 
                            });
                        }
                    }, 100);
                    
                    // Initialize pagination for comments
                    try {
                        // First clear the loading message
                        const commentsListDiv = document.getElementById(`comments-list-${postId}`);
                        if (commentsListDiv) {
                            commentsListDiv.classList.remove('comments-loading');
                            commentsListDiv.innerHTML = ''; // Clear loading message
                        }
                        
                        // Initialize pagination which will load the first page
                        await CommentController.initializeCommentsPagination(postId);
                        
                    } catch (error) {
                        const commentsListDiv = document.getElementById(`comments-list-${postId}`);
                        if (commentsListDiv) {
                            commentsListDiv.innerHTML = '<p class="comments-error">Error loading comments.</p>';
                        }
                    }
                } else {
                    // Collapse the comments section
                    commentsDiv.classList.remove('expanded');
                }
            });
        });
    }
}
