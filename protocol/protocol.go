package protocol

import (
	"encoding/binary"
	"fmt"

	"github.com/google/uuid"
)

const (
	commandFieldLen = 4
	IDLen           = 36
	FromLen         = 170
	contentFieldLen = 2
)

const (
	CmdHand = "HAND"
	CmdResp = "RESP"
)

var headerLen = commandFieldLen + IDLen + FromLen

type Message struct {
	Command []byte
	ID      []byte
	From    []byte
	Length  uint16
	Content []byte
}

func NewMessage(cmd string) *Message {
	return &Message{
		Command: []byte(cmd)[:commandFieldLen],
		ID:      []byte(uuid.New().String()),
		From:    make([]byte, FromLen),
	}
}

func (m *Message) Serialize() []byte {
	result := make([]byte, 0, headerLen+contentFieldLen+len(m.Content))

	result = append(result, m.Command[:commandFieldLen]...)
	result = append(result, m.ID[:IDLen]...)
	result = append(result, m.From[:FromLen]...)

	contentLengthBytes := make([]byte, contentFieldLen)
	binary.BigEndian.PutUint16(contentLengthBytes, m.Length)

	result = append(result, contentLengthBytes...)
	result = append(result, m.Content...)

	return result
}

func Unserialize(b []byte) (msg *Message) {
	contentLength := binary.BigEndian.Uint16(b[headerLen : headerLen+commandFieldLen])
	fmt.Println(contentLength)
	if contentLength > 65535 {
		return nil
	}

	msg = &Message{
		Command: b[0:commandFieldLen],
		ID:      b[commandFieldLen : commandFieldLen+IDLen],
		From:    b[commandFieldLen+IDLen : commandFieldLen+IDLen+FromLen],
		Length:  contentLength,
	}

	if len(b) == (headerLen + contentFieldLen + int(contentLength)) {
		msg.Content = b[headerLen+contentFieldLen:]
	} else {
		msg.Content = make([]byte, contentLength)
	}

	return
}
