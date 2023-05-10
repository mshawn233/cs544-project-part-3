package main

import "log"

//adapted from https://github.com/quic-go/quic-go/blob/master/example/echo/echo.go

const addr = "localhost:4242"

// We start a server echoing data on the first stream the client opens,
// then connect with a client, send the message, and wait for its receipt.
func main() {
	const message = "foobar"

	go func() { log.Fatal(EchoServer()) }()

	err := ClientMain(addr, message)
	if err != nil {
		panic(err)
	}
}
