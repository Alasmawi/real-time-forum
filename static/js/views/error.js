import AbstractView from "./abstract-view.js";

export default class ErrorView extends AbstractView {
    constructor(params) {
        super();
        this.setTitle("Error");
        this.params = params;
        this.errorData = null;
    }

    async getHtml() {
        return `
        <div class="error-page-container">
            <div class="error-page-content">
                <h1 id="error-code"></h1>
                <h2 id="error-message"></h2>
                <div class="error-page-actions">
                    <a href="/" class="error-page-button" data-link>Go Back</a>
                </div>
            </div>
        </div>
        `;
    }

    async getData() {
        // Get error data from ErrorHandler
        this.errorData = window.ErrorHandler?.getStoredError();
        
        if (this.errorData) {
            document.getElementById("error-code").textContent = this.errorData.code;
            document.getElementById("error-message").textContent = this.errorData.message;
        } else {
            history.pushState(null, null, '/home');
            window.dispatchEvent(new PopStateEvent('popstate'));
        }
    }
}
