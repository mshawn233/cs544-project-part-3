package main

import (
	"context"
	"fmt"
	"io"

	"github.com/quic-go/quic-go"
)

// Start a server that echos all data on the first stream opened by the client
func EchoServer() error {
	listener, err := quic.ListenAddr(addr, GenerateTLSConfig(), nil)
	if err != nil {
		return err
	}
	conn, err := listener.Accept(context.Background())
	if err != nil {
		return err
	}
	stream, err := conn.AcceptStream(context.Background())
	if err != nil {
		panic(err)
	}
	// Echo through the loggingWriter
	_, err = io.Copy(loggingWriter{stream}, stream)
	return err
}

// A wrapper for io.Writer that also logs the message.
type loggingWriter struct{ io.Writer }

func (w loggingWriter) Write(b []byte) (int, error) {
	fmt.Printf("Server: Got '%s'\n", string(b))
	return w.Writer.Write(b)
}
