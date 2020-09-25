package main

import (
	mclog "github.com/masp/mcgo/log"
	mcnet "github.com/masp/mcgo/net"
	"github.com/masp/mcgo/pstn"
	"github.com/masp/mcgo/worlds"
	log "github.com/sirupsen/logrus"
	"net"
)

func init() {
	// log.SetLevel(log.DebugLevel)
}

func main() {
	mclog.SetupLogging()
	log.Info("Opening server on port 25565")

	ln, err := net.Listen("tcp", ":25565")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	log.Info("Generating spawn...")
	world := worlds.New(pstn.Block{X: 0, Y: 65, Z: 0})
	log.Info("Finished generating spawn")
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Error(err)
			continue
		}

		go mcnet.HandlePlayer(&world, conn)
	}
}
