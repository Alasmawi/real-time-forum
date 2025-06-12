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
        <div>
            <h1 id="error-code"></h1>
            <h2 id="error-message"></h2>
            <button onclick="window.history.back()">Go Back</button>
            <button onclick="window.location.href = '/'">Home</button>
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
            document.getElementById("error-code").textContent = "500";
            document.getElementById("error-message").textContent = "Unknown error occurred";
        }
    }
}
