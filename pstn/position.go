package pstn

type Chunk struct {
	X, Z int32
}

type Block struct {
	X, Y, Z int32
}

type Entity struct {
	X, Y, Z float64
}

const (
	chunkSize = 16
)

func BlockToChunk(pos Block) Chunk {
	return Chunk{X: pos.X / chunkSize, Z: pos.Z / chunkSize}
}

func EntityToChunk(pos Entity) Chunk {
	return Chunk{X: int32(pos.X) / 16, Z: int32(pos.Z) / 16}
}

func BlockToEntity(pos Block) Entity {
	return Entity{X: float64(pos.X), Y: float64(pos.Y), Z: float64(pos.Z)}
}
