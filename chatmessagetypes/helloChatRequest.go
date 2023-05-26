package chatmessagetypes

import (
	"bytes"
	"encoding/gob"
)

type HelloChatRequest struct {
	ChatMessageType      uint8
	ClientChatSesstionId uint32
	Username             string
	Password             string
	ChatPartner          string
}

func NewHelloChatRequest(username string, password string, chatPartner string) *HelloChatRequest {

	return &HelloChatRequest{
		ChatMessageType:      0x00,
		ClientChatSesstionId: 0x00,
		Username:             username,
		Password:             password,
		ChatPartner:          chatPartner,
	}
}

func HelloChatRequestToBytes(hcr *HelloChatRequest) ([]byte, error) {
	buff := bytes.Buffer{}
	encoder := gob.NewEncoder(&buff)

	err := encoder.Encode(hcr)
	return buff.Bytes(), err
}

func HelloChatRequestFromBytes(b []byte) (*HelloChatRequest, error) {
	buff := bytes.NewBuffer(b)
	decoder := gob.NewDecoder(buff)

	myHelloChatRequest := &HelloChatRequest{}
	err := decoder.Decode(myHelloChatRequest)
	return myHelloChatRequest, err
}
