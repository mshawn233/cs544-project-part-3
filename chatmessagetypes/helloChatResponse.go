package chatmessagetypes

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

type HelloChatResponse struct {
	ChatMessageType uint8
	ClientChatSesstionId uint32
}

func NewHelloChatResponse(clientChatSesstionId uint32) *HelloChatResponse {
	return &HelloChatResponse{
		ChatMessageType: 0x1,
		ClientChatSesstionId: clientChatSesstionId
	}
}

func ToBytes(hcr *HelloChatResponse) ([]byte, error) {
	buff := bytes.Buffer{}
	encoder := gob.NewEncoder(&buff)

	err := encoder.Encode(hcr)
	return buff.Bytes(), err
}

func FromBytes(b []byte) (*HelloChatResponse, error) {
	buff := bytes.NewBuffer(b)
	decoder := gob.NewDecoder(buff)

	myHelloChatResponse := &HelloChatResponse{}
	err := decoder.Decode(myHelloChatResponse)
	return myHelloChatResponse, err
}
