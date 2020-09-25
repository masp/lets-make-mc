package net

import (
	"errors"
	"github.com/masp/mcgo/proto"
)

const handshakeID = 0x00

type playerState int32

const (
	status playerState = 1
	login  playerState = 2
)

type handshake struct {
	ProtocolVersion int32
	NextState       playerState
}

func readHandshake(p proto.RecvPacket) handshake {
	if p.ID != handshakeID {
		panic(errors.New("invalid packet ID, expected handshake"))
	}

	h := handshake{}
	h.ProtocolVersion = p.ReadVar32()
	_ = p.ReadString()
	_ = p.ReadU16()
	h.NextState = playerState(p.ReadVar32())
	return h
}

func handleHandshake(player *Player) playerState {
	h := readHandshake(player.readPacket())
	if h.NextState == status || h.NextState == login {
		player.Version = int(h.ProtocolVersion)
		return h.NextState
	} else {
		panic(errors.New("invalid next status sent in handshake (needs to be 1 status or 2 login)"))
	}
}
