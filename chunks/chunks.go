package chunks

import (
	"github.com/masp/mcgo/blocks"
	"github.com/masp/mcgo/pstn"
)

const (
	Width           = 16
	Depth           = 16
	Height          = 256
	SectionHeight   = 16
	SectionsInChunk = Height / SectionHeight
	BlocksInSection = Width * SectionHeight * Depth
)

type BlockState uint16 // TODO: This should definitely not be an int

type ChunkLighting struct {
	SkyLight   LightingArray
	BlockLight LightingArray
}

type ChunkSection struct {
	blocks            [BlocksInSection]BlockState
	Lighting          ChunkLighting
	totalNonAirBlocks int
}

func (s *ChunkSection) index(x int, y int, z int) int {
	return (y&0xf)<<8 | z<<4 | x
}

func (s *ChunkSection) BlockAt(x int, y int, z int) BlockState {
	return s.blocks[s.index(x, y, z)]
}

func (s *ChunkSection) SetBlockAt(x int, y int, z int, newBlock BlockState) {
	prev := s.BlockAt(x, y, z)
	if prev == blocks.Air && newBlock != blocks.Air {
		s.totalNonAirBlocks++
	} else if prev != blocks.Air && newBlock == blocks.Air {
		s.totalNonAirBlocks--
	}
	s.blocks[s.index(x, y, z)] = newBlock
}

func (s *ChunkSection) IsAir() bool {
	return s.totalNonAirBlocks == 0
}

type ChunkColumn struct {
	Pos pstn.Chunk

	Sections    [SectionsInChunk]ChunkSection
	VoidSection ChunkLighting // y=-16 to y=-1
	SkySection  ChunkLighting // y=256 to y=271
}

func NewChunk(pos pstn.Chunk) *ChunkColumn {
	return &ChunkColumn{Pos: pos}
}

func (c *ChunkColumn) BlockAt(x int, y int, z int) BlockState {
	return c.Sections[y/16].BlockAt(x, y, z)
}

func (c *ChunkColumn) SetBlockAt(x int, y int, z int, newBlock BlockState) {
	c.Sections[y/16].SetBlockAt(x, y, z, newBlock)
}

type Heightmap struct {
	MotionBlocking []int64 `nbt:"MOTION_BLOCKING"`
}

func (c *ChunkColumn) highestBlock(x int, z int) int {
	for y := Height - 1; y >= 0; y-- {
		if c.BlockAt(x, y, z) != blocks.Air {
			return y
		}
	}
	return 0
}

func (c *ChunkColumn) HeightMap() Heightmap {
	highestBlocks := make([]BlockState, Width*Depth)
	for x := 0; x < Width; x++ {
		for z := 0; z < Depth; z++ {
			highestBlocks[x+z*Width] = BlockState(c.highestBlock(x, z))
		}
	}
	return Heightmap{
		MotionBlocking: NewPaddedBlockArray(9, highestBlocks).Data,
	}
}
