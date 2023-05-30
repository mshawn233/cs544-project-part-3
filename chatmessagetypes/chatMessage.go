package chatmessagetypes

import (
	"bytes"
	"encoding/gob"
)

type ChatMessage struct {
	ChatMessageType      uint8
	ClientChatSesstionId uint8
	ChatMessageData      string
}

func NewChatMessage(clientChatSesstionId uint8, chatMessageData string) *ChatMessage {
	return &ChatMessage{
		ChatMessageType:      0x2,
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
