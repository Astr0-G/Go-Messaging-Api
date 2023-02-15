document.addEventListener("DOMContentLoaded", () => {
  // Create a WebSocket connection to the chat room
  const socket = new WebSocket("ws://localhost:8080/chat");

  // Handle incoming messages from the server
  socket.onmessage = (event) => {
    const chatRoom = JSON.parse(event.data);
    const messages = chatRoom.messages;

    // Update the chat box with the new messages
    const chatBox = document.getElementById("chat-box");
    chatBox.innerHTML = "";
    messages.forEach((message) => {
      const messageElement = document.createElement("div");
      messageElement.textContent = message.username + ": " + message.text;
      chatBox.appendChild(messageElement);
    });
  };

  // Send a message to the server when the user clicks the send button
  const sendButton = document.getElementById("send-button");
  sendButton.addEventListener("click", () => {
    const nameInput = document.getElementById("name-input");
    const messageInput = document.getElementById("message-input");

    const message = {
      username: nameInput.value,
      text: messageInput.value,
    };

    socket.send(JSON.stringify(message));

    // Clear the message input field
    messageInput.value = "";
  });
});
