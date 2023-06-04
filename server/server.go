package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"os"
	"time"

	"mshawn233/cs544-project-part-3/chatmessagetypes"
	"mshawn233/cs544-project-part-3/configs"
	"mshawn233/cs544-project-part-3/helpers"
	"mshawn233/cs544-project-part-3/tls"

	"github.com/quic-go/quic-go"
)

const MAX_PACKET_SZ = 1 << 16

var recvBuffer = make([]byte, MAX_PACKET_SZ)

func InitServer() (quic.Connection, error) {
	log.Printf("Server is initializing")

	//Read the config file
	file, err := os.ReadFile("../configs/config.json")
	if err != nil {
		return nil, err
	}

	//Get the host and port from the config file
	var config configs.Config
	err = json.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}

	//Create the address string
	var addr string = config.Host + ":" + config.Port

	listener, err := quic.ListenAddr(addr, tls.GenerateTLSConfig(), nil)
	if err != nil {
		return nil, err
	}

	log.Printf("Getting ready to accept connections")
	return listener.Accept(context.Background())
}

func HandleHelloChatRequest(stream quic.Stream, bytesRead int, userCredentialsMap map[string]string) error {

	//Create empty HelloChatResponse
	hcr := &chatmessagetypes.HelloChatResponse{}

	pckt, _ := chatmessagetypes.HelloChatRequestFromBytes(recvBuffer[0:bytesRead])
	log.Printf("Received request from user %s\n", pckt.Username)

	//Authenticate the user
	userAuthenticatedBool := authenticateClientUser(pckt.Username, pckt.Password, userCredentialsMap)

	if !userAuthenticatedBool {
		log.Printf("User %s failed authentication", pckt.Username)
		//Generate HelloChatResponse with failed authentication
		hcr = chatmessagetypes.NewHelloChatResponse(0x1, 0x0)
	} else {
		log.Printf("User %s successfully authenticated", pckt.Username)
		//Generate HelloChatResponse with successful authentication
		hcr = chatmessagetypes.NewHelloChatResponse(0x1, 0x1)
	}

	//Now lets convert into bytes
	netBytes, err := helpers.ToBytes(hcr)
	if err != nil {
		log.Printf("Error serializing: %+v", err)
		return err
	}

	//Send response to the client
	n, err := stream.Write(netBytes)
	if err != nil {
		log.Printf("Error writing to client: %+v", err)
		return err
	}
	log.Printf("Just wrote %d bytes to user %s\n", n, pckt.Username)

	return nil
}

func HandleChatMessage(stream quic.Stream, bytesRead int) error {

	pckt, _ := chatmessagetypes.ChatMessageFromBytes(recvBuffer[0:bytesRead])
	log.Printf("Received chat message %s\n", pckt.ChatMessageData)

	//Send response to the client
	n, err := stream.Write(recvBuffer[0:bytesRead])
	if err != nil {
		log.Printf("Error writing to client: %+v", err)
		return err
	}
	log.Printf("Just wrote %d bytes to client\n", n)

	return nil

}

func createUserCredentialsMap() map[string]string {

	userCredentials := make(map[string]string)

	userCredentials["Shawn"] = "pass1"
	userCredentials["Jon"] = "pass2"
	userCredentials["Steve"] = "pass3"

	return userCredentials

}

func authenticateClientUser(username string, password string, userCredentials map[string]string) bool {

	return userCredentials[username] == password

}

func main() {

	userCredentialsMap := createUserCredentialsMap()

	c, err := InitServer()
	log.Printf("Server just initialized, error is %+v", err)

	log.Printf("Server is waiting for a connection")
	stream, err := c.AcceptStream(context.Background())
	if err != nil {
		log.Printf("Error accepting stream: %+v", err)
	}
	log.Printf("Just Got A Connection")

	//Read the request from the client
	n, err := stream.Read(recvBuffer)
	if (err != nil) && (err != io.ErrUnexpectedEOF) {
		log.Printf("GOT ERROR READING FROM CLIENT %+v", err)
	}
	log.Printf("Server: Received %d bytes", n)

	time.Sleep(15 * time.Second)
	HandleHelloChatRequest(stream, n, userCredentialsMap)

	//Read the request from the client
	n, err = stream.Read(recvBuffer)
	if (err != nil) && (err != io.ErrUnexpectedEOF) {
		log.Printf("GOT ERROR READING FROM CLIENT %+v", err)
	}
	log.Printf("Server: Received %d bytes", n)

	time.Sleep(15 * time.Second)
	HandleChatMessage(stream, n)

	//Give some time for the client to respond with a ChatDisconnect
	time.Sleep(15 * time.Second)

	//This is a stream so we need to have a way for the client to have the opportunity to receive
	//the message and then close the connection, we use context for this
	log.Printf("Server: Closing connection with client")
	connCtx := c.Context()
	<-connCtx.Done()

}
