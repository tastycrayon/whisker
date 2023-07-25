# Chat App with PocketBase (Golang & Svelte)

![Whisker Logo](/client/static/favicon.png)

## Description

The Whisker Chat App is a Websocket based messaging application built using Golang and SvelteKit. It allows users to connect and exchange messages in real-time, leveraging the power of PocketBase as the backend. The app offers a simple and intuitive interface, supporting group chats for smooth and seamless communication.

## Features

- Real-time messaging: Instantly send and receive messages in real-time.
- Group chats: Engage in groups to chat with multiple users.
- User-friendly interface: The app's intuitive design makes it easy for users of all levels to navigate and use.
- Online status: Know when other users are online and available for chat.
- Image sharing: Share images or youtube videos with your friends to make conversations more engaging.
- Message history: Access your chat history and previous conversations (up to 256 messages).

## Installation

1. Clone the repository: `git clone https://github.com/tastycrayon/whisker`

2. Navigate to the project directory: `cd whisker`

3. Build the Docker image:

   ```bash
   docker build -t whisker .
   ```

4. Run the Docker container:

   ```bash
   docker run -p 8080:8080 whisker
   ```

## Configuration

### Environment Variables

The Chat App can be customized using the following environment variables:

- `PUBLIC_POCKETBASE_URL`: Set the public URL for PocketBase. By default, it uses `https://whisker.fly.dev`.
- `PUBLIC_WEBSOCKET_URL`: Set the public WebSocket URL for real-time communication. By default, it uses `wss://whisker.fly.dev/ws`.

Environment variable can be edited in `Dockerfile` and `./client/.env`.

## Developement Environment

### Go/Pocketbase

1. Start pocketbase server with debug mode

   ```bash
   go run main.go serve --dir ./pb_data --debug
   ```

   Or. Use Air for hot reload

   ```bash
   go install github.com/cosmtrek/air@latest
   air serve --dir ./pb_data --debug
   ```

2. Import collection schema from `pb.collection.json`

### SvelteKit

1. Enter client folder and install packages

   ```bash
   cd client
   npm i
   ```

2. Run devlopement server on http://127.0.0.1:5173
   ```bash
   npm run dev -- --host 127.0.0.1
   ```

## Technologies Used

- Golang: Programming language used to build the backend server for the Chat App.
- Svelte: Frontend framework used to build the user interface for the web application.
- PocketBase: Backend database for storing and managing chat data in real-time.

## Contributing

Contributions are welcome! If you find any issues or have suggestions for improvement, feel free to submit a pull request.

## License

This project is licensed under the [MIT License](LICENSE).

## Contact

For any inquiries or feedback, please email me at [mosiurdev@gmail.com](mailto:mosiurdev@gmail.com)
