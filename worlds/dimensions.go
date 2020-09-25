package worlds

import (
	"bytes"
	"github.com/masp/mcgo/blocks"
	"github.com/masp/mcgo/chunks"
	"github.com/masp/mcgo/pstn"
)

type Type int

const (
	Nether    Type = -1
	Overworld Type = 0
	End       Type = 1
)

type Dimension struct {
	Spawn  pstn.Block
	chunks map[pstn.Chunk]*chunks.ChunkColumn
}

func New(spawn pstn.Block) Dimension {
	w := Dimension{
		Spawn:  spawn,
		chunks: make(map[pstn.Chunk]*chunks.ChunkColumn),
	}
	w.generateSpawn()
	return w
}

func (w *Dimension) ChunkAt(p pstn.Chunk) *chunks.ChunkColumn {
	return w.chunks[p]
}

func (w *Dimension) ChunkAtBlock(p pstn.Block) *chunks.ChunkColumn {
	return w.chunks[pstn.BlockToChunk(p)]
}

func (w *Dimension) LoadChunk(chunk *chunks.ChunkColumn) {
	w.chunks[chunk.Pos] = chunk
}

func generateFlatChunk(p pstn.Chunk) *chunks.ChunkColumn {
	chunk := chunks.NewChunk(p)
	for y := 0; y < chunks.Height; y++ {
		for x := 0; x < chunks.Width; x++ {
			for z := 0; z < chunks.Depth; z++ {
				if y < 63 {
					chunk.SetBlockAt(x, y, z, blocks.Stone)
				}
			}
		}
	}
	return chunk
}

const (
	SpawnSize = 16 // chunks
)

func (w *Dimension) generateSpawn() {
	spawnChunk := pstn.BlockToChunk(w.Spawn)
	for x := spawnChunk.X - SpawnSize; x <= spawnChunk.X+SpawnSize; x++ {
		for z := spawnChunk.Z - SpawnSize; z <= spawnChunk.Z+SpawnSize; z++ {
			p := pstn.Chunk{X: x, Z: z}
			w.LoadChunk(generateFlatChunk(p))
		}
	}
	w.generateLighting()
}

func (w *Dimension) generateLighting() {
	allLight := bytes.Repeat([]byte{15}, chunks.BlocksInSection)
	for _, chunk := range w.chunks {
		for _, sec := range chunk.Sections {
			sec.Lighting.SkyLight.Init(allLight)
		}
		chunk.SkySection.SkyLight.Init(allLight)
	}
}
