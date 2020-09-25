package net

import (
	"github.com/masp/mcgo/proto"
	uuid "github.com/satori/go.uuid"
)

const (
	loginStartID = 0x00
)

type LoginSuccess struct {
	player *Player
}

func (l LoginSuccess) EncodeTo(e *proto.PacketEncoder) {
	e.WriteUUID(l.player.UUID)
	e.WriteString(l.player.Username)
}

func handleLogin(player *Player) {
	start := player.readPacket()
	if start.ID != loginStartID {
		panic("expected login start packet ID")
	}

	player.Username = start.ReadString()
	player.UUID = uuid.NewV4()
	player.sendPacketImmediately(0x02, LoginSuccess{player})
}
