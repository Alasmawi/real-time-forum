import AbstractView from "../abstract-view.js";
import { LoginController } from "../../controllers/login-controller.js";

export default class LoginView extends AbstractView {
    constructor() {
        super();
        this.setTitle("Login");
    }

    async getHtml() {
        return `
        <div class="form-container">
            <h2 class="form-title">Login</h2>
            
            <div id="error-messages" class="general-error"></div>
            
            <form id="login-form">
                <div class="form-group">
                    <label for="identifier" class="form-label">Email or Username</label>
                    <input type="text" id="identifier" name="identifier" class="form-input" placeholder="Enter your email or username">
                    <span id="identifier-error" class="form-error"></span>
                </div>
                
                <div class="form-group">
                    <label for="password" class="form-label">Password</label>
                    <input type="password" id="password" name="password" class="form-input" placeholder="Enter your password">
                    <span id="password-error" class="form-error"></span>
                </div>
                
                <button type="submit" class="form-button full-width">Login</button>
            </form>
            
            <p class="form-footer">
                Don't have an account? <a href="/register" data-link>Register here</a>
            </p>
        </div>
        `;
    }

    async getData() {
        // Initialize the form using the controller
        LoginController.initializeForm();
    }
}
