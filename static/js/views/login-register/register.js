import AbstractView from "../abstract-view.js";
import { RegisterController } from "../../controllers/register-controller.js";

export default class RegisterView extends AbstractView {
    constructor() {
        super();
        this.setTitle("Register");
    }

    async getHtml() {
        return `
        <div class="form-container register">
            <h2 class="form-title">Create Account</h2>
            
            <div id="error-messages" class="general-error"></div>
            
            <form id="register-form">
                <div class="form-group">
                    <label for="first-name" class="form-label">First Name</label>
                    <input type="text" id="first-name" name="first-name" class="form-input" placeholder="Enter your first name">
                    <span id="first-name-error" class="form-error"></span>
                </div>
                
                <div class="form-group">
                    <label for="last-name" class="form-label">Last Name</label>
                    <input type="text" id="last-name" name="last-name" class="form-input" placeholder="Enter your last name">
                    <span id="last-name-error" class="form-error"></span>
                </div>
                
                <div class="form-group">
                    <label for="username" class="form-label">Username</label>
                    <input type="text" id="username" name="username" class="form-input" placeholder="Choose a username">
                    <span id="username-error" class="form-error"></span>
                </div>
                
                <div class="form-group">
                    <label for="email" class="form-label">Email</label>
                    <input type="email" id="email" name="email" class="form-input" placeholder="Enter your email">
                    <span id="email-error" class="form-error"></span>
                </div>
                
                <div class="form-group">
                    <label for="password" class="form-label">Password</label>
                    <input type="password" id="password" name="password" class="form-input" placeholder="Create a password">
                    <span id="password-error" class="form-error"></span>
                </div>
                
                <div class="form-group">
                    <label for="age" class="form-label">Age</label>
                    <input type="number" id="age" name="age" class="form-input" placeholder="Enter your age">
                    <span id="age-error" class="form-error"></span>
                </div>
                
                <div class="form-group">
                    <label for="sex" class="form-label">Gender</label>
                    <select id="sex" name="sex" class="form-input">
                        <option value="">Select your gender</option>
                        <option value="0">Male</option>
                        <option value="1">Female</option>
                    </select>
                    <span id="sex-error" class="form-error"></span>
                </div>
                
                <button type="submit" class="form-button full-width">Register</button>
            </form>
            
            <p class="form-footer">
                Already have an account? <a href="/" data-link>Login here</a>
            </p>
        </div>
        `;
    }

    async getData() {
        // Initialize the form using the controller
        RegisterController.initializeForm();
    }
}
