// Main Pagination Manager - Orchestrates all pagination modules with extensible page sizes
import CacheManager from './cache-manager.js';
import ScrollManager from './scroll-manager.js';
import ApiFetcher from './api-fetcher.js';

class PaginationManager {
    constructor() {
        this.cacheManager = new CacheManager();
        this.scrollManager = new ScrollManager();
        this.apiFetcher = new ApiFetcher();
        this.loadingStates = new Map();
        
        // Default page sizes for different data types based on overhead
        this.pageSizes = {
            posts: 20,      
            users: 50,    
            messages: 25,   // Changed from 30 to 25 to match private chat expectation
            comments: 25,  
            categories: 100 
        };
        
        this.config = {
            maxCacheSize: 500,
            defaultPageSize: 25 // Fallback for unspecified types
        };
    }

    // Set custom page size for a specific type
    setPageSize(type, pageSize) {
        this.pageSizes[type] = pageSize;
        this.apiFetcher.setPageSizeForType(type, pageSize);
    }

    // Get page size for a specific type
    getPageSize(type) {
        // Check if exact type exists first
        if (this.pageSizes[type]) {
            return this.pageSizes[type];
        }
        
        // For specific message types like 'messages-123', check base 'messages' type
        if (type.startsWith('messages-')) {
            return this.pageSizes['messages'] || this.config.defaultPageSize;
        }
        
        // For other types, use default
        return this.config.defaultPageSize;
    }

    // Initialize pagination for a specific type
    async initializePagination(type, endpoint, containerId, renderCallback, params = {}) {
        // Set the page size for this type in the API fetcher
        const pageSize = this.getPageSize(type);
        this.apiFetcher.setPageSizeForType(type, pageSize);
        
        // Initialize cache
        this.cacheManager.initializeCache(type, endpoint, params);
        
        // Try to load from localStorage first
        const cachedData = this.cacheManager.loadFromLocalStorage(type);
        if (cachedData && cachedData.items.length > 0) {
            // Use cached data for immediate display
            renderCallback(cachedData.items.slice(0, pageSize));
        }
        
        // Load fresh data from server
        const initialItems = await this.loadInitialData(type);
        if (initialItems && (Array.isArray(initialItems) ? initialItems.length > 0 : initialItems.items && initialItems.items.length > 0)) {
            renderCallback(initialItems, true); // true = replace existing
        }
        
        // Set up infinite scroll with type-specific options
        const scrollOptions = type === 'users' ? 
            { useThrottle: true, useDebounce: false } : 
            {}; // Default options for other types
            
        // Skip setting up normal scroll for message types (they use reverse scroll)
        const isMessageType = type.startsWith('messages-');
        if (!isMessageType) {
            this.scrollManager.setupInfiniteScroll(type, containerId, async () => {
                const newItems = await this.loadMoreData(type);
                if (newItems && ((Array.isArray(newItems) && newItems.length > 0) || (newItems.items && newItems.items.length > 0))) {
                    renderCallback(newItems, false); // false = append
                }
            }, scrollOptions);
        }
    }

    // Load initial page of data
    async loadInitialData(type) {
        const cache = this.cacheManager.getCache(type);
        if (!cache) return [];

        this.setLoading(type, true);
        
        try {
            let response;
            const pageSize = this.getPageSize(type);
            
            // Route to appropriate fetcher method based on type
            switch (true) {
                case type === 'posts':
                    response = await this.apiFetcher.fetchPosts(0, cache.params.category);
                    break;
                case type === 'users':
                    response = await this.apiFetcher.fetchUsers(0);
                    break;
                case type.startsWith('messages'):
                    response = await this.apiFetcher.fetchMessages(0, cache.params.user_id);
                    break;
                default:
                    response = await this.apiFetcher.fetchPage(cache.endpoint, 0, cache.params);
            }
            
            const items = response.items || [];
            const totalCount = response.total;
            
            // Reset cache with fresh data
            cache.items = items;
            cache.currentPage = 1;
            cache.totalCount = totalCount;
            cache.hasMore = items.length < totalCount; // Has more if loaded items < total available
            
            console.log(`Initial load for ${type}:`, {
                itemsLoaded: items.length,
                totalCount: totalCount,
                hasMore: cache.hasMore,
                pageSize: pageSize
            });
            
            this.cacheManager.saveToLocalStorage(type, cache);
            
            // For users, return the full response to preserve metadata, otherwise just items
            return type === 'users' ? response : items;
        } catch (error) {
            console.error(`Error loading initial ${type}:`, error);
            return [];
        } finally {
            this.setLoading(type, false);
        }
    }

