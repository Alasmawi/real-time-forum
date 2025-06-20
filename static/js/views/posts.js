import AbstractView from "./abstract-view.js";
import { CommentsManager } from "../modules/comments.js";
import { loadCategories, generateCategoryFilters, attachCategoryClickListeners, filterItemsByCategory } from "../utils/category-utils.js";
import { loadPosts, generatePostsHTML, attachPostClickListeners } from "../utils/post-utils.js";

export default class PostsView extends AbstractView {
    constructor() {
        super();
        this.setTitle("Posts");
        this.categories = [];
        this.posts = [];
        this.selectedCategory = 'all';
        this.commentsManager = new CommentsManager();
    }

    async getHtml() {
        return `
        <div id="posts-container">
            <div class="categories-filter-section">
                <h3>Filter by Category</h3>
                <div id="categories-filter">
                    <!-- Categories will be loaded here -->
                </div>
            </div>

            <div class="posts-section">
                <h2>Posts</h2>
                <div id="posts-list">
                    <!-- Posts will be loaded here -->
                </div>
            </div>
        </div>
        `;
    }

    async getData() {
        await this.loadData();
        this.renderContent();
        this.attachEventListeners();
    }

    async loadData() {
        this.categories = await loadCategories();
        this.posts = await loadPosts();
    }

    renderContent() {
        // Render categories filter
        document.getElementById('categories-filter').innerHTML = generateCategoryFilters(this.categories, this.selectedCategory);
        
        // Render posts
        document.getElementById('posts-list').innerHTML = this.renderPosts();
    }

    renderPosts() {
        const filteredPosts = filterItemsByCategory(this.posts, this.selectedCategory);
        
        if (filteredPosts.length === 0) {
            return '<p>No posts found for the selected category.</p>';
        }

        return generatePostsHTML(filteredPosts, this.commentsManager);
    }

    attachEventListeners() {
        // Category filter listeners
        attachCategoryClickListeners((categoryId) => {
            this.selectedCategory = categoryId;
            this.updatePostsList();
        });

        // Post click listeners
        attachPostClickListeners(this.commentsManager);
        
        // Comment form listeners
        this.commentsManager.attachCommentFormListeners();
    }

    updatePostsList() {
        document.getElementById('posts-list').innerHTML = this.renderPosts();
        attachPostClickListeners(this.commentsManager);
    }
}
