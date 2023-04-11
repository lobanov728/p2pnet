package protocol

import (
	"testing"
)

func TestNewMessage(t *testing.T) {
	nMsg := NewMessage("MESS")
	nMsg.Length = uint16(12)
	nMsg.Content = []byte("some message")
	ser := nMsg.Serialize()

	msg := Unserialize(ser)
	if string(msg.Content) != "some message" {
		t.Error("wrong content", string(msg.Content))
	}
}

func TestEncryptDecryptMessage(t *testing.T) {
	privKey, pubKey := GenerateKeys()

	nMsg := NewMessage("MESS")
	nMsg.From = pubKey

	encryptedMsg, err := EncryptMessage(pubKey, []byte("qwertyuiopasdfghjklzxcvbnm"))
	if err != nil {
		t.Error("got error", err.Error())
	}
	nMsg.Content = encryptedMsg
	nMsg.Length = uint16(len(encryptedMsg))

	msg := Unserialize(nMsg.Serialize())
	decryptedMsg, err := DecryptMessage(privKey, msg.Content)
	if err != nil {
		t.Error("got error", err.Error())
	}
	if string(decryptedMsg) != "qwertyuiopasdfghjklzxcvbnm" {
		t.Error("wrong content", string(msg.Content))
	}
}
