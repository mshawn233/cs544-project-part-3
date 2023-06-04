# SM Chat Protocal Implementation using Go and quic-go library

The Makefile can be used to compile the server and client:

`make build` : This build 2 binaries, `./cli` and `./svr`

Or to compile each separately, navigate to the "/server" and "/client" directories and run, respectively:
    - "go build server.go"
    - "go build client.go"

Once the executables are created, each can be run in the terminal by:
    - "./server.exe"
    - "./client.exe"

It is recommended to start the server (server.go) first and then start the client (client.go). The server and the client both use the quic-go library (https://pkg.go.dev/github.com/quic-go/quic-go@v0.34.0) to transport chat messages between themselves. When the server starts, it begins waiting for a connection from the client. The server process is blocked until a client connection is made. Once a connection is made the client will send a HelloChatRequest over an open stream to the server. The server will process the reques and return a HelloChatResponse through the same stream back to the client. Once the initial request and response round trip is completed, the client will send a chat message to the server. The client will ask for the username, password, and chatpartner from the client. (A working username and password is "Shawn" and "pass1"). These values are taken from the terminal input and are sent by pressing "Enter". When the server receives a chat message, it will echo the contents of the message back to the client. The client then sends a ChatDisconnect and closes it's connection with the server.

Note, currently the server uses go code to manage the required TLS crypto for QUIC, down the road ill likely enable option to use generated certs (e.g., `make gen-certs`)

Areas for improvement:
1. Allow the client and server to send more than one chat message back and forth