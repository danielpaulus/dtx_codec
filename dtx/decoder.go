package dtx

import (
	"encoding/binary"
	"fmt"
)

type DtxMessage struct {
	MessageType   int
	MessageLength int
	Fragments     uint16
	FragmentIndex uint16
	ExpectsReply  bool
}

func (d DtxMessage) String() string {
	return "message"
}

const (
	MethodInvocationWithExpectedReply    = 0x3
	MethodinvocationWithoutExpectedReply = 0x2
)

const (
	DtxMessageMagic uint32 = 0x795B3D1F
	DtxHeaderLength uint32 = 32
	DtxReservedBits uint32 = 0x0
)

func Decode(messageBytes []byte) (DtxMessage, error) {

	if binary.BigEndian.Uint32(messageBytes) != DtxMessageMagic {
		return DtxMessage{}, fmt.Errorf("Wrong Magic: %x", messageBytes[0:4])
	}
	if binary.LittleEndian.Uint32(messageBytes[4:]) != DtxHeaderLength {
		return DtxMessage{}, fmt.Errorf("Incorrect Header length, should be 32: %x", messageBytes[4:8])
	}
	result := DtxMessage{}
	result.FragmentIndex = binary.LittleEndian.Uint16(messageBytes[8:])
	result.Fragments = binary.LittleEndian.Uint16(messageBytes[10:])
	result.MessageLength = int(binary.LittleEndian.Uint32(messageBytes[12:]))
	result.MessageType = int(binary.LittleEndian.Uint32(messageBytes[16:]))
	if binary.LittleEndian.Uint64(messageBytes[20:]) != 0 {
		return DtxMessage{}, fmt.Errorf("Reserved bits should be 0 but are %x", messageBytes[20:28])
	}
	result.ExpectsReply = binary.LittleEndian.Uint32(messageBytes[28:]) == uint32(1)
	return result, nil
}
