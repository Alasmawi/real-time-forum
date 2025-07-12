import AbstractView from "./abstract-view.js";

export default class LogoutView extends AbstractView {
    constructor(params) {
        super();
        this.setTitle("Logout");
        this.params = params;
    }

    async getHtml() {
        return `
        <div class="logout-container">
            <div id="logout-status">
                <div class="loading">
                    <h2>Logging you out...</h2>
                    <p>Please wait...</p>
                </div>
            </div>
        </div>
        `;
    }

    async getData() {
        try {
            const response = await fetch('/protected/v1/logout', {
                method: 'POST',
                credentials: 'include'
            });
            
            const statusDiv = document.getElementById('logout-status');
            
            // Let ErrorHandler process the response
            const errorHandled = await window.ErrorHandler.handleResponse(response);
            if (errorHandled) {
                // ErrorHandler took care of error (redirect, show error page, etc.)
                return;
            }
            
            // Logout successful
            statusDiv.innerHTML = `
                <div class="logout-success">
                    <h2>✓ Logged out successfully!</h2>
                    <p>You have been logged out of your account.</p>
                    <div class="logout-buttons">
                        <a href="/" data-link class="logout-btn login-btn">Login</a>
                        <a href="/register" data-link class="logout-btn register-btn">Register</a>
                    </div>
                </div>
            `;
        } catch (error) {
            // Handle network/fetch errors
            console.error('Logout error:', error);
            const statusDiv = document.getElementById('logout-status');
            statusDiv.innerHTML = `
                <div class="logout-error">
                    <h2>✗ Network Error</h2>
                    <p>Could not connect to server. Please check your connection and try again.</p>
                    <div class="logout-buttons">
                        <a href="/" data-link class="logout-btn login-btn">Login</a>
                    </div>
                </div>
            `;
        }
    }
}
