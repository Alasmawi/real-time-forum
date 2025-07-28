// Scroll Manager - Handles infinite scroll with throttle and debounce
class ScrollManager {
    constructor() {
        this.scrollListeners = new Map();
        this.config = {
            loadTriggerPercent: 95,        // Load more when 95% down
            reverseLoadTriggerPercent: 5,  // Load older when within 5% of top
            throttleDelay: 300,
            debounceDelay: 500,
        };
    }

    // Set up infinite scroll for a container with configurable throttle/debounce
    setupInfiniteScroll(type, containerId, loadMoreCallback, options = {}) {
        const container = document.getElementById(containerId);
        if (!container) {
            console.error(`Container '${containerId}' not found`);
            return;
        }

        // Remove existing listener if any
        this.removeScrollListener(type);

        // Merge options with global defaults (options override defaults)
        const useThrottle = options.useThrottle !== false; // Default: true
        const useDebounce = options.useDebounce === true;  // Default: false
        const throttleDelay = options.throttleDelay || this.config.throttleDelay;
        const debounceDelay = options.debounceDelay || this.config.debounceDelay;
        const isReverse = options.isReverse === true; // Default: false (normal scroll down)

        // Track scroll direction for reverse scrolling
        let lastScrollTop = container.scrollTop;
        let isLoading = false;

        // Create base scroll handler
        let scrollHandler = async () => {
            if (isLoading) return;
            
            if (isReverse) {
                // Reverse scrolling: load when scrolling up and near top
                const currentScrollTop = container.scrollTop;
                const scrollingUp = currentScrollTop < lastScrollTop;
                
                if (this.shouldLoadReverse(container) && scrollingUp) {
                    isLoading = true;
                    try {
                        await loadMoreCallback();
                    } finally {
                        isLoading = false;
                    }
                }
                
                lastScrollTop = currentScrollTop;
            } else {
                // Normal scrolling: load when near bottom
                if (this.shouldLoadMore(container)) {
                    isLoading = true;
                    try {
                        await loadMoreCallback();
                    } finally {
                        isLoading = false;
                    }
                }
            }
        };

        // Apply debounce first if enabled
        if (useDebounce) {
            scrollHandler = this.debounce(scrollHandler, debounceDelay);
        }

        // Apply throttle if enabled (wraps debounce if both used)
        if (useThrottle) {
            scrollHandler = this.throttle(scrollHandler, throttleDelay);
        }

        // Store listener reference for cleanup
        this.scrollListeners.set(type, { container, handler: scrollHandler });
        
        // Add scroll listener
        container.addEventListener('scroll', scrollHandler);
    }

    // Check if should load more content based on scroll position (bottom)
    shouldLoadMore(container) {
        const scrollTop = container.scrollTop;
        const scrollHeight = container.scrollHeight;
        const clientHeight = container.clientHeight;
        
        const scrollPercent = (scrollTop + clientHeight) / scrollHeight * 100;
        return scrollPercent >= this.config.loadTriggerPercent;
    }

    // Check if should load reverse content based on scroll position (top)
    shouldLoadReverse(container) {
        const scrollTop = container.scrollTop;
        const scrollHeight = container.scrollHeight;
        const clientHeight = container.clientHeight;
        
        // Calculate how far from top as percentage (0% = at top, 100% = at bottom)
        const scrollPercent = (scrollTop) / (scrollHeight - clientHeight) * 100;
        
        // Load when within the reverse trigger percentage from top
        return scrollPercent <= this.config.reverseLoadTriggerPercent;
    }

    // Remove scroll listener for cleanup
    removeScrollListener(type) {
        const listener = this.scrollListeners.get(type);
        if (listener) {
            listener.container.removeEventListener('scroll', listener.handler);
            this.scrollListeners.delete(type);
        }
    }

    // Remove all scroll listeners
    removeAllListeners() {
        for (const [type] of this.scrollListeners) {
            this.removeScrollListener(type);
        }
    }

    // Throttle function - limits function calls to at most once per delay period
    throttle(func, delay) {
        let timeoutId;
        let lastExecTime = 0;
        
        return function (...args) {
            const currentTime = Date.now();
            
            if (currentTime - lastExecTime > delay) {
                func.apply(this, args);
                lastExecTime = currentTime;
            } else {
                clearTimeout(timeoutId);
                timeoutId = setTimeout(() => {
                    func.apply(this, args);
                    lastExecTime = Date.now();
                }, delay - (currentTime - lastExecTime));
            }
        };
    }

    // Debounce function - delays function execution until after delay period of inactivity
    debounce(func, delay) {
        let timeoutId;
        
        return function (...args) {
            clearTimeout(timeoutId);
            timeoutId = setTimeout(() => func.apply(this, args), delay);
        };
    }

    // Update configuration
    updateConfig(newConfig) {
        this.config = { ...this.config, ...newConfig };
    }
}

export default ScrollManager;
