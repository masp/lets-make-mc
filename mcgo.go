package main

import (
	"fmt"
	"log"
	"net"

	"github.com/masp/mcgo/protocol"
)

func main() {
	fmt.Println("Opening server on port 25565")

	ln, err := net.Listen("tcp", ":25565")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}

		rd := protocol.NewPacketReader(conn)
		for {
			packet, err := rd.ReadPacket()

			if packet == nil {
				log.Printf("Connection closed by player")
				break
			}

			if err != nil {
				log.Fatal(err)
			}

		}
	}

}
