// Category Controller - Handles category-related operations
import { CategoriesModel } from "../models/category-model.js";
import { CategoryFiltersView } from "../views/category-filters.js";

export class CategoryController {
    static categories = [];
    static selectedCategory = 'all';

    // Load categories using CategoriesModel
    static async loadCategories() {
        try {
            // Use CategoriesModel to fetch categories
            this.categories = await CategoriesModel.fetchCategories();
            return this.categories;
        } catch (error) {
            console.error('Error loading categories:', error);
            this.categories = [];
            return [];
        }
    }

    // Generate category filter buttons HTML using view
    static generateCategoryFilters(categories, selectedCategory) {
        return CategoryFiltersView.generateCategoryFiltersHTML(categories, selectedCategory);
    }

    // Filter posts by category
    static filterPosts(posts, selectedCategory) {
        if (selectedCategory === 'all') {
            return posts;
        }
        
        return posts.filter(post => 
            post.categories && post.categories.some(cat => cat.id.toString() === selectedCategory)
        );
    }

    // Setup category filter event listeners
    static setupCategoryFilters(onCategoryChange) {
        document.addEventListener('click', (e) => {
            if (e.target.classList.contains('category-filter')) {
                const categoryId = e.target.getAttribute('data-category');
                this.selectedCategory = categoryId;
                
                // Update active button styles - remove active from all, add to clicked
                document.querySelectorAll('.category-filter').forEach(btn => {
                    btn.classList.remove('active');
                });
                
                e.target.classList.add('active');
                
                // Call the callback function
                if (onCategoryChange) {
                    onCategoryChange(categoryId);
                }
            }
        });
    }
}
