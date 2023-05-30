package chatmessagetypes

import (
	"bytes"
	"encoding/gob"
)

type HelloChatResponse struct {
	ChatMessageType     uint8
	ClientChatSessionId uint8
}

func NewHelloChatResponse(chatmessagetype uint8, clientChatSessionId uint8) *HelloChatResponse {
	return &HelloChatResponse{
		ChatMessageType:     chatmessagetype,
		ClientChatSessionId: clientChatSessionId,
	}
}

func HelloChatResponseToBytes(hcr *HelloChatResponse) ([]byte, error) {
	buff := bytes.Buffer{}
	encoder := gob.NewEncoder(&buff)

	err := encoder.Encode(hcr)
	return buff.Bytes(), err
}

func HelloChatResponseFromBytes(b []byte) (*HelloChatResponse, error) {
	buff := bytes.NewBuffer(b)
	decoder := gob.NewDecoder(buff)

	myHelloChatResponse := &HelloChatResponse{}
	err := decoder.Decode(myHelloChatResponse)
	return myHelloChatResponse, err
}
