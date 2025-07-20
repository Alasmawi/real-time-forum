// Character Counter Utility - Reusable character counting functionality
export class CharacterCounter {
    
    /**
     * Initialize character counter for a textarea element
     * @param {HTMLElement} textarea - The textarea element to monitor
     * @param {HTMLElement} counterElement - The element to display the count
     * @param {Object} options - Configuration options
     * @param {number} options.maxLength - Maximum character limit (default: 500)
     * @param {number} options.warningThreshold - When to show orange warning (default: 400)
     * @param {number} options.dangerThreshold - When to show red warning (default: 450)
     * @param {string} options.defaultColor - Default counter color (default: '#666')
     * @param {string} options.warningColor - Warning counter color (default: '#f39c12')
     * @param {string} options.dangerColor - Danger counter color (default: '#e74c3c')
     */
    static initialize(textarea, counterElement, options = {}) {
        if (!textarea || !counterElement) {
            console.error('CharacterCounter: textarea and counterElement are required');
            return null;
        }

        const config = {
            maxLength: 500,
            warningThreshold: 400,
            dangerThreshold: 450,
            defaultColor: '#666',
            warningColor: '#f39c12',
            dangerColor: '#e74c3c',
            ...options
        };

        // Set initial counter text
        counterElement.textContent = `0/${config.maxLength} characters`;
        counterElement.style.color = config.defaultColor;

        // Add input event listener
        const updateCounter = (e) => {
            const count = e.target.value.length;
            const trimmedCount = e.target.value.trim().length;
            
            counterElement.textContent = `${count}/${config.maxLength} characters`;
            
            // Update color based on character count and validation
            if (count > config.dangerThreshold) {
                counterElement.style.color = config.dangerColor;
            } else if (count > config.warningThreshold) {
                counterElement.style.color = config.warningColor;
            } else if (trimmedCount === 0 && count > 0) {
                // Show warning if only whitespace characters
                counterElement.style.color = config.warningColor;
            } else {
                counterElement.style.color = config.defaultColor;
            }
        };

        textarea.addEventListener('input', updateCounter);

        // Return cleanup function
        return () => {
            textarea.removeEventListener('input', updateCounter);
        };
    }

    /**
     * Reset counter to initial state
     * @param {HTMLElement} counterElement - The counter element to reset
     * @param {Object} options - Configuration options
     */
    static reset(counterElement, options = {}) {
        if (!counterElement) return;

        const config = {
            maxLength: 500,
            defaultColor: '#666',
            ...options
        };

        counterElement.textContent = `0/${config.maxLength} characters`;
        counterElement.style.color = config.defaultColor;
    }

    /**
     * Get current character count from textarea
     * @param {HTMLElement} textarea - The textarea element
     * @returns {number} Current character count
     */
    static getCount(textarea) {
        return textarea ? textarea.value.length : 0;
    }
}
