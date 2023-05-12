package myproto

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

const PROTO_TYPE_MY_FANCY_QUIC = 0x01
const DEF_PROTO_FLAG = 0x0

const MAX_PACKET_SZ = 1 << 16

const (
	FLG_CLEAR   = 0x0                //In bits 00000000
	FLG_SEND    = 0x1                //In bits 00000001
	FLG_ACK     = 0x2                //In bits 00000010
	FLG_SND_ACK = FLG_SEND & FLG_ACK //In bits 00000011
)

type MyPDU struct {
	ProtoType uint8
	ProtoFlag uint8
	Length    int
}

type MyPacket struct {
	MyPDU
	Payload []byte
}

func NewPacket(data []byte) *MyPacket {

	return NewPacketWithFlags(data, FLG_CLEAR)
}

func NewPacketWithFlags(data []byte, flags uint8) *MyPacket {

	return &MyPacket{
		MyPDU: MyPDU{ProtoType: PROTO_TYPE_MY_FANCY_QUIC,
			ProtoFlag: flags,
			Length:    len(data)},
		Payload: data,
	}
}

func (p *MyPacket) ToString() string {
	return fmt.Sprintf("Type:\t 0x%.2x\n"+
		"Flags:\t 0x%.2x\n"+
		"Length:\t %d\n"+
		"Data:\t [%.25s] (up to first 25 bytes)",
		p.ProtoType, p.ProtoFlag, p.Length, p.Payload)
}

func ToBytes(p *MyPacket) ([]byte, error) {
	buff := bytes.Buffer{}
	encoder := gob.NewEncoder(&buff)

	err := encoder.Encode(p)
	return buff.Bytes(), err
}

func FromBytes(b []byte) (*MyPacket, error) {
	buff := bytes.NewBuffer(b)
	decoder := gob.NewDecoder(buff)

	myPacket := &MyPacket{}
	err := decoder.Decode(myPacket)
	return myPacket, err
}
