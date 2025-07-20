// API Fetcher - Handles paginated API requests with type-specific page sizes
import ErrorHandler from '../utils/errors/error-handler.js';

class ApiFetcher {
    constructor() {
        this.typeSizes = new Map(); // Store page sizes per type
        this.defaultPageSize = 25;
    }

    // Set page size for a specific type
    setPageSizeForType(type, pageSize) {
        this.typeSizes.set(type, pageSize);
    }

    // Get page size for a specific type
    getPageSizeForType(type) {
        return this.typeSizes.get(type) || this.defaultPageSize;
    }

    // Fetch a page of data from an endpoint with type-specific page size
    async fetchPage(endpoint, page, additionalParams = {}, type = 'default') {
        const pageSize = this.getPageSizeForType(type);
        const pageNumber = page + 1; // Convert 0-based to 1-based pagination
        
        const url = new URL(endpoint, window.location.origin);
        url.searchParams.set('page', pageNumber);
        url.searchParams.set('limit', pageSize);
        
        // Add additional parameters
        Object.entries(additionalParams).forEach(([key, value]) => {
            url.searchParams.set(key, value);
        });

        const response = await fetch(url, {
            credentials: 'include'
        });

        // Use your error handler for API responses
        const errorHandled = await ErrorHandler.handleResponse(response);
        if (errorHandled) {
            throw new Error(`API Error ${response.status}`);
        }

        const data = await response.json();
        
        // Handle new paginated response format
        if (data.posts) {
            return { items: data.posts, total: data.pagination.total_count };
        } else if (data.comments) {
            return { items: data.comments, total: data.pagination.total_count };
        } else if (data.users) {
            return { 
                items: data.users, 
                total: data.pagination.total_count,
                pagination: data.pagination, // Pass through full pagination object
                online_count: data.online_count,
                offline_count: data.offline_count,
                rawData: data // Include full response for user list
            };
        } else if (data.messages) {
            return { 
                items: data.messages, 
                total: data.pagination.total_count,
                target_user_id: data.target_user_id
            };
        } else if (Array.isArray(data)) {
            return { items: data, total: data.length };
        } else {
            return { items: data, total: data.length };
        }
    }

    // Fetch posts with pagination
    async fetchPosts(page, categoryFilter = null) {
        const params = {};
        if (categoryFilter) {
            params.category = categoryFilter;
        }
        return await this.fetchPage('/guest/v1/posts', page, params, 'posts');
    }

    // Fetch users with pagination
    async fetchUsers(page) {
        return await this.fetchPage('/protected/v1/user-list', page, {}, 'users');
    }

    // Fetch messages with pagination
    async fetchMessages(page, userId) {
        const params = { user_id: userId };
        return await this.fetchPage('/protected/v1/message-history', page, params, 'messages');
    }

    // Fetch comments with pagination
    async fetchComments(page, postId) {
        const params = { post_id: postId };
        return await this.fetchPage('/guest/v1/comments', page, params, 'comments');
    }

    // Fetch categories with pagination
    async fetchCategories(page) {
        return await this.fetchPage('/guest/v1/categories', page, {}, 'categories');
    }

    // Update default page size
    updateDefaultPageSize(pageSize) {
        this.defaultPageSize = pageSize;
    }

    // Get all configured page sizes
    getAllPageSizes() {
        const sizes = {};
        for (const [type, size] of this.typeSizes) {
            sizes[type] = size;
        }
        return sizes;
    }
}

export default ApiFetcher;
