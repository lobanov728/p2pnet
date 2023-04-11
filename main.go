package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/Lobanov728/p2pnet/client"
	"github.com/Lobanov728/p2pnet/messenger"
	"github.com/Lobanov728/p2pnet/server"
)

const (
	buffSize = 1024
	protocol = "tcp"
)

var (
	port, name string
)

func main() {
	var nextArgName, nextArgPort bool
	for _, val := range os.Args {
		if nextArgName {
			nextArgName = false
			name = val
			continue
		}
		if nextArgPort {
			nextArgPort = false
			port = val
			continue
		}

		if val == "-n" || val == "--name" {
			nextArgName = true
		}

		if val == "-p" || val == "--port" {
			nextArgPort = true
		}
	}

	if port == "" {
		log.Fatalln("port undefined")
	}

	// buffer := make([]byte, buffSize)

	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:"+port)
	if err != nil {
		log.Fatalln(err.Error())
	}
	messanger := messenger.NewMessenger(protocol, addr, name)

	go client.Run(messanger)

	go server.Run(messanger)

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT)
	<-ch
}