    // Load more data (next page)
    async loadMoreData(type) {
        const cache = this.cacheManager.getCache(type);
        console.log(`loadMoreData called for type: ${type}`, { 
            cacheExists: !!cache, 
            isLoading: this.isLoading(type), 
            hasMore: cache?.hasMore,
            currentPage: cache?.currentPage,
            itemsCount: cache?.items?.length 
        });
        
        if (!cache || this.isLoading(type) || !cache.hasMore) {
            console.log(`Returning empty array for ${type}:`, {
                noCache: !cache,
                isLoading: this.isLoading(type),
                noMore: !cache?.hasMore
            });
            return [];
        }

        this.setLoading(type, true);
        
        try {
            let response;
            const pageSize = this.getPageSize(type);
            
            // Route to appropriate fetcher method based on type
            switch (true) {
                case type === 'posts':
                    response = await this.apiFetcher.fetchPosts(cache.currentPage, cache.params.category);
                    break;
                case type === 'users':
                    response = await this.apiFetcher.fetchUsers(cache.currentPage);
                    break;
                case type.startsWith('messages'):
                    response = await this.apiFetcher.fetchMessages(cache.currentPage, cache.params.user_id);
                    break;
                default:
                    response = await this.apiFetcher.fetchPage(cache.endpoint, cache.currentPage, cache.params);
            }
            
            const newItems = response.items || [];
            const totalCount = response.total;
            
            // Update cache
            this.cacheManager.updateCache(type, newItems, totalCount);
            
            // Update hasMore based on total items loaded vs total available
            const totalLoadedAfterUpdate = cache.items.length;
            cache.hasMore = totalLoadedAfterUpdate < totalCount;
            
            console.log(`Loaded more for ${type}:`, {
                newItemsCount: newItems.length,
                totalLoadedNow: totalLoadedAfterUpdate,
                totalAvailable: totalCount,
                hasMore: cache.hasMore,
                currentPage: cache.currentPage
            });
            
            // Trim cache if too large
            if (cache.items.length > this.config.maxCacheSize) {
                cache.items = cache.items.slice(-this.config.maxCacheSize);
                this.cacheManager.saveToLocalStorage(type, cache);
            }
            
            // For users, return the full response to preserve metadata, otherwise just items
            return type === 'users' ? response : newItems;
        } catch (error) {
            console.error(`Error loading more ${type}:`, error);
            return [];
        } finally {
            this.setLoading(type, false);
        }
    }

    // Get cached items for display
    getCachedItems(type, count = null) {
        return this.cacheManager.getCachedItems(type, count);
    }

    // Clear pagination data for a type
    clearPagination(type) {
        this.cacheManager.clearCache(type);
        this.scrollManager.removeScrollListener(type);
        this.loadingStates.delete(type);
    }

    // Check if currently loading
    isLoading(type) {
        return this.loadingStates.get(type) || false;
    }

    // Set loading state
    setLoading(type, loading) {
        this.loadingStates.set(type, loading);
    }

    // Get pagination stats
    getStats(type) {
        const cache = this.cacheManager.getCache(type);
        if (!cache) return null;
        
        return {
            itemCount: cache.items.length,
            currentPage: cache.currentPage,
            totalCount: cache.totalCount,
            hasMore: cache.hasMore,
            isLoading: this.isLoading(type),
            lastFetchTime: cache.lastFetchTime,
            pageSize: this.getPageSize(type)
        };
    }

    // Update configuration
    updateConfig(newConfig) {
        this.config = { ...this.config, ...newConfig };
    }

    // Get all configured page sizes
    getAllPageSizes() {
        return { ...this.pageSizes };
    }

    // Batch update page sizes
    updatePageSizes(newPageSizes) {
        this.pageSizes = { ...this.pageSizes, ...newPageSizes };
        
        // Update API fetcher for all modified types
        Object.entries(newPageSizes).forEach(([type, pageSize]) => {
            this.apiFetcher.setPageSizeForType(type, pageSize);
        });
    }
}

export default PaginationManager;
