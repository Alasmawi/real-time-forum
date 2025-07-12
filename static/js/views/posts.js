import AbstractView from "./abstract-view.js";
import { formatTime } from "../utils/time-formatter.js";

export default class PostsView extends AbstractView {
    constructor() {
        super();
        this.setTitle("Posts");
        this.categories = [];
        this.posts = [];
        this.selectedCategory = 'all';
    }

    async getHtml() {
        await this.loadCategories();
        await this.loadPosts();
        
        return `
        <div id="posts-container" style="margin-top: 30px; padding: 20px;">
            <!-- Categories Filter Section -->
            <div style="border: 2px solid #007bff; padding: 15px; margin-bottom: 20px; border-radius: 5px;">
                <h3>Filter by Category</h3>
                <div id="categories-filter">
                    ${this.generateCategoryFilters()}
                </div>
            </div>

            <!-- Posts Section -->
            <div style="border: 3px solid black; padding: 20px; border-radius: 5px;">
                <h2>Posts</h2>
                <div id="posts-list">
                    ${this.generatePostsHTML()}
                </div>
            </div>
        </div>
        `;
    }

    async loadCategories() {
        try {
            const response = await fetch('/protected/v1/categories');
            
            // Let ErrorHandler process the response
            const errorHandled = await window.ErrorHandler.handleResponse(response);
            if (errorHandled) {
                this.categories = [];
                return; // ErrorHandler took care of error (redirect, show error page, etc.)
            }
            
            this.categories = await response.json();
        } catch (error) {
            console.error('Error loading categories:', error);
            this.categories = [];
        }
    }

    async loadPosts() {
        try {
            const response = await fetch('/protected/v1/posts', {
                credentials: 'include'
            });
            
            // Let ErrorHandler process the response
            const errorHandled = await window.ErrorHandler.handleResponse(response);
            if (errorHandled) {
                this.posts = [];
                return; // ErrorHandler took care of error (redirect, show error page, etc.)
            }
            
            this.posts = await response.json();
        } catch (error) {
            console.error('Error loading posts:', error);
            this.posts = [];
        }
    }

    generateCategoryFilters() {
        let html = `
            <button class="category-filter ${this.selectedCategory === 'all' ? 'active' : ''}" 
                    data-category="all" 
                    style="margin: 5px; padding: 8px 15px; border: 1px solid #ddd; border-radius: 3px; cursor: pointer; background: ${this.selectedCategory === 'all' ? '#007bff' : '#fff'}; color: ${this.selectedCategory === 'all' ? '#fff' : '#000'};">
                All Posts
            </button>
        `;
        
        this.categories.forEach(category => {
            const isActive = this.selectedCategory === category.id.toString();
            html += `
                <button class="category-filter ${isActive ? 'active' : ''}" 
                        data-category="${category.id}" 
                        style="margin: 5px; padding: 8px 15px; border: 1px solid #ddd; border-radius: 3px; cursor: pointer; background: ${isActive ? '#007bff' : '#fff'}; color: ${isActive ? '#fff' : '#000'};">
                    ${category.name}
                </button>
            `;
        });
        
        return html;
    }

