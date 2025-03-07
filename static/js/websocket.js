
//Javascript that is used to Connect to Websocket and Handle New messages

  // selectedchat is by default General.
  var selectedchat = "general";

  /**
   * changeChatRoom will update the value of selectedchat
   * and also notify the server that it changes chatroom
   * */
  function changeChatRoom() {
      // Change Header to reflect the Changed chatroom
      var newchat = document.getElementById("chatroom");
      if (newchat != null && newchat.value != selectedchat) {
          console.log(newchat);
      }
      return false;
  }
  /**
   * sendMessage will send a new message onto the Websocket
   * */
  function sendMessage() {
      var newmessage = document.getElementById("message");
      if (newmessage != null) {
          console.log(newmessage);
      }
      return false;
  }
  /**
   * Once the website loads, we want to apply listeners and connect to websocket
   * */
  window.onload = function () {
      // Apply our listener functions to the submit event on both forms
      // we do it this way to avoid redirects
      document.getElementById("chatroom-selection").onsubmit = changeChatRoom;
      document.getElementById("chatroom-message").onsubmit = sendMessage;

      // Check if the browser supports WebSocket
      if (window["WebSocket"]) {
          console.log("supports websockets");
      } else {
          alert("Not supporting websockets");
      }
  };