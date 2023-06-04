package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"mshawn233/cs544-project-part-3/chatmessagetypes"
	"mshawn233/cs544-project-part-3/configs"
	"mshawn233/cs544-project-part-3/helpers"

	"github.com/quic-go/quic-go"
)

const MAX_PACKET_SZ = 1 << 16

var _recvBuffer = make([]byte, MAX_PACKET_SZ)

func InitConnection() (quic.Stream, error) {

	//Initialize TLS config
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-security-setup"},
	}

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

	//Dial the server
	conn, err := quic.DialAddr(addr, tlsConf, nil)
	if err != nil {
		return nil, err
	}

	//Return the QUIC stream
	return conn.OpenStreamSync(context.Background())
}

func RecieveChatMessage(stream quic.Stream) error {

	//Read the response from the server
	n, err := stream.Read(_recvBuffer)

	//Check for errors reading from stream
	if (err != nil) && (err != io.ErrUnexpectedEOF) {
		log.Printf("ERROR GETTING SERVER RESPONSE %+v", err)
		return err
	}

	//Convert bytes to ChatMessage
	chatMessage, _ := chatmessagetypes.ChatMessageFromBytes(_recvBuffer[0:n])
	log.Printf("Chat Message Received--------%s\n", chatMessage.ChatMessageData)

	return nil
}

func ReceiveHelloChatResponse(stream quic.Stream) (uint8, error) {

	//Read the response from the server
	n, err := stream.Read(_recvBuffer)

	//Check for errors reading from stream
	if (err != nil) && (err != io.ErrUnexpectedEOF) {
		log.Printf("ERROR GETTING SERVER RESPONSE %+v", err)
		return 0x0, err
	}

	//Convert bytes to HelloChatResponse
	hcr, _ := chatmessagetypes.HelloChatResponseFromBytes(_recvBuffer[0:n])
	log.Printf("HelloChatResponse Received\n")

	return hcr.ClientChatSessionId, nil
}

func SendChatDisconnect(stream quic.Stream, chatSessionId uint8) (int, error) {

	//Create new ChatDisconnect
	chatDisconnect := chatmessagetypes.NewChatDisconnect(chatSessionId)

	//Convert the ChatDisconnect to bytes
	netBytes, err := helpers.ToBytes(chatDisconnect)

	//Check for errors converting to bytes
	if err != nil {
		log.Printf("Error serializing: %+v", err)
		return 0, err
	}

	//Send byte array (ChatDisconnect) to the server
	n, err := stream.Write(netBytes)

	//Check for errors writing to stream
	if err != nil {
		log.Printf("Error writing to server: %+v", err)
		return 0, err
	}

	log.Printf("Just wrote %d bytes to server\n", n)

	return 0, nil

}

func SendChatMessage(stream quic.Stream, msg string, chatSessionId uint8) (int, error) {

	//Create new ChatMessage
	chatmessage := chatmessagetypes.NewChatMessage(chatSessionId, msg)

	//Convert the ChatMessage to bytes
	netBytes, err := helpers.ToBytes(chatmessage)

	//Check for errors converting to bytes
	if err != nil {
		log.Printf("Error serializing: %+v", err)
		return 0, err
	}

	//Send byte array (ChatMessage) to the server
	n, err := stream.Write(netBytes)

	//Check for errors writing to stream
	if err != nil {
		log.Printf("Error writing to server: %+v", err)
		return 0, err
	}
	log.Printf("Just wrote %d bytes to server\n", n)

	return n, err
}

func SendHelloChatMessageRequest(stream quic.Stream, username string, password string, chatPartner string) (int, error) {

	//Create new HelloChatRequest
	hcr := chatmessagetypes.NewHelloChatRequest(username, password, chatPartner)

	//Convert the HelloChatRequest to bytes
	netBytes, err := helpers.ToBytes(hcr)

	//Check for errors converting to bytes
	if err != nil {
		log.Printf("Error serializing: %+v", err)
		return 0, err
	}

	//Send byte array (HelloChatMessageRequest) to the server
	n, err := stream.Write(netBytes)

	//Check for errors writing to stream
	if err != nil {
		log.Printf("Error writing to server: %+v", err)
		return 0, err
	}
	log.Printf("Just wrote %d bytes to server\n", n)

	return n, err
}

func main() {
	c, _ := InitConnection()

	//Read username from console
	reader := bufio.NewReader(os.Stdin)
	log.Printf("Enter username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimRight(username, "\r\n")

	//Read password from console
	log.Printf("Enter password: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimRight(password, "\r\n")

	//Read chat partner username from console
	log.Printf("Enter chat partner: ")
	chatPartner, _ := reader.ReadString('\n')
	chatPartner = strings.TrimRight(chatPartner, "\r\n")

	//Use with debugger
	/*username := "Shawn"
	password := "pass1"
	chatPartner := "Jon"*/

	SendHelloChatMessageRequest(c, username, password, chatPartner)

	sessionId, _ := ReceiveHelloChatResponse(c)

	//Read chat text from console
	log.Printf("%s: ", username)
	chatText, _ := reader.ReadString('\n')
	chatText = strings.TrimRight(chatText, "\r\n")

	//Use with debugger
	//chatText := "Hello World!"

	SendChatMessage(c, chatText, sessionId)

	time.Sleep(20 * time.Second)
	RecieveChatMessage(c)

	SendChatDisconnect(c, sessionId)

	c.Close()
	<-c.Context().Done()

}
