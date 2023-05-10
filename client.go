package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"

	"github.com/quic-go/quic-go"
)

func ClientMain(addr string, sndMessage string) error {
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-security-setup"},
	}
	conn, err := quic.DialAddr(addr, tlsConf, nil)
	if err != nil {
		return err
	}

	stream, err := conn.OpenStreamSync(context.Background())
	if err != nil {
		return err
	}

	fmt.Printf("Client: Sending '%s'\n", sndMessage)
	_, err = stream.Write([]byte(sndMessage))
	if err != nil {
		return err
	}

	buf := make([]byte, len(sndMessage))
	_, err = io.ReadFull(stream, buf)
	if err != nil {
		return err
	}
	fmt.Printf("Client: Got '%s'\n", buf)

	return nil
}
