// Register Controller - Handles registration form logic
import { FormErrorHelpers } from "../modules/utils/errors/form-error-helpers.js";
import { UserModel } from "../models/user-model.js";
import { FormHandler } from "../modules/utils/form-handler.js";

export class RegisterController {
    
    static async handleFormSubmission(formData) {
        try {
            // Use UserModel for registration (now returns Response object or null)
            const response = await UserModel.register(formData);

            console.log("Form Data:", formData);
            console.log("Response:", response);

            if (response === null) {
                // ErrorHandler already handled the error (redirect, error page, etc.)
                return { success: false, type: 'general', message: "Registration failed. Please try again." };
            }

            if (!response.ok) {
                const errorData = await response.json();
                console.error("Error:", errorData);
                
                if (response.status === 422) {
                    // Handle validation errors
                    return { success: false, type: 'validation', data: errorData };
                } else {
                    // Handle other errors
                    return { success: false, type: 'general', message: "An error occurred during registration. Please try again." };
                }
            } else {
                console.log("Registration successful");
                return { success: true };
            }
        } catch (error) {
            console.error("Error:", error);
            return { success: false, type: 'general', message: "An error occurred during registration. Please try again." };
        }
    }


    static collectFormData() {
        return {
            "first-name": document.getElementById("first-name").value,
            "last-name": document.getElementById("last-name").value,
            "username": document.getElementById("username").value,
            "email": document.getElementById("email").value,
            "password": document.getElementById("password").value,
            "age": document.getElementById("age").value,
            "sex": document.getElementById("sex").value,
        };
    }

    static initializeForm() {
        FormHandler.handleFormSubmission(
            '#register-form',
            async (form) => {
                const formData = this.collectFormData();
                const result = await this.handleFormSubmission(formData);
                
                if (result.success) {
                    // Redirect on success
                    window.location.href = "/";
                    return { success: true };
                } else if (result.type === 'validation') {
                    return { success: false, validationErrors: result.data };
                } else {
                    return { success: false, error: result.message };
                }
            },
            {
                loadingText: 'Registering...',
                defaultText: 'Register',
                clearOnSuccess: true
            }
        );
    }
}
