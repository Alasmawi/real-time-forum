// Category Filters View - Handles category filter HTML generation
export class CategoryFiltersView {
    
    // Generate category filter buttons HTML
    static generateCategoryFiltersHTML(categories, selectedCategory) {
        let html = `
            <button class="category-filter ${selectedCategory === 'all' ? 'active' : ''}" 
                    data-category="all">
                All Posts
            </button>
        `;
        
        categories.forEach(category => {
            const isActive = selectedCategory === category.id.toString();
            html += `
                <button class="category-filter ${isActive ? 'active' : ''}" 
                        data-category="${category.id}">
                    ${category.name}
                </button>
            `;
        });
        
        return html;
    }
    
    // Generate complete categories card HTML
    static generateCategoriesCardHTML(categories, selectedCategory = 'all') {
        return `
            <div class="categories-card">
                <div class="categories-header">
                    <button id="toggle-categories-btn" class="toggle-categories-btn">
                        Show Filters
                    </button>
                </div>
                <div class="category-filters" id="category-filters">
                    ${this.generateCategoryFiltersHTML(categories, selectedCategory)}
                </div>
            </div>
        `;
    }
}
