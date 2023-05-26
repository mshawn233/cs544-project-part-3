package chatmessagetypes

import (
	"bytes"
	"encoding/gob"
)

type ChatMessage struct {
	ChatMessageType      uint8
	ClientChatSesstionId uint32
	ChatMessageData      []byte
}

func NewChatMessage(clientChatSesstionId uint32, chatMessageData []byte) *ChatMessage {
	return &ChatMessage{
		ClientChatSesstionId: clientChatSesstionId,
		ChatMessageData:      chatMessageData,
	}
}

func ChatMessageToBytes(cm *ChatMessage) ([]byte, error) {
	buff := bytes.Buffer{}
	encoder := gob.NewEncoder(&buff)

	err := encoder.Encode(cm)
	return buff.Bytes(), err
}

func ChatMessageFromBytes(b []byte) (*ChatMessage, error) {
	buff := bytes.NewBuffer(b)
	decoder := gob.NewDecoder(buff)

	myChatMessage := &ChatMessage{}
	err := decoder.Decode(myChatMessage)
	return myChatMessage, err
}
