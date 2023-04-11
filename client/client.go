package client

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/Lobanov728/p2pnet/messenger"
	"github.com/Lobanov728/p2pnet/protocol"
)

var (
	localIP        = getLocalIP()
	addrList       []string
	connectionList map[string]*net.TCPConn
)

func Run(messanger *messenger.Messenger) {
	connectionList := make(map[string]*net.TCPConn, 10)
	for {
		var (
			reader  *bufio.Reader
			command string
		)
		reader = bufio.NewReader(os.Stdin)
		command, _ = reader.ReadString('\n')
		message := strings.Replace(command, "\n", "", -1)
		messageTrim := strings.Replace(message, " ", "", -1)
		messageSplited := strings.Split(messageTrim, "->")

		fmt.Println("message", message)

		switch messageSplited[0] {
		case ":exit":
			os.Exit(0)
		case ":connect":
			if len(messageSplited) == 2 {
				addrList = append(addrList, messageSplited[1])
				addr, err := net.ResolveTCPAddr("tcp", messageSplited[1])
				if err != nil {
					continue
				}
				conn, err := net.DialTCP(messanger.Protocol, nil, addr)
				if err != nil {
					continue
				}

				msg := protocol.NewMessage(protocol.CmdHand, []byte(messanger.Name))
				msg.From = messanger.PubKey

				conn.Write(msg.Serialize())
			}

		case ":disconnect":
			if len(messageSplited) == 1 {
				for i := range addrList {
					addrList = append(addrList[:0], addrList[1:]...)
					_, ok := connectionList[addrList[i]]
					if ok {
						delete(connectionList, addrList[i])
					}
				}

			} else {
				for i := range addrList {
					if messageSplited[1] == addrList[i] {
						addrList = append(addrList[:i], addrList[i+1:]...)
						_, ok := connectionList[addrList[i]]
						if ok {
							delete(connectionList, addrList[i])
						}
						break
					}
				}
			}
		default:
			for i := range addrList {
				addr, err := net.ResolveTCPAddr("tcp", addrList[i])
				if err != nil {
					addrList = append(addrList[:i], addrList[i+1:]...)
					continue
				}
				conn, err := net.DialTCP(messanger.Protocol, nil, addr)
				// _, ok := connectionList[addrList[i]]
				// if !ok {
				// 	addr, err := net.ResolveTCPAddr("tcp", addrList[i])
				// 	conn, err := net.DialTCP(protocol, nil, addr)
				// 	if err != nil {
				// 		addrList = append(addrList[:i], addrList[i+1:]...)
				// 		continue
				// 	}
				// 	connectionList[addrList[i]] = conn
				// 	fmt.Println(conn == nil)
				// }
				// content := fmt.Sprintf("[%s/%s]: %s\n", localIP, messanger.Name, message)
				msg := protocol.NewMessage("SUCK")
				conn.Write(msg.Serialize())
			}
		}
	}
}

func getLocalIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:8080")
	defer conn.Close()
	if err != nil {
		log.Fatalln(err.Error())
	}

	return "Your local IP adress: " + strings.Split(conn.LocalAddr().String(), ":")[0]
}
