// Post utility functions for reuse across components
import { formatTime } from "./time-formatter.js";

export async function loadPosts() {
    try {
        const response = await fetch('/protected/v1/posts', {
            credentials: 'include'
        });
        if (response.ok) {
            return await response.json();
        } else {
            console.error('Failed to load posts');
            return [];
        }
    } catch (error) {
        console.error('Error loading posts:', error);
        return [];
    }
}

export function generatePostsHTML(posts, commentsManager) {
    if (posts.length === 0) {
        return '<p>No posts available.</p>';
    }

    return posts.map(post => {
        const categories = post.categories && post.categories.length > 0 
            ? post.categories.map(cat => `<span class="category-tag">${cat.name}</span>`).join('')
            : '<span>No categories</span>';
        
        const commentsHTML = commentsManager.generateCommentsHTML(post.comments);
        const commentForm = commentsManager.generateCommentForm(post.id);
        
        return `
            <div class="post-item" data-post-id="${post.id}">
                <div class="post-header">
                    <h4>By ${post.username}</h4>
                    <small>${formatTime(post.created_at)}</small>
                </div>
                
                <div class="post-content">
                    <p>${post.content}</p>
                </div>
                
                <div class="post-categories">
                    <strong>Categories: </strong>
                    ${categories}
                </div>
                
                <div class="post-stats">
                    <span>💬 Comments</span>
                    <span class="click-hint">Click to view comments</span>
                </div>
                
                <div class="post-comments" id="comments-${post.id}" style="display: none;">
                    ${commentsHTML}
                    ${commentForm}
                </div>
            </div>
        `;
    }).join('');
}

export function attachPostClickListeners(commentsManager) {
    document.querySelectorAll('.post-item').forEach(postElement => {
        postElement.addEventListener('click', async (e) => {
            await handlePostClick(e, postElement, commentsManager);
        });
    });
}

export async function handlePostClick(e, postElement, commentsManager) {
    if (e.target.classList.contains('category-filter') || 
        commentsManager.isCommentFormElement(e.target)) {
        return;
    }
    
    const postId = parseInt(postElement.getAttribute('data-post-id'));
    const commentsDiv = document.getElementById(`comments-${postId}`);
    
    if (commentsDiv.style.display === 'none') {
        commentsDiv.style.display = 'block';
        await commentsManager.displayCommentsForPost(postId, commentsDiv);
    } else {
        commentsDiv.style.display = 'none';
    }
}
