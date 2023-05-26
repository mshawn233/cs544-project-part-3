package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"io"
	"log"
	"os"

	"mshawn233/cs544-project-part-3/chatmessagetypes"
	"mshawn233/cs544-project-part-3/myproto"

	"github.com/quic-go/quic-go"
)

const MAX_PACKET_SZ = 1 << 16

var _recvBuffer = make([]byte, myproto.MAX_PACKET_SZ)

func InitConnection(host string) (quic.Stream, error) {
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-security-setup"},
	}
	conn, err := quic.DialAddr(host, tlsConf, nil)
	if err != nil {
		return nil, err
	}

	return conn.OpenStreamSync(context.Background())
}

func RecieveChatMessage(stream quic.Stream) error {

	n, err := stream.Read(_recvBuffer)
	if (err != nil) && (err != io.ErrUnexpectedEOF) {
		log.Printf("ERROR GETTING SERVER RESPONSE %+v", err)
		return err
	}

	chatMessage, _ := chatmessagetypes.ChatMessageFromBytes(_recvBuffer[0:n])
	log.Printf("Chat Message Received--------%s\n", chatMessage.ChatMessageData)
	return nil
}

func ReceiveHelloChatResponse(stream quic.Stream) (uint32, error) {
	n, err := stream.Read(_recvBuffer)
	if (err != nil) && (err != io.ErrUnexpectedEOF) {
		log.Printf("ERROR GETTING SERVER RESPONSE %+v", err)
		return 0x0, err
	}

	hcr, _ := chatmessagetypes.HelloChatResponseFromBytes(_recvBuffer[0:n])
	log.Printf("HelloChatResponse Received\n")
	return hcr.ClientChatSessionId, nil
}

func SendChatMessage(stream quic.Stream, msg string, chatSessionId uint32) (int, error) {
	chatmessage := chatmessagetypes.NewChatMessage(chatSessionId, msg)

	//Convert to bytes
	netBytes, err := chatmessagetypes.ChatMessageToBytes(chatmessage)
	if err != nil {
		log.Printf("Error serializing: %+v", err)
		return 0, err
	}

	//Send to the server
	n, err := stream.Write(netBytes)
	if err != nil {
		log.Printf("Error writing to server: %+v", err)
		return 0, err
	}
	log.Printf("Just wrote %d bytes to server\n", n)
	return n, err
}

func SendHelloChatMessageRequest(stream quic.Stream, username string, password string, chatPartner string) (int, error) {
	hcr := chatmessagetypes.NewHelloChatRequest(username, password, chatPartner)

	//Convert to bytes
	netBytes, err := chatmessagetypes.HelloChatRequestToBytes(hcr)
	if err != nil {
		log.Printf("Error serializing: %+v", err)
		return 0, err
	}

	//Send to the server
	n, err := stream.Write(netBytes)
	if err != nil {
		log.Printf("Error writing to server: %+v", err)
		return 0, err
	}
	log.Printf("Just wrote %d bytes to server\n", n)
	return n, err
}

func main() {
	c, _ := InitConnection("localhost:4242")
	//defer c.Close()

	//Read input from console
	reader := bufio.NewReader(os.Stdin)
	log.Printf("Enter username: ")
	username, _ := reader.ReadString('\n')

	//Read password from console
	log.Printf("Enter password: ")
	password, _ := reader.ReadString('\n')

	//Read usename from console
	log.Printf("Enter chat partner: ")
	chatPartner, _ := reader.ReadString('\n')

	SendHelloChatMessageRequest(c, username, password, chatPartner)
	sessionId, _ := ReceiveHelloChatResponse(c)

	//Read chat text from console
	log.Printf("%s: ", username)
	chatText, _ := reader.ReadString('\n')

	SendChatMessage(c, chatText, sessionId)
	RecieveChatMessage(c)

	c.Close()
	<-c.Context().Done()
}
