// DOM Utilities - Only genuinely useful DOM manipulation helpers
export class DOMUtilities {
    
    // Toggle element visibility - useful helper
    static toggleVisibility(element, show = null) {
        if (!element) return;
        
        if (show === null) {
            element.style.display = element.style.display === 'none' ? 'block' : 'none';
        } else {
            element.style.display = show ? 'block' : 'none';
        }
    }
       // Add HTML to element with position option
    static insertHTML(element, html, position = 'beforeend') {
        if (!element) return;
        
        element.insertAdjacentHTML(position, html);
    }
    
    // Clear element content
    static clearContent(element) {
        if (!element) return;
        
        element.innerHTML = '';
    }
    
    // Focus element safely
    static focusElement(element) {
        if (!element) return;
        
        try {
            element.focus();
        } catch (error) {
            console.error('Error focusing element:', error);
        }
    }
    
    // Remove event listeners by cloning (prevents duplicates) - useful helper
    static removeEventListeners(element) {
        if (!element) return null;
        
        const newElement = element.cloneNode(true);
        element.parentNode.replaceChild(newElement, element);
        return newElement;
    }
}
