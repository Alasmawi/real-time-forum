// New Post Card View - Handles HTML generation for inline post creation
import { AvatarGenerator } from '../../modules/utils/UI-UX/avatar-generator.js';

export class NewPostCardView {
    
    // Generate HTML for the new post creation card
    static async generateHTML(user, categories) {
        const avatarHTML = AvatarGenerator.generateAvatarHTML(user.username, { size: 'medium' });
        
        return `
            <div class="new-post-card">
                <div id="new-post-error-messages" class="general-error"></div>
                
                <div class="post-header">
                    <div class="post-avatar">
                        ${avatarHTML}
                    </div>
                    <div class="post-username">${user.username}</div>
                </div>
                
                <form id="new-post-form">
                    <div class="form-group">
                        <textarea 
                            id="post-content" 
                            name="content" 
                            class="form-input post-textarea" 
                            placeholder="What's on your mind?" 
                            maxlength="500"
                        ></textarea>
                        <span id="post-content-error" class="form-error"></span>
                        <div class="character-counter">
                            <span id="char-count">0/500 characters</span>
                        </div>
                    </div>
                    
                    <div class="form-group expandable-content">
                        <label class="form-label">Categories (optional):</label>
                        <div class="category-selection">
                            ${categories.map(cat => `
                                <button type="button" class="category-btn" data-category-id="${cat.id}">
                                    ${cat.name}
                                </button>
                            `).join('')}
                        </div>
                        <span id="categories-error" class="form-error"></span>
                    </div>
                    
                    <div class="form-group expandable-content">
                        <button type="submit" class="form-button post-submit-btn">Post</button>
                    </div>
                </form>
            </div>
        `;
    }

    // Render the new post card to the DOM
    static async render(user, categories, container) {
        if (!container) return false;
        
        const html = await this.generateHTML(user, categories);
        container.insertAdjacentHTML('afterbegin', html);
        return true;
    }

    // Check if new post card already exists
    static exists(container) {
        return container?.querySelector('.new-post-card') !== null;
    }
}
