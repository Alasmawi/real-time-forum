class ErrorHandler {
    static storedError = null;

    /**
     * @param {Response} response - Fetch API response object
     * @returns {Promise<boolean>} - Returns true if error was handled, false if response was ok
     */
    static async handleResponse(response) {
        if (response.ok) {
            return;
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
                    break;

                case 404:
                    this.showErrorPage(response.status, message);
                    break;
              
                case 422:

                case 500:
                    // Show error page with backend's message
                    this.showErrorPage(response.status, message);
                    break;

                default:
                    // Log other errors but don't break user flow
                    console.error(`Unhandled error ${response.status}:`, message);
            }

        } catch (parseError) {
            console.error('Failed to parse error response:', parseError);
            this.showErrorPage(500, 'An error occurred. Please try again.');
        }

        return;
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
