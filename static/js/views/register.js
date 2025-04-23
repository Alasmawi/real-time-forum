import AbstractView from "./AbstarctView";

export default class RegisterView extends AbstractView {
    constructor() {
        super();
        this.setTitle("Register");
    }

    async getHtml() {
        return `
        <div style="border: 3px solid black;margin-top: 30px;">
            <form id="register-form">
                <label for="username">username:</label>
                <input type="text" id="username" name="username"><br>
                <label for="password">password:</label>
                <input type="password" id="password" name="password"><br><br>
                <input type="submit" value="Register">
            </form>
        </div>
        `;
    }

    async getData() {
        // const response = await fetch("/register");
        // const data = await response.json();
        // return data;
        let formData = {
            "username": document.getElementById("username").value,
            "email": document.getElementById("email").value,
            "password": document.getElementById("password").value,
            "age": document.getElementById("age").value,
            "sex": document.getElementById("sex").value,
        }
        const response = await fetch("/register", {
            method: 'post',
            body: JSON.stringify(formData),
            mode: 'cors',
        });
        const data = await response.json();
        return data;
        // Send the request
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
    }
}