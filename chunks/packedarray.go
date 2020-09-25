package chunks

import (
	"fmt"
	"github.com/masp/mcgo/proto"
)

type PaddedBlockArray struct {
	bitsPerItem int
	Data        []int64
}

// ceildiv divides two integers and if the division is not even adds 1 to ceil the result
func ceildiv(numerator, denominator int) int {
	quotient := numerator / denominator
	isEven := numerator%denominator == 0
	if !isEven {
		quotient++
	}
	return quotient
}

func NewPaddedBlockArray(bitsPerItem int, values []BlockState) PaddedBlockArray {
	if bitsPerItem >= 64 || bitsPerItem == 0 {
		panic("Invalid bits per item, must be between [1-64]")
	}

	itemsPerLong := 64 / bitsPerItem
	longsNeeded := ceildiv(len(values), itemsPerLong)
	array := PaddedBlockArray{
		bitsPerItem: bitsPerItem,
		Data:        make([]int64, longsNeeded),
	}

	for i, v := range values {
		longIndex := i / itemsPerLong
		bitPosStart := (i % itemsPerLong) * bitsPerItem
		array.Data[longIndex] |= int64(v) << bitPosStart
	}
	return array
}

type LightingArray struct {
	// Even blocks are stored in first nibble
	// Odd blocks are stored in second nibble
	Data [2048]byte
}

func (l *LightingArray) EncodeTo(e *proto.PacketEncoder) {
	e.WriteVar32(int32(len(l.Data)))
	e.WriteBytes(l.Data[:])
}

func NewLightingArray(data []byte) LightingArray {
	array := LightingArray{}
	array.Init(data)
	return array
}

func (arr *LightingArray) Init(data []byte) {
	if len(data) != BlocksInSection {
		panic("invalid Lighting array size (must be 4096)")
	}

	for i, v := range data {
		if v > 15 || v < 0 {
			panic(fmt.Errorf("invalid Lighting value %d", v))
		}
		if i%2 == 0 {
			arr.Data[i/2] |= v
		} else {
			arr.Data[i/2] |= v << 4
		}
	}
}
