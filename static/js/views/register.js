import AbstractView from "./abstract-view.js";

export default class RegisterView extends AbstractView {
    constructor() {
        super();
        this.setTitle("Register");
    }

    async getHtml() {
        return `
        <div style="border: 3px solid black;margin-top: 30px;">
            <form id="register-form">
                <label for="first-name">First Name:</label>
                <input type="text" id="first-name" name="first-name"><br>
                <label for="last-name">Last Name:</label>
                <input type="text" id="last-name" name="last-name"><br>
                <label for="email">Email:</label>
                <input type="text" id="email" name="email"><br>
                <label for="username">Username:</label>
                <input type="text" id="username" name="username"><br>
                <label for="password">Password:</label>
                <input type="password" id="password" name="password"><br>
                <label for="age">Age:</label>
                <input type="number" id="age" name="age" min="12" max="90" step="1">
                <label for="sex">Sex:</label>
                <select id="sex" name="sex" required>
                    <option value="">Select Gender</option>
                    <option value="male">Male</option>
                    <option value="female">Female</option>
                </select><br>
                <input type="submit" value="Register">
            </form>
        </div>
        `;
    }

    // form on submit!
    async getData() {
        // Set up registration form submission
        document.getElementById("register-form").addEventListener("submit", async (e) => {
            e.preventDefault();

            let sexValue;
            const sexSelection = document.getElementById("sex").value;
            if (sexSelection === "male") {
                sexValue = true;  
            } else if (sexSelection === "female") {
                sexValue = false; 
            } else {
                alert("Please select your gender");
                return;
            }

            const ageValue = document.getElementById("age").value;
            const age = ageValue ? parseInt(ageValue) : 0;
            
            // Validate age before sending
            if (!ageValue || isNaN(age) || age < 1 || age > 120) {
                alert("Please enter a valid age between 1 and 120");
                return;
            }

            let formData = {
                "f_name": document.getElementById("first-name").value.trim(),
                "l_name": document.getElementById("last-name").value.trim(),
                "email": document.getElementById("email").value.trim(),
                "username": document.getElementById("username").value.trim(),
                "password": document.getElementById("password").value,
                "age": age,
                "sex": sexValue,
            };

            try {
                const response = await fetch("/v1/register", {
                    method: 'POST',
                    headers: {
                        "Content-Type": "application/json",
                    },
                    body: JSON.stringify(formData),
                });

                if (!response.ok) {
                    const errorData = await response.json();
                    console.error("Validation errors:", errorData);
                    
                    // Display validation errors to the user
                    if (errorData.FieldErrors) {
                        let errorMessage = "Registration failed:\n\n";
                        for (const [field, message] of Object.entries(errorData.FieldErrors)) {
                            errorMessage += `• ${message}\n`;
                        }
                        alert(errorMessage);
                    } else {
                        alert("Registration failed. Please check your input and try again.");
                    }
                } else {
                    console.log("Registration successful");
                    alert("Registration successful! Redirecting to login page...");
                    // Redirect to login page after successful registration
                    window.location.href = "/login";
                }
            } catch (error) {
                console.error("Error:", error);
                alert("An error occurred during registration. Please try again.");
            }
        });
    }
}

// async getData() {
//     const response = await fetch("/register");
//     const data = await response.json();
//     return data;

// fetch("register", {
//     method: 'post',
//     body: JSON.stringify(formData),
//     mode: 'cors',
// }).then((response) => {
//     if (response.ok) {
//         return response.json();
//     } else {
//         throw 'unauthorized';
//     }
// })
// return false;
// }

// const form = document.getElementById("register-form");

// form.addEventListener("submit", async (e) => {
//     e.preventDefault();

//     let formData = {
//         "username": document.getElementById("username").value,
//         "email": document.getElementById("email").value,
//         "password": document.getElementById("password").value,
//         "age": document.getElementById("age").value,
//         "sex": document.getElementById("sex").value,
//     };

//     try {
//         const response = await fetch("/register", {
//             method: 'POST',
//             headers: {
//                 "Content-Type": "application/json",
//             },
//             body: JSON.stringify(formData),
//         });

//         console.log("Form Data:", formData);
//         console.log("Response:", response);

//         if (!response.ok) {
//             const errorData = await response.json();
//             console.error("Error:", errorData);
//             alert("Registration failed: " + errorData.message);
//         } else {
//             const data = await response.json();
//             console.log("Registration successful:", data);
//             alert("Registration successful!");
//         }
//     } catch (error) {
//         console.error("Error:", error);
//         alert("An error occurred during registration.");
//     }
// });