    generatePostsHTML() {
        if (this.posts.length === 0) {
            return '<p>No posts available.</p>';
        }

        const filteredPosts = this.filterPosts();
        
        if (filteredPosts.length === 0) {
            return '<p>No posts found for the selected category.</p>';
        }

        return filteredPosts.map(post => `
            <div class="post-item" style="border: 1px solid #ddd; margin: 15px 0; padding: 15px; border-radius: 5px; cursor: pointer; transition: background-color 0.2s;" 
                 data-post-id="${post.id}">
                <div class="post-header" style="margin-bottom: 10px;">
                    <h4 style="margin: 0 0 5px 0;">By ${post.username}</h4>
                    <small style="color: #666;">${formatTime(post.created_at)}</small>
                </div>
                
                <div class="post-content" style="margin: 10px 0;">
                    <p>${post.content}</p>
                </div>
                
                <div class="post-categories" style="margin: 10px 0;">
                    <strong>Categories: </strong>
                    ${post.categories && post.categories.length > 0 
                        ? post.categories.map(cat => `<span style="background: #e9ecef; padding: 2px 6px; border-radius: 3px; margin-right: 5px; font-size: 0.9em;">${cat.name}</span>`).join('')
                        : '<span style="color: #666;">No categories</span>'
                    }
                </div>
                
                <div class="post-stats" style="margin: 10px 0; color: #666;">
                    <span>💬 ${post.comments ? post.comments.length : 0} comments</span>
                    <span style="float: right; font-size: 0.9em;">Click to view comments</span>
                </div>
                
                <div class="post-comments" id="comments-${post.id}" style="display: none; margin-top: 15px; padding-top: 15px; border-top: 1px solid #eee;">
                    ${this.generateCommentsHTML(post.comments)}
                </div>
            </div>
        `).join('');
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

    filterPosts() {
        if (this.selectedCategory === 'all') {
            return this.posts;
        }
        
        return this.posts.filter(post => 
            post.categories && post.categories.some(cat => cat.id.toString() === this.selectedCategory)
        );
    }

    async getData() {
        // Category filter event listeners
        document.addEventListener('click', (e) => {
            if (e.target.classList.contains('category-filter')) {
                const categoryId = e.target.getAttribute('data-category');
                this.selectedCategory = categoryId;
                
                // Update active button styles
                document.querySelectorAll('.category-filter').forEach(btn => {
                    btn.style.background = '#fff';
                    btn.style.color = '#000';
                    btn.classList.remove('active');
                });
                
                e.target.style.background = '#007bff';
                e.target.style.color = '#fff';
                e.target.classList.add('active');
                
                // Re-render posts
                document.getElementById('posts-list').innerHTML = this.generatePostsHTML();
                this.attachPostClickListeners();
            }
        });

        // Post click event listeners for expanding comments
        this.attachPostClickListeners();
    }

    attachPostClickListeners() {
        document.querySelectorAll('.post-item').forEach(postElement => {
            postElement.addEventListener('click', async (e) => {
                // Don't trigger if clicking on category buttons
                if (e.target.classList.contains('category-filter')) {
                    return;
                }
                
                const postId = parseInt(postElement.getAttribute('data-post-id'));
                const commentsDiv = document.getElementById(`comments-${postId}`);
                
                if (commentsDiv.style.display === 'none') {
                    // Show loading state
                    commentsDiv.innerHTML = '<p style="color: #666; font-style: italic;">Loading comments...</p>';
                    commentsDiv.style.display = 'block';
                    postElement.style.backgroundColor = '#f8f9fa';
                    
                    // Fetch comments from the backend
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
                        
                        // Let ErrorHandler process the response
                        const errorHandled = await window.ErrorHandler.handleResponse(response);
                        if (errorHandled) {
                            commentsDiv.innerHTML = '<p style="color: #dc3545; font-style: italic;">Failed to load comments.</p>';
                            return; // ErrorHandler took care of error (redirect, show error page, etc.)
                        }
                        
                        const comments = await response.json();
                        commentsDiv.innerHTML = this.generateCommentsHTML(comments);
                    } catch (error) {
                        console.error('Error loading comments:', error);
                        commentsDiv.innerHTML = '<p style="color: #dc3545; font-style: italic;">Error loading comments.</p>';
                    }
                } else {
                    commentsDiv.style.display = 'none';
                    postElement.style.backgroundColor = '#fff';
                }
            });
            
            // Add hover effect
            postElement.addEventListener('mouseenter', () => {
                if (postElement.style.backgroundColor !== 'rgb(248, 249, 250)') {
                    postElement.style.backgroundColor = '#f5f5f5';
                }
            });
            
            postElement.addEventListener('mouseleave', () => {
                if (postElement.style.backgroundColor !== 'rgb(248, 249, 250)') {
                    postElement.style.backgroundColor = '#fff';
                }
            });
        });
    }
}
