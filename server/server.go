package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"mshawn233/cs544-project-part-3/myproto"

	"github.com/quic-go/quic-go"
)

var recvBuffer = make([]byte, myproto.MAX_PACKET_SZ)

func InitServer(addr string) (quic.Connection, error) {
	log.Printf("Server is initializing")
	listener, err := quic.ListenAddr(addr, myproto.GenerateTLSConfig(), nil)
	if err != nil {
		return nil, err
	}

	log.Printf("Getting ready to accept connections")
	return listener.Accept(context.Background())
}

func HandleStream(conn quic.Connection) error {
	log.Printf("Server is waiting for a connection")
	stream, err := conn.AcceptStream(context.Background())
	if err != nil {
		return err
	}
	log.Printf("Just Got A Connection")

	//n, err := io.ReadFull(stream, recvBuffer)
	n, err := stream.Read(recvBuffer)
	log.Printf("Post Readfull")
	if (err != nil) && (err != io.ErrUnexpectedEOF) {
		log.Printf("GOT ERROR READING FROM CLIENT %+v", err)
		return err
	}
	log.Printf("SVR: Received %d bytes", n)

	pckt, _ := myproto.FromBytes(recvBuffer[0:n])
	log.Printf("RECEIVED\n----------\n%s\n----------", pckt.ToString())

	//Now Lets Send A Response, First Lets Build An Echo Packet
	pckt.ProtoFlag = myproto.FLG_SND_ACK
	pckt.Payload = []byte(fmt.Sprintf("ECHO->%0.25s", pckt.Payload))
	pckt.Length = len(pckt.Payload)

	//Now lets convert into bytes
	netBytes, err := myproto.ToBytes(pckt)
	if err != nil {
		log.Printf("Error serializing: %+v", err)
		return err
	}

	//Send to the server
	n, err = stream.Write(netBytes)
	if err != nil {
		log.Printf("Error writing to client: %+v", err)
		return err
	}
	log.Printf("Just wrote %d bytes to client\n%s\n", n, pckt.ToString())

	//This is a stream so we need to have a way for the client to have the opportunity to receive
	//the message and then close the connection, we use context for this
	connCtx := conn.Context()
	<-connCtx.Done()

	return nil
}

func main() {
	c, err := InitServer("localhost:4242")
	log.Printf("Server just initialized, error is %+v", err)
	HandleStream(c)

	//time.Sleep(time.Second * 2)
}
