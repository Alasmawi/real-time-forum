// Form Error Helpers - Reusable utility for smart form error handling
export class FormErrorHelpers {
    
    /**
     * Smart error clearing - only clears errors that are not in the new error response
     * This prevents jarring layout jumps when forms are re-submitted
     * @param {Object} newErrors - New error data from server response
     */
    static smartClearErrors(newErrors = {}) {
        // Only clear general errors if there are no new general errors
        if (!newErrors.Errors || newErrors.Errors.length === 0) {
            const generalErrorElement = document.getElementById("error-messages");
            if (generalErrorElement && generalErrorElement.innerHTML !== "") {
                generalErrorElement.innerHTML = "";
            }
        }
        
        // Clear field errors that are not in the new errors
        const errorElements = document.querySelectorAll('.form-error');
        errorElements.forEach(element => {
            const fieldName = element.id.replace('-error', '');
            if (!newErrors.FieldErrors || !newErrors.FieldErrors[fieldName]) {
                if (element.innerHTML !== "") {
                    element.innerHTML = "";
                }
            }
        });
        
        // Remove error classes from inputs that don't have errors anymore
        const inputElements = document.querySelectorAll('.form-input');
        inputElements.forEach(element => {
            const fieldName = element.id;
            if (!newErrors.FieldErrors || !newErrors.FieldErrors[fieldName]) {
                element.classList.remove('error');
            }
        });
    }

    /**
     * Display field-specific validation errors
     * @param {Object} errorData - Server error response with FieldErrors
     */
    static displayFieldErrors(errorData) {
        if (!errorData.FieldErrors) return;

        for (const [field, message] of Object.entries(errorData.FieldErrors)) {
            const errorElement = document.getElementById(`${field}-error`);
            const inputElement = document.getElementById(field);
            
            if (errorElement && errorElement.innerHTML !== message) {
                errorElement.innerHTML = message;
            }
            
            if (inputElement && !inputElement.classList.contains('error')) {
                inputElement.classList.add('error');
            }
        }
    }

    /**
     * Display general form errors (non-field-specific)
     * @param {Object} errorData - Server error response with Errors array
     */
    static displayGeneralErrors(errorData) {
        if (!errorData.Errors || errorData.Errors.length === 0) return;

        const generalErrors = errorData.Errors.join("<br>");
        const errorMessages = document.getElementById("error-messages");
        if (errorMessages && errorMessages.innerHTML !== generalErrors) {
            errorMessages.innerHTML = generalErrors;
        }
    }

    /**
     * Display a single general error message
     * @param {String} message - Error message to display
     */
    static displayGeneralError(message) {
        const errorMessages = document.getElementById("error-messages");
        if (errorMessages && errorMessages.innerHTML !== message) {
            errorMessages.innerHTML = message;
        }
    }

    /**
     * Handle validation errors from server response
     * Combines field and general error display with smart clearing
     * @param {Object} errorData - Complete server error response
     */
    static handleValidationErrors(errorData) {
        // Smart clear - only clear errors not in new response
        this.smartClearErrors(errorData);
        
        // Display new field errors
        this.displayFieldErrors(errorData);
        
        // Display new general errors
        this.displayGeneralErrors(errorData);
    }

    /**
     * Clear all errors completely (for successful submissions)
     */
    static clearAllErrors() {
        // Clear general errors
        const generalErrorElement = document.getElementById("error-messages");
        if (generalErrorElement) {
            generalErrorElement.innerHTML = "";
        }
        
        // Clear all field errors
        const errorElements = document.querySelectorAll('.form-error');
        errorElements.forEach(element => {
            element.innerHTML = "";
        });
        
        // Remove error classes from all inputs
        const inputElements = document.querySelectorAll('.form-input');
        inputElements.forEach(element => {
            element.classList.remove('error');
        });
    }

    /**
     * Add error class to specific input field
     * @param {String} fieldName - Name of the field to add error class to
     */
    static addInputErrorClass(fieldName) {
        const inputElement = document.getElementById(fieldName);
        if (inputElement) {
            inputElement.classList.add('error');
        }
    }

    /**
     * Remove error class from specific input field
     * @param {String} fieldName - Name of the field to remove error class from
     */
    static removeInputErrorClass(fieldName) {
        const inputElement = document.getElementById(fieldName);
        if (inputElement) {
            inputElement.classList.remove('error');
        }
    }

    /**
     * Get current error state for comparison
     * @returns {Object} Current error state
     */
    static getCurrentErrorState() {
        const state = {
            general: "",
            fields: {}
        };

        // Get general error
        const generalErrorElement = document.getElementById("error-messages");
        if (generalErrorElement) {
            state.general = generalErrorElement.innerHTML;
        }

        // Get field errors
        const errorElements = document.querySelectorAll('.form-error');
        errorElements.forEach(element => {
            const fieldName = element.id.replace('-error', '');
            state.fields[fieldName] = element.innerHTML;
        });

        return state;
    }
}

// Export for global access if needed
window.FormErrorHelpers = FormErrorHelpers;
