import { CharacterCounter } from '../character-counter.js';
import { FormErrorHelpers } from '../errors/form-error-helpers.js';

export class ReplyBoxManager {
    static characterCounterCleanups = new Map();
    
    static show(postId) {
        const replyBox = document.getElementById(`reply-box-${postId}`);
        if (!replyBox) return;
        
        replyBox.style.display = 'block';
        
        const textarea = replyBox.querySelector('.reply-textarea');
        if (textarea) {
            textarea.focus();
            this.initializeCharacterCounter(postId, textarea, replyBox);
        }
    }
    
    static hide(postId) {
        const replyBox = document.getElementById(`reply-box-${postId}`);
        if (!replyBox) return;
        
        replyBox.style.display = 'none';
        
        // Clear textarea
        const textarea = replyBox.querySelector('.reply-textarea');
        if (textarea) {
            textarea.value = '';
        }
        
        // Clean up character counter
        this.cleanupCharacterCounter(postId, replyBox);
        
        // Clear errors
        FormErrorHelpers.clearAllErrors();
    }
    
    static initializeCharacterCounter(postId, textarea, replyBox) {
        const counterElement = replyBox.querySelector('.character-counter');
        if (!counterElement) return;
        
        // Clean up existing counter
        const existingCleanup = this.characterCounterCleanups.get(postId);
        if (existingCleanup) {
            existingCleanup();
        }
        
        // Initialize new character counter
        const cleanup = CharacterCounter.initialize(textarea, counterElement, {
            maxLength: 500,
            warningThreshold: 400,
            dangerThreshold: 450
        });
        
        if (cleanup) {
            this.characterCounterCleanups.set(postId, cleanup);
        }
    }
    
    static cleanupCharacterCounter(postId, replyBox) {
        const cleanup = this.characterCounterCleanups.get(postId);
        if (cleanup) {
            cleanup();
            this.characterCounterCleanups.delete(postId);
        }
        
        // Reset counter display
        const counterElement = replyBox.querySelector('.character-counter');
        if (counterElement) {
            CharacterCounter.reset(counterElement, { maxLength: 500 });
        }
    }
}
