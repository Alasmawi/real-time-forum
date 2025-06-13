import AbstractView from "./abstract-view.js";

export default class LoginView extends AbstractView {
    constructor() {
        super();
        this.setTitle("Login");
    }

    async getHtml() {
        return `
        <div style="border: 3px solid black;margin-top: 30px;">
            <div id="error-messages" style="color: red; margin-bottom: 10px;"></div>
            <form id="login-form">
                <label for="identifier">Email or Username:</label>
                <input type="text" id="identifier" name="identifier"><br>
                <div id="identifier-error" style="color: red; font-size: 12px;"></div>
                
                <label for="password">Password:</label>
                <input type="password" id="password" name="password"><br>
                <div id="password-error" style="color: red; font-size: 12px;"></div>
                <br>
                <input type="submit" value="Login">
            </form>
        </div>
        `;
    }

    async getData() {
        document.getElementById("login-form").addEventListener("submit", async (e) => {
            e.preventDefault();
            
            // Clear previous errors
            this.clearErrors();
            
            let formData = {
                "identifier": document.getElementById("identifier").value,
                "password": document.getElementById("password").value,
            };

            try {
                const response = await fetch("/v1/login", {
                    method: 'POST',
                    headers: {
                        "Content-Type": "application/json",
                    },
                    body: JSON.stringify(formData),
                });

                console.log("Form Data:", formData);
                console.log("Response:", response);

                if (!response.ok) {
                    const errorData = await response.json();
                    console.error("Error:", errorData);
                    
                    if (response.status === 422) {
                        // Handle validation errors
                        this.displayValidationErrors(errorData);
                    } else {
                        // Handle other errors
                        this.displayGeneralError("An error occurred during login. Please try again.");
                    }
                } else {
                    console.log("Login successful");
                    // Redirect to posts page after successful login
                    window.location.href = "/posts";
                }
            } catch (error) {
                console.error("Error:", error);
                this.displayGeneralError("An error occurred during login. Please try again.");
            }
        });
    }

    clearErrors() {
        document.getElementById("error-messages").innerHTML = "";
        document.getElementById("identifier-error").innerHTML = "";
        document.getElementById("password-error").innerHTML = "";
    }

    displayValidationErrors(errorData) {
        if (errorData.FieldErrors) {
            // Display field-specific errors
            for (const [field, message] of Object.entries(errorData.FieldErrors)) {
                const errorElement = document.getElementById(`${field}-error`);
                if (errorElement) {
                    errorElement.innerHTML = message;
                }
            }
        }
        
        if (errorData.Errors && errorData.Errors.length > 0) {
            // Display general errors
            const generalErrors = errorData.Errors.join("<br>");
            document.getElementById("error-messages").innerHTML = generalErrors;
        }
    }

    displayGeneralError(message) {
        document.getElementById("error-messages").innerHTML = message;
    }
}
