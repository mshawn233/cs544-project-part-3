package main

import (
	"context"
	"crypto/tls"
	"io"
	"log"

	"mshawn233/cs544-project-part-3/myproto"

	"github.com/quic-go/quic-go"
)

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

func RecvPacket(stream quic.Stream) error {

	n, err := stream.Read(_recvBuffer)
	if (err != nil) && (err != io.ErrUnexpectedEOF) {
		log.Printf("ERROR GETTING SERVER RESPONSE %+v", err)
		return err
	}

	pckt, _ := myproto.FromBytes(_recvBuffer[0:n])
	log.Printf("RECEIVED\n----------\n%s\n----------", pckt.ToString())
	return nil
}

func SendMessage(stream quic.Stream, msg string) (int, error) {
	pckt := myproto.NewPacketWithFlags([]byte(msg), myproto.FLG_SEND)

	//Convert to bytes
	netBytes, err := myproto.ToBytes(pckt)
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
	log.Printf("Just wrote %d bytes to server\n%s\n", n, pckt.ToString())
	return n, err
}

func main() {
	c, _ := InitConnection("localhost:4242")
	//defer c.Close()
	SendMessage(c, "Hello there")
	RecvPacket(c)

	c.Close()
	<-c.Context().Done()
}
