import AbstractView from "./abstract-view.js";

export default class LoginView extends AbstractView {
    constructor() {
        super();
        this.setTitle("Login");
    }

    async getHtml() {
        return `
        <div style="border: 3px solid black;margin-top: 30px;">
            <form id="login-form">
                <label for="identifier">Email or Username:</label>
                <input type="text" id="identifier" name="identifier"><br>
                <label for="password">Password:</label>
                <input type="password" id="password" name="password"><br><br>
                <input type="submit" value="Login">
            </form>
        </div>
        `;
    }

    async getData() {
        document.getElementById("login-form").addEventListener("submit", async (e) => {
            e.preventDefault();
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
                } else {
                    const data = await response.json();
                    console.log("Login successful:", data);
                }
            } catch (error) {
                console.error("Error:", error);
                console.log("An error occurred during login.");
            }
        });
    }
}
