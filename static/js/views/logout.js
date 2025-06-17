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
            
            if (response.ok) {
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
            } else {
                // Logout failed
                statusDiv.innerHTML = `
                    <div class="logout-error">
                        <h2>⚠ Logout Failed</h2>
                        <p>There was an issue logging you out. Please try again.</p>
                        <div class="logout-buttons">
                            <a href="/" data-link class="logout-btn login-btn">Login</a>
                        </div>
                    </div>
                `;
            }
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
