// Login Controller - Handles login form logic
import { FormErrorHelpers } from "../modules/utils/errors/form-error-helpers.js";
import { UserModel } from "../models/user-model.js";
import { FormHandler } from "../modules/utils/form-handler.js";

export class LoginController {
    
    static async handleFormSubmission(formData) {
        try {
            // Use UserModel for login (now returns Response object or null)
            const response = await UserModel.login(formData);

            console.log("Form Data:", formData);
            console.log("Response:", response);

            if (response === null) {
                // ErrorHandler already handled the error (redirect, error page, etc.)
                return { success: false, type: 'general', message: "Login failed. Please try again." };
            }

            if (!response.ok) {
                const errorData = await response.json();
                console.error("Error:", errorData);
                
                if (response.status === 422) {
                    // Handle validation errors
                    return { success: false, type: 'validation', data: errorData };
                } else {
                    // Handle other errors
                    return { success: false, type: 'general', message: "An error occurred during login. Please try again." };
                }
            } else {
                console.log("Login successful");
                return { success: true };
            }
        } catch (error) {
            console.error("Error:", error);
            return { success: false, type: 'general', message: "An error occurred during login. Please try again." };
        }
    }

    static collectFormData() {
        return {
            "identifier": document.getElementById("identifier").value,
            "password": document.getElementById("password").value,
        };
    }

    static initializeForm() {
        FormHandler.handleFormSubmission(
            '#login-form',
            async (form) => {
                const formData = this.collectFormData();
                const result = await this.handleFormSubmission(formData);
                
                if (result.success) {
                    // Redirect on success
                    window.location.href = "/home";
                    return { success: true };
                } else if (result.type === 'validation') {
                    return { success: false, validationErrors: result.data };
                } else {
                    return { success: false, error: result.message };
                }
            },
            {
                loadingText: 'Logging in...',
                defaultText: 'Login',
                clearOnSuccess: true
            }
        );
    }
}
