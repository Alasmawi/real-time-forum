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
                <input type="text" id="age" name="age"><br>
                <label for="sex">Sex:</label>
                <input type="sex" id="sex" name="sex"><br>
                <input type="submit" value="Register">
            </form>
        </div>
        `;
    }

    // form on submit!
    async getData() {
        document.getElementById("register-form").addEventListener("submit", async (e) => {
            e.preventDefault();
            let formData = {
                "first-name": document.getElementById("first-name").value,
                "last-name": document.getElementById("last-name").value,
                "email": document.getElementById("email").value,
                "username": document.getElementById("username").value,
                "password": document.getElementById("password").value,
                "age": document.getElementById("age").value,
                "sex": document.getElementById("sex").value,
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
                    console.error("Error:", errorData);
                } else {
                    const data = await response.json();
                    console.log("Registration successful:", data);
                }
            } catch (error) {
                console.error("Error:", error);
                alert("An error occurred during registration.");
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
