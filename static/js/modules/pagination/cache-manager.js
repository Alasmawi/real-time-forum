// Cache Manager - Handles localStorage caching and cache validation
class CacheManager {
    constructor() {
        this.caches = new Map();
        this.cacheTimeout = 5 * 60 * 1000; // 5 minutes
    }

    // Initialize a cache for a specific type
    initializeCache(type, endpoint, params = {}) {
        const cache = {
            type,
            endpoint,
            params,
            items: [],
            currentPage: 0,
            totalCount: null,
            isLoading: false,
            hasMore: true,
            lastFetchTime: null
        };
        
        this.caches.set(type, cache);
        this.saveToLocalStorage(type, cache);
        
        return cache;
    }

    // Get cache by type
    getCache(type) {
        return this.caches.get(type);
    }

    // Update cache with new items
    updateCache(type, newItems, totalCount = null) {
        const cache = this.caches.get(type);
        if (!cache) return;

        cache.items.push(...newItems);
        cache.currentPage += 1;
        cache.totalCount = totalCount || cache.totalCount;
        cache.hasMore = newItems.length > 0;
        cache.lastFetchTime = Date.now();

        this.saveToLocalStorage(type, cache);
    }

    // Save cache to localStorage
    saveToLocalStorage(type, cache) {
        try {
            const cacheData = {
                ...cache,
                timestamp: Date.now()
            };
            localStorage.setItem(`cache_${type}`, JSON.stringify(cacheData));
        } catch (error) {
            console.warn(`Failed to save cache to localStorage:`, error);
        }
    }

    // Load cache from localStorage
    loadFromLocalStorage(type) {
        try {
            const cached = localStorage.getItem(`cache_${type}`);
            if (!cached) return null;

            const cacheData = JSON.parse(cached);
            if (this.isCacheValid(cacheData)) {
                this.caches.set(type, cacheData);
                return cacheData;
            }
        } catch (error) {
            console.warn(`Failed to load cache from localStorage:`, error);
        }
        return null;
    }

    // Check if cached data is still valid
    isCacheValid(cacheData) {
        if (!cacheData || !cacheData.timestamp) return false;
        const cacheAge = Date.now() - cacheData.timestamp;
        return cacheAge < this.cacheTimeout;
    }

    // Get cached items
    getCachedItems(type, count = null) {
        const cache = this.caches.get(type);
        if (!cache) return [];
        
        return count ? cache.items.slice(0, count) : cache.items;
    }

    // Clear cache
    clearCache(type) {
        this.caches.delete(type);
        localStorage.removeItem(`cache_${type}`);
    }
}

export default CacheManager;
