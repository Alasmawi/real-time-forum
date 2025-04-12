  /**
         * login will send a login request to the server and then 
         * connect websocket
         * */
  function login() {
    let formData = {
        "username-email": document.getElementById("username-email").value,
        "password": document.getElementById("password").value
    }
    // Send the request
    fetch("login", {
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