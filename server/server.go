package server

import (
	"bufio"
	"fmt"
	"log"
	"net"

	"github.com/Lobanov728/p2pnet/messenger"
	"github.com/Lobanov728/p2pnet/protocol"
)

func Run(messenger *messenger.Messenger) {
	listener, err := net.ListenTCP(messenger.Protocol, messenger.Addr)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer listener.Close()
	fmt.Println("before accept")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln(err.Error())
		}
		go func() {
			message, err := bufio.NewReader(conn).ReadBytes('\n')
			if err != nil {
				fmt.Println(err)
			} else {
				msg := protocol.Unserialize(message)

				if err != nil {
					log.Fatalln(err.Error())
				}
				log.Println("Message Received:", string(msg.ID))

				if string(msg.Command) == protocol.CmdHand {
					messenger.AddPeer(string(msg.Content), msg.From)

					resp := protocol.NewMessage(protocol.CmdResp, []byte(messenger.Name))
					resp.From = messenger.PubKey

					conn.Write(resp.Serialize())
				}

				if string(msg.Command) == protocol.CmdResp {
					messenger.AddPeer(string(msg.Content), msg.From)
				}

			}
		}()
	}
}
