# Chat App with PocketBase (Golang & Node.js)

![Chat App Logo](/path/to/logo.png)

## Description

The Whisker Chat App with PocketBase is a Websocket based messaging application built using Golang and SvelteKit. It allows users to connect and exchange messages in real-time, leveraging the power of PocketBase as the backend database. The app offers a simple and intuitive interface, supporting both individual and group chats for smooth and seamless communication.

## Features

- Real-time messaging: Instantly send and receive messages in real-time.
- Individual and group chats: Engage in groups to chat with multiple users.
- User-friendly interface: The app's intuitive design makes it easy for users of all levels to navigate and use.
- Online status: Know when other users are online and available for chat.
- Image sharing: Share images with your friends to make conversations more engaging.
- Message history: Access your chat history and previous conversations (up to 256 messages).

## Installation

1. Clone the repository: `git clone https://github.com/yourusername/chat-app.git`
2. Navigate to the project directory: `cd chat-app`
3. Build the Docker image:
   ```bash
   docker build -t chat-app .
   ```
4. Run the Docker container:
   ```bash
   docker run -p 8080:8080 chat-app
   ```

## Configuration

### Environment Variables

The Chat App can be customized using the following environment variables:

- `PUBLIC_POCKETBASE_URL`: Set the public URL for PocketBase. By default, it uses `https://whisker.fly.dev`.
- `PUBLIC_WEBSOCKET_URL`: Set the public WebSocket URL for real-time communication. By default, it uses `wss://whisker.fly.dev/ws`.

## Technologies Used

- Golang: Programming language used to build the backend server for the Chat App.
- Node.js and npm: For building and serving the frontend of the web application.
- PocketBase: Backend database for storing and managing chat data in real-time.

## Contributing

Contributions are welcome! If you find any issues or have suggestions for improvement, feel free to submit a pull request.

## License

This project is licensed under the [MIT License](LICENSE).

## Contact

For any inquiries or feedback, please email us at [contact@example.com](mailto:contact@example.com) or join our community Slack channel.

---

Please make sure to replace `/path/to/logo.png` with the actual path to your Chat App's logo and customize the template with specific details about your Golang and Node.js-based chat application, including the build and run instructions and any other relevant information. This version of the README includes information about PocketBase as the backend database for the Chat App.
