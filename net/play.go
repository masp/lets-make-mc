package net

import (
	"context"
	"errors"
	"github.com/masp/mcgo/biome"
	"github.com/masp/mcgo/chunks"
	"github.com/masp/mcgo/proto"
	"github.com/masp/mcgo/pstn"
	"github.com/masp/mcgo/worlds"
	log "github.com/sirupsen/logrus"
	"time"
)

type gamemode uint8

const (
	survival  gamemode = 0
	creative  gamemode = 1
	adventure gamemode = 2
	spectator gamemode = 3
)

// IDs as of 1.16.2

// Serverbound
const (
	teleportConfirmServerboundID = 0x00
	clientSettingsID             = 0x05
	keepAliveServerboundID       = 0x10
	playerPosID                  = 0x12
	playerPosAndRotID            = 0x13
	playerRotID                  = 0x14
	playerMovementID             = 0x15
)

// Clientbound
const (
	joinGameID                    = 0x24
	keepAliveClientboundID        = 0x1F
	heldItemChangeID              = 0x3F
	declareRecipesID              = 0x5A
	updateTagsID                  = 0x5B
	entityStatusID                = 0x1A
	declareCommandsID             = 0x10
	unlockRecipesID               = 0x35
	playerInfoID                  = 0x32
	updateViewPositionID          = 0x40
	spawnPositionID               = 0x42
	playerPosAndLookClientboundID = 0x34
	chunkDataID                   = 0x20
	updateLightID                 = 0x23
)

type JoinGame struct {
	player *Player
	mode   gamemode
}

func (j JoinGame) EncodeTo(e *proto.PacketEncoder) {
	e.WriteI32(j.player.EID)
	e.WriteBool(false)
	e.WriteU8(uint8(j.mode)) // gamemode
	e.WriteU8(uint8(j.mode)) // previous gamemode

	e.WriteVar32(1) // number of worlds
	e.WriteString("minecraft:overworld")

	e.WriteNBT(biome.BuildRegistry())
	e.WriteNBT(biome.OverworldDimension())
	e.WriteString("minecraft:overworld")

	e.WriteI64(0)      // hashed seed
	e.WriteVar32(1337) // unused
	e.WriteVar32(16)   // max view distance
	e.WriteBool(false) // reduced debug info
	e.WriteBool(true)  // show respawn screen
	e.WriteBool(false) // is debug
	e.WriteBool(false) // is flat
}

type keepAlive struct {
	Id int64
}

func (k keepAlive) EncodeTo(e *proto.PacketEncoder) {
	e.WriteI64(k.Id)
}

var (
	ErrTimeout = errors.New("timed out")
)

func handleSendingPackets(ctx context.Context, player *Player) {
	defer catchPlayerPanic(player)

	flushTicker := time.NewTicker(time.Second / 20) // Force flush every 50ms for lower latency
	heartbeatTicker := time.NewTicker(time.Second)
	for {
		select {
		case <-ctx.Done():
			flushTicker.Stop()
			heartbeatTicker.Stop()
			return
		case <-heartbeatTicker.C:
			if time.Now().Sub(player.LastKeepAlive) > 30*time.Second {
				player.Disconnect(ErrTimeout)
			}
			player.writePacket(proto.EncodePacket(keepAliveClientboundID, keepAlive{time.Now().Unix()}))
		case <-flushTicker.C:
			player.socketEncoder.Flush()
		case p := <-player.packetsToSend:
			player.writePacket(p)
		}
	}
}

func handlePlay(ctx context.Context, world *worlds.Dimension, player *Player) {
	player.EID = 999 // TODO: register entities
	player.LastKeepAlive = time.Now()
	player.sendPacketImmediately(joinGameID, JoinGame{
		player: player,
		mode:   creative,
	})

	spawnPlayer(world, player)
	go handleSendingPackets(ctx, player)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			p := player.readPacket()
			player.handlePacket(p)
		}
	}
}

func (p *Player) handlePacket(packet proto.RecvPacket) {
	switch packet.ID {
	case keepAliveServerboundID:
		p.LastKeepAlive = time.Now()
	case teleportConfirmServerboundID:
		log.Info("TODO: TeleportConfirm packet")
		// TODO
	case clientSettingsID:
		log.Info("TODO: ClientSettings packet")
		// TODO
	case playerPosID:
		play.Handle()
	case playerPosAndRotID:
	case playerRotID:
	case playerMovementID:

	default:
		log.Infof("Received unknown packet 0x%2x, ignoring", packet.ID)
	}
}

