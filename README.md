# Chitty-Chat
Server-chat implementation using gRPC

Start the server with: "go run Chitty_Chat.go" - Windows will ask for permission to run it, click "Allow"

Start any number of clients in different terminals with: "go run Chitty_Chat_Client.go" - It will ask for a port number, type in: "5000"

Using any of the clients opened in the terminal, you can type any message and click enter - it will then be broadcastet to all the other clients

When done with a client, type in "quit" to leave the server and terminate the client.


