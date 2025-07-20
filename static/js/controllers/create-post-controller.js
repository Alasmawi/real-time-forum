// New Post Controller - Handles inline post creation business logic
import { CategoryController } from "./category-controller.js";
import { PostsModel } from "../models/posts-model.js";
import { UserModel } from "../models/user-model.js";
import { FormErrorHelpers } from "../modules/utils/errors/form-error-helpers.js";
import { CharacterCounter } from "../modules/utils/character-counter.js";
import { generatePostHTML } from "../views/posts/post-item.js";
import { FormHandler } from "../modules/utils/form-handler.js";

export class NewPostController {

    // Setup event listeners for new post form (called after HTML is rendered by view)
    static setupEventListeners() {
        const form = document.getElementById('new-post-form');
        const textarea = document.getElementById('post-content');
        const charCount = document.getElementById('char-count');
        const categoryBtns = document.querySelectorAll('.new-post-card .category-btn');
        const expandableContent = document.querySelectorAll('.new-post-card .expandable-content');

        // Initialize character counter using utility
        CharacterCounter.initialize(textarea, charCount, {
            maxLength: 500,
            warningThreshold: 400,
            dangerThreshold: 450
        });

        // Expand new post card when textarea is focused
        if (textarea) {
            textarea.addEventListener('focus', () => {
                this.expandNewPostCard();
                // Hide filters if they're shown
                this.hideFiltersIfShown();
            });
        }

        // Category selection
        categoryBtns.forEach(btn => {
            btn.addEventListener('click', (e) => {
                e.preventDefault();
                e.target.classList.toggle('selected');
            });
        });

        // Form submission using FormHandler
        if (form) {
            FormHandler.handleFormSubmission(
                '#new-post-form',
                async (form) => this.submitPost(form),
                {
                    loadingText: 'Posting...',
                    defaultText: 'Post',
                    clearOnSuccess: false // We handle clearing ourselves
                }
            );
        }
    }

    // Expand the new post card to show categories and post button
    static expandNewPostCard() {
        const expandableContent = document.querySelectorAll('.new-post-card .expandable-content');
        expandableContent.forEach(element => {
            element.classList.add('expanded');
        });
    }

    // Collapse the new post card to minimized state
    static collapseNewPostCard() {
        const expandableContent = document.querySelectorAll('.new-post-card .expandable-content');
        expandableContent.forEach(element => {
            element.classList.remove('expanded');
        });
    }

    // Hide filters if they're currently shown
    static hideFiltersIfShown() {
        const categoryFilters = document.getElementById('category-filters');
        const toggleBtn = document.getElementById('toggle-categories-btn');
        
        if (categoryFilters && categoryFilters.classList.contains('expanded')) {
            categoryFilters.classList.remove('expanded');
            if (toggleBtn) {
                toggleBtn.textContent = 'Show Filters';
            }
        }
    }

    // Submit post (used by FormHandler)
    static async submitPost(form) {
        // Collect form data
        const content = document.getElementById('post-content').value;
        const selectedCategories = Array.from(document.querySelectorAll('.category-btn.selected'))
            .map(btn => parseInt(btn.dataset.categoryId));

        const postData = {
            content: content,
            categories: selectedCategories
        };

        // Use PostsModel to create post
        const response = await PostsModel.createPost(postData);

        if (response === null) {
            return { success: false, error: "Failed to create post. Please try again." };
        }

        if (!response.ok) {
            const errorData = await response.json();
            if (response.status === 422) {
                return { success: false, validationErrors: errorData };
            }
            return { success: false, error: "An error occurred while creating the post. Please try again." };
        }

        const postResponse = await response.json();
        console.log("Post created successfully");
        
        // Clear form and add to feed
        this.clearPostForm();
        await this.addNewPostToList(postResponse);
        
        return { success: true, data: postResponse };
    }

    // Clear the post form
    static clearPostForm() {
        const textarea = document.getElementById('post-content');
        const charCount = document.getElementById('char-count');
        const categoryBtns = document.querySelectorAll('.new-post-card .category-btn.selected');

        if (textarea) {
            textarea.value = '';
            textarea.blur(); // Remove focus from textarea
        }

        // Reset character counter using utility
        CharacterCounter.reset(charCount, { maxLength: 500 });

        // Deselect all categories only on successful submission
        categoryBtns.forEach(btn => {
            btn.classList.remove('selected');
        });

        // Collapse the new post card back to minimized state
        this.collapseNewPostCard();
    }

    // Add new post to the top of the feed
    static async addNewPostToList(postResponse) {
        try {
            // Get current user info
            const user = await UserModel.fetchCurrentUser();
            if (!user) return;

            // Get categories for the post
            const categories = await CategoryController.loadCategories();
            const postCategories = postResponse.categories || [];
            const categoryNames = postCategories.map(catId => {
                const category = categories.find(cat => cat.id === catId);
                return category ? category.name : '';
            }).filter(name => name);

            // Create a post object that matches the expected format
            const newPost = {
                id: postResponse.post_id,
                content: postResponse.content,
                username: user.username, // Changed from 'author' to 'username' to match view expectations
                created_at: new Date().toISOString(),
                categories: categoryNames.map(name => ({ name })), // Convert to objects to match view expectations
                comment_count: 0
            };
            
            // Generate HTML for the new post using view
            const postHTML = generatePostHTML(newPost);
            
            // Find the container and new post card
            const postsContainer = document.getElementById('posts-list');
            const newPostCard = postsContainer.querySelector('.new-post-card');
            
            if (postsContainer && newPostCard) {
                // Insert the new post right after the new post card
                newPostCard.insertAdjacentHTML('afterend', postHTML);
                
                // Setup comment listeners for the new post
                const { CommentController } = await import('./comments-controller.js');
                CommentController.setupCommentClickListeners();
            }
        } catch (error) {
            console.error('Error adding new post to list:', error);
        }
    }
}
