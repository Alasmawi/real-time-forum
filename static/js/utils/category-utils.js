// Category utility functions for reuse across components

export async function loadCategories() {
    try {
        const response = await fetch('/protected/v1/categories');
        if (response.ok) {
            return await response.json();
        } else {
            console.error('Failed to load categories');
            return [];
        }
    } catch (error) {
        console.error('Error loading categories:', error);
        return [];
    }
}

export function generateCategoryFilters(categories, selectedCategory = 'all') {
    const createButton = (id, name, isActive) => {
        return `<button class="category-filter ${isActive ? 'active' : ''}" data-category="${id}">${name}</button>`;
    };
    
    let html = createButton('all', 'All Posts', selectedCategory === 'all');
    
    categories.forEach(category => {
        const isActive = selectedCategory === category.id.toString();
        html += createButton(category.id, category.name, isActive);
    });
    
    return html;
}

export function handleCategoryFilter(e, callback) {
    const categoryId = e.target.getAttribute('data-category');
    
    // Update active button classes
    document.querySelectorAll('.category-filter').forEach(btn => {
        btn.classList.remove('active');
    });
    
    e.target.classList.add('active');
    
    // Execute callback with selected category
    if (callback) {
        callback(categoryId);
    }
    
    return categoryId;
}

export function attachCategoryClickListeners(callback) {
    document.addEventListener('click', (e) => {
        if (e.target.classList.contains('category-filter')) {
            handleCategoryFilter(e, callback);
        }
    });
}

export function filterItemsByCategory(items, selectedCategory, categoryField = 'categories') {
    if (selectedCategory === 'all') {
        return items;
    }
    
    return items.filter(item => {
        const categories = item[categoryField];
        return categories && categories.some(cat => cat.id.toString() === selectedCategory);
    });
}
