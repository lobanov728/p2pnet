package messenger

import (
	"net"

	"github.com/Lobanov728/p2pnet/protocol"
)

type Messenger struct {
	Addr     *net.TCPAddr
	Protocol string
	Name     string
	Peers    map[string]*Peer
	PubKey   []byte
	privKey  []byte
}

type Peer struct {
	Name   string
	PubKey []byte
}

func NewMessenger(protocolName string, addr *net.TCPAddr, name string) *Messenger {
	priv, pub := protocol.GenerateKeys()

	return &Messenger{
		Addr:     addr,
		Protocol: protocolName,
		Name:     name,
		Peers:    make(map[string]*Peer, 10),
		PubKey:   pub,
		privKey:  priv,
	}
}

func (m *Messenger) AddPeer(name string, pubKey []byte) {
	m.Peers[name] = &Peer{
		Name:   name,
		PubKey: pubKey,
	}
}
