// Form Handler Utility - Reusable form operations
import { FormErrorHelpers } from './errors/form-error-helpers.js';

export class FormHandler {
    
    // Generic form submission handler
    static async handleFormSubmission(formSelector, submitCallback, options = {}) {
        const form = document.querySelector(formSelector);
        if (!form) return;

        // Add event listener for form submission
        return new Promise((resolve) => {
            form.addEventListener('submit', async (e) => {
                e.preventDefault();
                
                try {
                    const submitBtn = form.querySelector(options.submitBtnSelector || 'button[type="submit"]');
                    
                    // Set loading state
                    this.setLoadingState(submitBtn, true, options.loadingText || 'Processing...');
                    
                    // Call the submit callback
                    const result = await submitCallback(form);
                    
                    if (result.success) {
                        // Clear form if specified
                        if (options.clearOnSuccess) {
                            this.clearForm(form);
                        }
                        resolve(result);
                    } else {
                        // Handle errors using smart error clearing
                        if (result.validationErrors) {
                            FormErrorHelpers.handleValidationErrors(result.validationErrors);
                        } else {
                            // For general errors, clear field errors but show general error
                            FormErrorHelpers.smartClearErrors({});
                            FormErrorHelpers.displayGeneralError(result.error || 'An error occurred');
                        }
                        resolve(null);
                    }
                    
                } catch (error) {
                    console.error('Form submission error:', error);
                    FormErrorHelpers.displayGeneralError('An error occurred. Please try again.');
                    resolve(null);
                } finally {
                    // Reset loading state
                    const submitBtn = form.querySelector(options.submitBtnSelector || 'button[type="submit"]');
                    this.setLoadingState(submitBtn, false, options.defaultText || 'Submit');
                }
            });
        });
    }
    
    // Set loading state for buttons
    static setLoadingState(button, isLoading, text) {
        if (!button) return;
        
        button.disabled = isLoading;
        button.textContent = text;
    }
    
    // Clear form inputs
    static clearForm(form) {
        if (!form) return;
        
        const inputs = form.querySelectorAll('input, textarea, select');
        inputs.forEach(input => {
            if (input.type === 'checkbox' || input.type === 'radio') {
                input.checked = false;
            } else {
                input.value = '';
            }
        });
        
        // Remove selected classes
        form.querySelectorAll('.selected').forEach(el => el.classList.remove('selected'));
    }
}
