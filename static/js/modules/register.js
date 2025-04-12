  /**
         * login will send a login request to the server and then 
         * connect websocket
         * */
  function register() {
    let formData = {
        "username": document.getElementById("username").value,
        "email": document.getElementById("email").value,
        "password": document.getElementById("password").value,
        "age": document.getElementById("age").value,
        "sex": document.getElementById("sex").value,
    }
    // Send the request
    fetch("register", {
        method: 'post',
        body: JSON.stringify(formData),
        mode: 'cors',
    }).then((response) => {
        if (response.ok) {
            return response.json();
        } else {
            throw 'unauthorized';
        }
    })
    // .then((data) => {
    //     // Now we have a OTP, send a Request to Connect to WebSocket
    //     connectWebsocket(data.otp);
    // }).catch((e) => { alert(e) });
    return false;
}