func spawnPlayer(world *worlds.Dimension, p *Player) {
	p.Pos = pstn.BlockToEntity(world.Spawn)
	// TODO: Send held item
	p.sendPacketImmediatelyUsing(heldItemChangeID, func(e *proto.PacketEncoder) {
		e.WriteI8(1)
	})

	// TODO: Send recipes
	p.sendPacketImmediatelyUsing(declareRecipesID, func(e *proto.PacketEncoder) {
		e.WriteVar32(0)
	})

	// TODO: Send tags
	/*p.sendPacketImmediatelyUsing(updateTagsID, func(e *proto.PacketEncoder) {
		e.WriteVar32(0) // Blocks
		e.WriteVar32(0) // Items
		e.WriteVar32(1) // Fluids
		e.WriteString("minecraft:water")
		e.WriteVar32(2)
		e.WriteVar32(blocks.Water)
		e.WriteVar32(blocks.WaterLow)
		e.WriteVar32(0) // Entities
	})*/

	// TODO: Send player status
	/*p.sendPacketImmediatelyUsing(entityStatusID, func(e *proto.PacketEncoder) {
		e.WriteI32(p.EID)
		e.WriteI8(2)
	})*/

	// TODO: Declare commands
	p.sendPacketImmediatelyUsing(declareCommandsID, func(e *proto.PacketEncoder) {
		e.WriteVar32(1) // Number of Command Nodes
		e.WriteI8(0x00) // Is root node not executable
		e.WriteVar32(0) // Has no children (0)

		e.WriteVar32(0) // Root Command Node is at 0
	})

	// TODO: Unlock recipes
	p.sendPacketImmediatelyUsing(unlockRecipesID, func(e *proto.PacketEncoder) {
		e.WriteVar32(0 /* init */)
		e.WriteBool(false) // crafting recipes
		e.WriteBool(false)
		e.WriteBool(false) // smelting recipes
		e.WriteBool(false)
		e.WriteBool(false) // blast furnace recipes
		e.WriteBool(false)
		e.WriteBool(false) // smoker recipes
		e.WriteBool(false)
		e.WriteVar32(0)
		e.WriteVar32(0)
	})

	// TODO: Player info - add player
	p.sendPacketImmediatelyUsing(playerInfoID, func(e *proto.PacketEncoder) {
		e.WriteVar32(0) // add player
		e.WriteVar32(1) // 1 player
		e.WriteUUID(p.UUID)
		e.WriteString(p.Username)
		e.WriteVar32(0)               // no properties
		e.WriteVar32(int32(creative)) // gamemode
		e.WriteVar32(0)               // ping
		e.WriteBool(false)            // has display name
	})

	// TODO: Player info - update latency
	p.sendPacketImmediatelyUsing(playerInfoID, func(e *proto.PacketEncoder) {
		e.WriteVar32(2) // update latency
		e.WriteVar32(1) // 1 player
		e.WriteUUID(p.UUID)
		e.WriteVar32(0) // ping
	})

	// TODO: Update View Position
	p.sendPacketImmediatelyUsing(updateViewPositionID, func(e *proto.PacketEncoder) {
		e.WriteVar32(p.ChunkPos().X)
		e.WriteVar32(p.ChunkPos().Z)
	})

	sendChunks(world, p)

	p.sendPacketImmediatelyUsing(spawnPositionID, func(e *proto.PacketEncoder) {
		e.WritePosition(world.Spawn)
	})

	// TODO: Player position and look
	p.sendPacketImmediatelyUsing(playerPosAndLookClientboundID, func(e *proto.PacketEncoder) {
		e.WriteFloat64(p.Pos.X) // X
		e.WriteFloat64(p.Pos.Y) // Y
		e.WriteFloat64(p.Pos.Z) // Z
		e.WriteFloat32(0)       // Yaw
		e.WriteFloat32(0)       // Pitch
		e.WriteI8(0)
		e.WriteVar32(0)
	})
}

const ViewingDistance = 8

func sendChunks(world *worlds.Dimension, p *Player) {
	center := p.ChunkPos()
	for x := center.X - ViewingDistance; x <= center.X+ViewingDistance; x++ {
		for z := center.Z - ViewingDistance; z <= center.Z+ViewingDistance; z++ {
			// Send chunk data
			chunk := world.ChunkAt(pstn.Chunk{X: x, Z: z})
			p.sendPacketImmediately(chunkDataID, chunk)
			// TODO: Send proper chunk lighting
			p.sendPacketImmediately(updateLightID, chunks.ChunkLightingPacket{Chunk: chunk})
		}
	}
}
