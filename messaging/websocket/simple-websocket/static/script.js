let lastSentMessage = ''; // Track the last sent message

document.addEventListener('DOMContentLoaded', () => {
    const ws = new WebSocket('ws://localhost:8080/ws');
    const messagesDiv = document.getElementById('messages');
    const messageInput = document.getElementById('message-input');
    const sendBtn = document.getElementById('send-btn');

    // Function to append a new message to the chat
    function appendMessage(message) {
        const messageElement = document.createElement('div');
        messageElement.textContent = message;
        messagesDiv.appendChild(messageElement);
        messagesDiv.scrollTop = messagesDiv.scrollHeight; // Auto-scroll to the bottom
    }

    // Listen for messages from the server
    ws.onmessage = (event) => {
        const message = event.data;

        // Check if the received message is the same as the last sent message
        if (message === lastSentMessage) {
            return; // Skip displaying the message again
        }

        appendMessage(message);
    };

    // Handle connection errors
    ws.onerror = (error) => {
        console.error('WebSocket error:', error);
        appendMessage('Error connecting to the server.');
    };

    // Handle connection close
    ws.onclose = () => {
        appendMessage('Disconnected from the server.');
    };

    // Send a message when the "Send" button is clicked
    sendBtn.addEventListener('click', () => {
        const message = messageInput.value.trim();
        if (message) {
            lastSentMessage = message; // Store the last sent message
            ws.send(message); // Send the message to the server
            appendMessage(`You: ${message}`); // Display the sent message locally
            messageInput.value = ''; // Clear the input field
        }
    });

    // Allow sending messages by pressing the "Enter" key
    messageInput.addEventListener('keypress', (event) => {
        if (event.key === 'Enter') {
            sendBtn.click();
        }
    });
});