package net

import (
	"context"
	"github.com/masp/mcgo/proto"
	"github.com/masp/mcgo/pstn"
	"github.com/masp/mcgo/worlds"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"net"
	"time"
)

type Player struct {
	Conn          net.Conn
	socketDecoder *proto.PacketDecoder
	socketEncoder *proto.PacketEncoder

	stopAll       context.CancelFunc // stops all goroutines that are spawned for handling this player
	packetsToSend chan []byte

	Version  int
	UUID     uuid.UUID
	EID      int32
	Username string
	Pos      pstn.Entity

	LastKeepAlive time.Time
}

func (p *Player) ChunkPos() pstn.Chunk {
	return pstn.EntityToChunk(p.Pos)
}

func (p *Player) Disconnect(err error) {
	log.Errorf("client '%s' disconnected with error: %v\n",
		p.Conn.RemoteAddr().String(), err)
	p.stopAll()
	p.Conn.Close()
}

func (p *Player) readPacket() proto.RecvPacket {
	return p.socketDecoder.ReadPacket()
}

func (p *Player) writePacket(packet []byte) {
	// No support for compression
	p.socketEncoder.WritePacket(packet)
}

// sendPacketImmediately is a silly helper to make it cleaner in handshaking to send packet synchronously
func (p *Player) sendPacketImmediately(id proto.PacketID, packet proto.EncodableAsPacket) {
	data := proto.EncodePacket(id, packet)
	log.Infof("Sent packet %d of len %d", id, len(data))
	p.writePacket(data)
	p.socketEncoder.Flush()
}

func (p *Player) sendPacketImmediatelyUsing(id proto.PacketID, encodeFunc proto.EncodeFunc) {
	p.sendPacketImmediately(id, proto.NewPacket(encodeFunc))
}

// SendPacket is a threadsafe way to send a packet to a player. If the buffer to send to a player is full, the packet
// is dropped and logged that the player is lagging.
func (p *Player) SendPacket(id proto.PacketID, packet proto.EncodableAsPacket) bool {
	select {
	case p.packetsToSend <- proto.EncodePacket(id, packet):
		return true
	default:
		return false
	}
}

func (p *Player) SendPacketUsing(id proto.PacketID, encodeFunc proto.EncodeFunc) bool {
	return p.SendPacket(id, proto.NewPacket(encodeFunc))
}

const defaultSendPacketsBuffered = 128

func newPlayer(conn net.Conn) (*Player, context.Context) {
	player := Player{}
	player.socketDecoder = proto.NewPacketDecoder(conn)
	player.socketEncoder = proto.NewEncoder(conn)
	player.Conn = conn

	player.packetsToSend = make(chan []byte, defaultSendPacketsBuffered)

	ctx := context.Background()
	ctx, player.stopAll = context.WithCancel(ctx)
	return &player, ctx
}

func catchPlayerPanic(player *Player) {
	err := recover()
	if err != nil {
		player.Disconnect(err.(error))
		if log.GetLevel() == log.DebugLevel {
			panic(err) /* if we're in debug, let's not handle crashes gracefully and make it easier to debug */
		}
	}
}

func HandlePlayer(world *worlds.Dimension, conn net.Conn) {
	defer conn.Close()

	player, ctx := newPlayer(conn)
	defer catchPlayerPanic(player)

	state := handleHandshake(player)
	if state == status {
		handleStatus(player)
		return
	} else if state == login {
		handleLogin(player)
		handlePlay(ctx, world, player)
	}
}
