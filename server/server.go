package main

import (
	"context"
	"io"
	"log"

	"mshawn233/cs544-project-part-3/chatmessagetypes"
	"mshawn233/cs544-project-part-3/tls"

	"github.com/quic-go/quic-go"
)

const MAX_PACKET_SZ = 1 << 16

var recvBuffer = make([]byte, MAX_PACKET_SZ)

func InitServer(addr string) (quic.Connection, error) {
	log.Printf("Server is initializing")
	listener, err := quic.ListenAddr(addr, tls.GenerateTLSConfig(), nil)
	if err != nil {
		return nil, err
	}

	log.Printf("Getting ready to accept connections")
	return listener.Accept(context.Background())
}

func HandleHelloChatRequest(stream quic.Stream) error {

	//Read the request from the client
	n, err := stream.Read(recvBuffer)
	log.Printf("Post Readfull")
	if (err != nil) && (err != io.ErrUnexpectedEOF) {
		log.Printf("GOT ERROR READING FROM CLIENT %+v", err)
		return err
	}
	log.Printf("SVR: Received %d bytes", n)

	pckt, _ := chatmessagetypes.HelloChatRequestFromBytes(recvBuffer[0:n])
	log.Printf("Received request from user %s\n", pckt.Username)

	//Generate HelloChatResponse
	hcr := chatmessagetypes.NewHelloChatResponse(0x01)

	//Now lets convert into bytes
	netBytes, err := chatmessagetypes.HelloChatResponseToBytes(hcr)
	if err != nil {
		log.Printf("Error serializing: %+v", err)
		return err
	}

	//Send response to the client
	n, err = stream.Write(netBytes)
	if err != nil {
		log.Printf("Error writing to client: %+v", err)
		return err
	}
	log.Printf("Just wrote %d bytes to user %s\n", n, pckt.Username)

	return nil
}

func HandleChatMessage(stream quic.Stream) error {

	//Read the request from the client
	n, err := stream.Read(recvBuffer)
	log.Printf("Post Readfull")
	if (err != nil) && (err != io.ErrUnexpectedEOF) {
		log.Printf("GOT ERROR READING FROM CLIENT %+v", err)
		return err
	}
	log.Printf("SVR: Received %d bytes", n)

	pckt, _ := chatmessagetypes.ChatMessageFromBytes(recvBuffer[0:n])
	log.Printf("Received chat message %s\n", pckt.ChatMessageData)

	//Send response to the client
	n, err = stream.Write(recvBuffer[0:n])
	if err != nil {
		log.Printf("Error writing to client: %+v", err)
		return err
	}
	log.Printf("Just wrote %d bytes to client\n", n)

	return nil

}

func main() {

	c, err := InitServer("localhost:4242")
	log.Printf("Server just initialized, error is %+v", err)

	log.Printf("Server is waiting for a connection")
	stream, err := c.AcceptStream(context.Background())
	if err != nil {
		log.Printf("Error accepting stream: %+v", err)
	}
	log.Printf("Just Got A Connection")

	HandleHelloChatRequest(stream)

	HandleChatMessage(stream)

	//This is a stream so we need to have a way for the client to have the opportunity to receive
	//the message and then close the connection, we use context for this
	connCtx := c.Context()
	<-connCtx.Done()

}
