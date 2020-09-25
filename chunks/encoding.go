package chunks

import (
	"bytes"
	"github.com/masp/mcgo/biome"
	"github.com/masp/mcgo/proto"
)

func (c *ChunkColumn) primaryBitmask() int32 {
	var mask int32
	for chunkY, sec := range c.Sections {
		if !sec.IsAir() {
			mask |= 1 << chunkY
		}
	}
	return mask
}

const (
	bitsPerBlock = 16
)

func (c *ChunkColumn) EncodeTo(enc *proto.PacketEncoder) {
	enc.WriteI32(c.Pos.X)
	enc.WriteI32(c.Pos.Z)
	fullChunk := true
	enc.WriteBool(fullChunk)
	enc.WriteVar32(c.primaryBitmask())
	enc.WriteNBT(c.HeightMap())
	if fullChunk {
		// Even though we specify the length, this must always be 1024 to match what the client expects
		const biomeSize = 1024
		enc.WriteVar32(biomeSize)
		for i := 0; i < biomeSize; i++ {
			enc.WriteVar32(biome.PlainsID)
		}
	}

	var secBuffer bytes.Buffer
	secEnc := proto.NewEncoder(&secBuffer)
	for _, sec := range c.Sections {
		sec.EncodeTo(secEnc)
	}

	// Chunk data
	enc.WriteVar32(int32(secBuffer.Len()))
	_, _ = enc.Write(secBuffer.Bytes())

	enc.WriteVar32(0) // Block entities
}

func (s *ChunkSection) EncodeTo(enc *proto.PacketEncoder) {
	if !s.IsAir() {
		enc.WriteI16(int16(s.totalNonAirBlocks))
		enc.WriteU8(bitsPerBlock)
		// No palette
		packed := NewPaddedBlockArray(bitsPerBlock, s.blocks[:])
		enc.WriteVar32(int32(len(packed.Data)))
		for _, packedBlock := range packed.Data {
			enc.WriteI64(packedBlock)
		}
	}
}

type ChunkLightingPacket struct {
	Chunk *ChunkColumn
}

func (c *ChunkColumn) skyLightBitmask() int32 {
	return 2<<18 - 1 // every chunk is 15
}

func (c *ChunkColumn) blockLightBitmask() int32 {
	return 0 // every chunk is 0
}

func (c ChunkLightingPacket) EncodeTo(e *proto.PacketEncoder) {
	e.WriteVar32(c.Chunk.Pos.X)
	e.WriteVar32(c.Chunk.Pos.Z)
	e.WriteBool(false) // trust edges

	// TODO: Don't send every chunk as fully lit
	skymask := c.Chunk.skyLightBitmask()
	blockmask := c.Chunk.blockLightBitmask()
	e.WriteVar32(skymask)
	e.WriteVar32(blockmask)
	e.WriteVar32(^skymask)
	e.WriteVar32(^blockmask)

	// Sky light
	c.Chunk.VoidSection.SkyLight.EncodeTo(e)
	for _, sec := range c.Chunk.Sections {
		light := sec.Lighting.SkyLight
		light.EncodeTo(e)
	}
	c.Chunk.SkySection.SkyLight.EncodeTo(e)
}
