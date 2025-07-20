class ErrorHandler {
    static storedError = null;

    /**
     * @param {Response} response - Fetch API response object
     * @returns {Promise<boolean>} - Returns true if error was handled, false if response was ok or 422
     */
    static async handleResponse(response) {
        if (response.ok) {
            return false; // No error, response was successful
        }

        // For 422 validation errors, delegate immediately without consuming the response body
        if (response.status === 422) {
            console.log('422 validation error - delegating to caller');
            return false; // Don't handle, let caller deal with it
        }

        try {
            const errorData = await response.json();
            const message = errorData.message /*|| 'Unknown error occurred'*/;

            console.error(`API Error ${response.status}:`, message);

            // Simple status-based handling
            switch (response.status) {
                case 401:
                    // Redirect to login when unauthorized
                    window.location.href = '/';
                    return true; // Error was handled

                case 404:
                    this.showErrorPage(response.status, message);
                    return true; // Error was handled

                case 409:
                    window.location.href = '/home';
                    return true

                case 500:
                    // Show error page with backend's message
                    this.showErrorPage(response.status, message);
                    return true; // Error was handled

                default:
                    // Log other errors but don't break user flow
                    console.error(`Unhandled error ${response.status}:`, message);
                    return false; // Don't handle, let caller deal with it
            }

        } catch (parseError) {
            console.error('Failed to parse error response:', parseError);
            this.showErrorPage(500, 'An error occurred. Please try again.');
            return true; // Error was handled
        }
    }

    static showErrorPage(code, message) {
        this.storedError = { code, message };
        history.pushState(null, null, `/error/${code}`);
        window.dispatchEvent(new PopStateEvent('popstate'));
    }

    static getStoredError() {
        return this.storedError;
    }

    static clearStoredError() {
        this.storedError = null;
    }
}

export default ErrorHandler;
