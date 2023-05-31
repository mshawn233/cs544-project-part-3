package helpers

import (
	"bytes"
	"encoding/gob"

	"mshawn233/cs544-project-part-3/chatmessagetypes"
)

func ToBytes[C *chatmessagetypes.ChatMessage | *chatmessagetypes.HelloChatRequest | *chatmessagetypes.HelloChatResponse | *chatmessagetypes.ChatDisconnect](cm C) ([]byte, error) {
	buff := bytes.Buffer{}
	encoder := gob.NewEncoder(&buff)

	err := encoder.Encode(cm)
	return buff.Bytes(), err
}
