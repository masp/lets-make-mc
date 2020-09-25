package net

import (
	"encoding/json"
	"errors"
	"github.com/masp/mcgo/proto"
)

type StatusResponse struct {
	Version       int
	VersionName   string
	MaxPlayers    int
	OnlinePlayers int
	Motd          string
}

func (s StatusResponse) EncodeTo(e *proto.PacketEncoder) {
	resp := map[string]interface{}{
		"version": map[string]interface{}{
			"name":     s.VersionName,
			"protocol": s.Version,
		},
		"players": map[string]interface{}{
			"max":    s.MaxPlayers,
			"online": s.OnlinePlayers,
			"sample": []string{},
		},
		"description": map[string]string{
			"text": s.Motd,
		},
	}

	str, err := json.Marshal(resp)
	if err != nil {
		panic(errors.New("status response: invalid JSON status response"))
	}
	e.WriteString(string(str))
}

type Pong struct {
	Payload int64
}

func (p Pong) EncodeTo(e *proto.PacketEncoder) {
	e.WriteI64(p.Payload)
}

const (
	RequestID = 0x00
	PingID    = 0x01
)

func handleStatus(player *Player) {
	req := player.readPacket()
	if req.ID != RequestID {
		panic(errors.New("invalid status packet: expected status request"))
	}

	resp := StatusResponse{
		Version:       751,
		VersionName:   "1.16.2",
		MaxPlayers:    1337,
		OnlinePlayers: 0,
		Motd:          "What a cool server!",
	}
	player.sendPacketImmediately(0x00, resp)

	req = player.readPacket()
	if req.ID == PingID {
		player.sendPacketImmediately(0x01, Pong{req.ReadI64()})
	}
}
