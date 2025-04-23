import AbstractView from "./AbstarctView";

export default class LoginView extends AbstractView {
    constructor() {
        super();
        this.setTitle("Login");
    }

    async getHtml() {
        return `
        <div style="border: 3px solid black;margin-top: 30px;">
            <form id="login-form">
                <label for="username-email">username/email:</label>
                <input type="text" id="username-email" name="username-email"><br>
                <label for="password">password:</label>
                <input type="password" id="password" name="password"><br><br>
                <input type="submit" value="Login">
            </form>
        </div>
        `;
    }

    async getData() {
        let formData = {
            "username-email": document.getElementById("username-email").value,
            "password": document.getElementById("password").value
        }
        const response = await fetch("/login", {
            method: 'post',
            body: JSON.stringify(formData),
            mode: 'cors',
        });
        const data = await response.json();
        return data;
        // const response = await fetch("/login");
        // const data = await response.json();
        // return data;
        // let formData = {
        //     "username-email": document.getElementById("username-email").value,
        //     "password": document.getElementById("password").value
        // }
        // // Send the request
        // fetch("login", {
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