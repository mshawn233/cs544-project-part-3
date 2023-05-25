package chatmessagetypes

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

type ChatDisconnect struct {
	ChatMessageType uint8
	ClientChatSesstionId uint32
}

func NewChatDisconnect(clientChatSesstionId uint32) *ChatDisconnect {
	return &ChatDisconnect{
		ChatMessageType: 0x3,
		ClientChatSesstionId: clientChatSesstionId,
	}
}

func ToBytes(cd *ChatDisconnect) ([]byte, error) {
	buff := bytes.Buffer{}
	encoder := gob.NewEncoder(&buff)

	err := encoder.Encode(cd)
	return buff.Bytes(), err
}

func FromBytes(b []byte) (*ChatDisconnect, error) {
	buff := bytes.NewBuffer(b)
	decoder := gob.NewDecoder(buff)

	myChatDisconnect := &ChatDisconnect{}
	err := decoder.Decode(myChatDisconnect)
	return myChatDisconnect, err
}
