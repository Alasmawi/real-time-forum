// Posts Controller - Handles all posts page functionality with PAGINATION ONLY
import PaginationManager from "../modules/pagination/pagination-manager.js";
import { AppInit } from "../modules/utils/app-init.js";
import { CategoryController } from "./category-controller.js";
import { CommentController } from "./comments-controller.js";
import { NewPostController } from "./create-post-controller.js";
import { generatePostHTML } from "../views/posts/post-item.js";
import { CategoryFiltersView } from "../views/category-filters.js";

export class PostsController {
    static paginationManager = new PaginationManager();


    // Generate HTML for multiple posts
    static generatePostsListHTML(posts) {
        if (posts.length === 0) {
            return '<p>No posts available.</p>';
        }

        return posts.map(post => generatePostHTML(post)).join('');
    }

    // Toggle comments visibility
    static toggleComments(postId) {
        const commentsDiv = document.getElementById(`comments-${postId}`);
        if (commentsDiv) {
            commentsDiv.style.display = commentsDiv.style.display === 'none' ? 'block' : 'none';
        }
    }

    // Render posts to the DOM (used by pagination manager)
    static renderPosts(posts, replace = false) {
        const postsContainer = document.getElementById('posts-list');
        if (!postsContainer) return;

        const postsHTML = this.generatePostsListHTML(posts);
        
        if (replace) {
            // Preserve the new post card when replacing content
            const newPostCard = postsContainer.querySelector('.new-post-card');
            postsContainer.innerHTML = postsHTML;
            
            // Re-insert the new post card at the beginning if it existed
            if (newPostCard) {
                postsContainer.insertBefore(newPostCard, postsContainer.firstChild);
            }
        } else {
            postsContainer.insertAdjacentHTML('beforeend', postsHTML);
        }
        
        // Setup comment click listeners after rendering
        CommentController.setupCommentClickListeners();
        
        // Make toggleComments available globally
        window.toggleComments = this.toggleComments;
    }

    // Initialize the entire posts page
    static async initializePage() {
        // Initialize navigation and user list
        await AppInit.initializeAuthenticatedApp();

        // Initialize categories card
        await this.initializeCategoriesCard();

        // Setup new post card event listeners (HTML is already rendered by view)
        NewPostController.setupEventListeners();

        // Initialize pagination manager 
        await this.paginationManager.initializePagination(
            'posts',                    // type
            '/guest/v1/posts',      // endpoint
            'content-card',             // containerId (the actual scrollable container)
            (posts, replace = false) => {
                const filteredPosts = CategoryController.filterPosts(posts, CategoryController.selectedCategory);
                PostsController.renderPosts(filteredPosts, replace);
            },
            { category: CategoryController.selectedCategory } // params
        );

        // Setup category filter listeners with callback
        CategoryController.setupCategoryFilters((categoryId) => {
            // Clear pagination and reinitialize with new category
            this.paginationManager.clearPagination('posts');
            this.paginationManager.initializePagination(
                'posts',
                '/guest/v1/posts', 
                'content-card',
                (posts, replace = false) => {
                    const filteredPosts = CategoryController.filterPosts(posts, categoryId);
                    PostsController.renderPosts(filteredPosts, replace);
                },
                { category: categoryId }
            );
        });

        // Setup comment click listeners
        CommentController.setupCommentClickListeners();
    }

    // Initialize categories card by injecting it into middle-column
    static async initializeCategoriesCard() {
        const middleColumn = document.querySelector('.middle-column');
        const contentCard = document.querySelector('.content-card');
        
        // Check if categories card already exists
        const existingCategoriesCard = middleColumn?.querySelector('.categories-card');
        if (existingCategoriesCard) {
            return; // Categories card already exists, don't create a new one
        }
        
        if (middleColumn && contentCard) {
            const categories = await CategoryController.loadCategories();
            
            // Create categories card element using view
            const categoriesCard = document.createElement('div');
            categoriesCard.innerHTML = CategoryFiltersView.generateCategoriesCardHTML(categories, 'all');
            
            // Insert categories card before the content card in middle column
            middleColumn.insertBefore(categoriesCard.firstElementChild, contentCard);
            
            // Setup toggle functionality
            this.setupCategoriesToggle();
        }
    }

    // Setup categories toggle functionality
    static setupCategoriesToggle() {
        const toggleBtn = document.getElementById('toggle-categories-btn');
        const categoryFilters = document.getElementById('category-filters');
        
        if (toggleBtn && categoryFilters) {
            toggleBtn.addEventListener('click', () => {
                const isVisible = categoryFilters.classList.contains('expanded');
                
                if (isVisible) {
                    // Hide categories
                    categoryFilters.classList.remove('expanded');
                    toggleBtn.textContent = 'Show Filters';
                    
                    // Don't collapse new post card when hiding filters
                } else {
                    // Show categories
                    categoryFilters.classList.add('expanded');
                    toggleBtn.textContent = 'Hide Filters';
                    
                    // Collapse new post card when showing filters
                    NewPostController.collapseNewPostCard();
                }
            });
        }
    }
}
