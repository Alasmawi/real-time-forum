/**
 * Centralized error handler for API responses
 * Maps HTTP status codes from errors.go functions to appropriate actions
 */

class ErrorHandler {
    static storedError = null;
    /**
     * Handle API response errors based on status codes from errors.go
     * @param {Response} response - Fetch API response object
     * @returns {Promise<boolean>} - Returns true if error was handled, false if response was ok
     */
    static async handleResponse(response) {
        if (response.ok) {
            return false; // No error to handle
        }

        console.error('API Error:', {
            status: response.status,
            statusText: response.statusText,
            url: response.url
        });

        try {
            const errorData = await response.json();
            await this.handleErrorByStatus(response.status, errorData);
        } catch (parseError) {
            console.error('Failed to parse error response:', parseError);
            this.showGenericError();
        }

        return true; // Error was handled
    }

    /**
     * Handle errors based on HTTP status code from errors.go functions
     * @param {number} status - HTTP status code
     * @param {Object} errorData - Parsed error response data
     */
    static async handleErrorByStatus(status, errorData) {
        const errorMessage = errorData.Error || 'Unknown error occurred';

        switch (status) {
            case 400: // badRequest from errors.go
                this.handleBadRequest(errorMessage);
                break;

            case 401: // invalidAuthenticationToken or authenticationRequired from errors.go
                this.handleAuthError(errorMessage);
                break;

            case 404: // notFound from errors.go
                this.handleNotFound(errorMessage);
                break;

            case 405: // methodNotAllowed from errors.go
                this.handleMethodNotAllowed(errorMessage);
                break;

            case 422: // failedValidation from errors.go
                this.handleValidationError(errorData);
                break;

            case 500: // serverError from errors.go
                this.handleServerError(errorMessage);
                break;

            default:
                this.handleUnknownError(status, errorMessage);
        }
    }

    /**
     * Handle 400 Bad Request errors
     */
    static handleBadRequest(message) {
        console.error('Bad Request:', message);
        this.storeErrorAndNavigate(400, message);
    }

    /**
     * Handle 401 Authentication errors (invalidAuthenticationToken, authenticationRequired)
     */
    static handleAuthError(message) {
        console.error('Authentication error:', message);
        
        // Clear any stored session data
        this.clearSession();
        
        // Redirect to login page
        window.location.href = '/login';
    }

    /**
     * Handle 404 Not Found errors
     */
    static handleNotFound(message) {
        console.error('Resource not found:', message);
        this.storeErrorAndNavigate(404, message);
    }

    /**
     * Handle 405 Method Not Allowed errors
     */
    static handleMethodNotAllowed(message) {
        console.error('Method not allowed:', message);
        this.storeErrorAndNavigate(405, message);
    }

    /**
     * Handle 422 Validation errors (failedValidation from errors.go)
     */
    static handleValidationError(errorData) {
        console.error('Validation error:', errorData);
        
        // Use the exact error message from JSON response
        const message = errorData.Error || 'Validation failed';
        this.storeErrorAndNavigate(422, message);
    }

    /**
     * Handle 500 Server errors
     */
    static handleServerError(message) {
        console.error('Server error:', message);
        this.storeErrorAndNavigate(500, message);
    }

    /**
     * Handle unknown error status codes
     */
    static handleUnknownError(status, message) {
        console.error(`Unknown error (${status}):`, message);
        this.storeErrorAndNavigate(status, message);
    }

    /**
     * Show a generic error message
     */
    static showGenericError() {
        console.error('Generic error occurred');
        this.storeErrorAndNavigate(500, 'An error occurred. Please try again.');
    }

    /**
     * Store error data and navigate to error page
     */
    static storeErrorAndNavigate(code, message) {
        this.storedError = { code, message };
        history.pushState(null, null, `/error/${code}`);
        window.dispatchEvent(new PopStateEvent('popstate'));
    }

    /**
     * Get stored error data for the error view
     */
    static getStoredError() {
        return this.storedError;
    }

    /**
     * Clear session data on authentication errors
     */
    static clearSession() {
        // Clear any stored authentication tokens or session data
        localStorage.removeItem('authToken');
        localStorage.removeItem('sessionId');
        sessionStorage.clear();
        
        console.log('Session data cleared due to authentication error');
    }

    /**
     * Log successful responses for debugging
     */
    static logSuccess(response) {
        console.log('API Success:', {
            status: response.status,
            statusText: response.statusText,
            url: response.url
        });
    }
}

export default ErrorHandler;
