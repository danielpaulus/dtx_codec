package dtx

import (
	"encoding/binary"
	"encoding/json"
	"fmt"

	"github.com/danielpaulus/nskeyedarchiver"
)

type DtxMessage struct {
	Fragments         uint16
	FragmentIndex     uint16
	MessageLength     int
	Identifier        int
	ConversationIndex int
	ChannelCode       int
	ExpectsReply      bool
	PayloadHeader     DtxPayloadHeader
	rawBytes          []byte
}

//16 Bytes
type DtxPayloadHeader struct {
	MessageType        int
	AuxiliaryLength    int
	TotalPayloadLength int
	Flags              int
}

func (d DtxMessage) String() string {
	var e = ""
	if d.ExpectsReply {
		e = "e"
	}
	msgtype := fmt.Sprintf("Unknown:%d", d.PayloadHeader.MessageType)
	if knowntype, ok := messageTypeLookup[d.PayloadHeader.MessageType]; ok {
		msgtype = knowntype
	}

	return fmt.Sprintf("i%d.%d%s c%d t:%s mlen:%d aux_len%d paylen%d", d.Identifier, d.ConversationIndex, e, d.ChannelCode, msgtype,
		d.MessageLength, d.PayloadHeader.AuxiliaryLength, d.PayloadHeader.PayloadLength())
}

func (d DtxMessage) StringDebug() string {
	auxHeaderBytes := make([]byte, 0)
	aux_bytes := make([]byte, 0)
	if d.PayloadHeader.HasAuxiliary() {
		auxHeaderBytes = d.rawBytes[48:64]
		aux_bytes = d.rawBytes[64 : 48+d.PayloadHeader.AuxiliaryLength]
		payloadBytes := make([]byte, 0)
		var b []byte
		if d.PayloadHeader.HasPayload() {
			payloadBytes = d.rawBytes[48+d.PayloadHeader.AuxiliaryLength:]
			payloadValue, _ := nskeyedarchiver.Unarchive(payloadBytes)
			b, _ = json.Marshal(payloadValue[0])

		}
		return fmt.Sprintf("auxheader:%x\naux:%x\npayload: %s \nrawbytes:%x", auxHeaderBytes, aux_bytes, b, d.rawBytes)
	}
	if d.PayloadHeader.HasPayload() {
		payloadBytes := d.rawBytes[48:]
		payloadValue, _ := nskeyedarchiver.Unarchive(payloadBytes)
		b, _ := json.Marshal(payloadValue[0])
		return fmt.Sprintf("no aux,payload: %s \nrawbytes:%x", b, d.rawBytes)
	}
	return fmt.Sprintf("\nrawbytes:%x", d.rawBytes)
}
func (d DtxMessage) GetPayloadBytes() []byte {

	if d.PayloadHeader.HasAuxiliary() {
		payloadBytes := make([]byte, 0)
		if d.PayloadHeader.HasPayload() {
			payloadBytes = d.rawBytes[48+d.PayloadHeader.AuxiliaryLength:]
			//			a, _ := archiver.Unarchive(payloadBytes)
			//			log.Fatal(a)
			return payloadBytes
		}

	}
	if d.PayloadHeader.HasPayload() {
		payloadBytes := d.rawBytes[48:]
		//a, _ := archiver.Unarchive(payloadBytes)
		//log.Fatal(a)
		return payloadBytes
	}
	return nil
}

func (d DtxPayloadHeader) PayloadLength() int {
	return d.TotalPayloadLength - d.AuxiliaryLength
}

func (d DtxPayloadHeader) HasAuxiliary() bool {
	return d.AuxiliaryLength > 0
}

func (d DtxPayloadHeader) HasPayload() bool {
	return d.PayloadLength() > 0
}

const (
	MethodInvocationWithExpectedReply    = 0x3
	MethodinvocationWithoutExpectedReply = 0x2
	Ack                                  = 0x0
)

var messageTypeLookup = map[int]string{
	MethodInvocationWithExpectedReply:    `rpc_void`,
	MethodinvocationWithoutExpectedReply: `rpc_asking_reply`,
	Ack:                                  `Ack`,
}

const (
	DtxMessageMagic uint32 = 0x795B3D1F
	DtxHeaderLength uint32 = 32
	DtxReservedBits uint32 = 0x0
)

func Decode(messageBytes []byte) (DtxMessage, []byte, error) {

	if binary.BigEndian.Uint32(messageBytes) != DtxMessageMagic {
		return DtxMessage{}, make([]byte, 0), fmt.Errorf("Wrong Magic: %x", messageBytes[0:4])
	}
	if binary.LittleEndian.Uint32(messageBytes[4:]) != DtxHeaderLength {
		return DtxMessage{}, make([]byte, 0), fmt.Errorf("Incorrect Header length, should be 32: %x", messageBytes[4:8])
	}
	result := DtxMessage{}
	result.FragmentIndex = binary.LittleEndian.Uint16(messageBytes[8:])
	result.Fragments = binary.LittleEndian.Uint16(messageBytes[10:])
	result.MessageLength = int(binary.LittleEndian.Uint32(messageBytes[12:]))
	result.Identifier = int(binary.LittleEndian.Uint32(messageBytes[16:]))
	result.ConversationIndex = int(binary.LittleEndian.Uint32(messageBytes[20:]))
	result.ChannelCode = int(binary.LittleEndian.Uint32(messageBytes[24:]))

	result.ExpectsReply = binary.LittleEndian.Uint32(messageBytes[28:]) == uint32(1)
	ph, err := parsePayloadHeader(messageBytes[32:48])
	if err != nil {
		return DtxMessage{}, make([]byte, 0), err
	}
	result.PayloadHeader = ph
	totalMessageLength := result.MessageLength + int(DtxHeaderLength)
	result.rawBytes = messageBytes[:totalMessageLength]
	remainingBytes := messageBytes[totalMessageLength:]
	return result, remainingBytes, nil
}

func parsePayloadHeader(messageBytes []byte) (DtxPayloadHeader, error) {
	result := DtxPayloadHeader{}
	result.MessageType = int(binary.LittleEndian.Uint32(messageBytes))
	result.AuxiliaryLength = int(binary.LittleEndian.Uint32(messageBytes[4:]))
	result.TotalPayloadLength = int(binary.LittleEndian.Uint32(messageBytes[8:]))
	result.Flags = int(binary.LittleEndian.Uint32(messageBytes[12:]))

	return result, nil
}
