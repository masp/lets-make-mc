package proto

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/masp/mcgo/pstn"
	uuid "github.com/satori/go.uuid"
	"io"
)

type PacketID int32

// PacketDecoder is a buffered reader that can frame and return individual packets from any Read interface
type PacketDecoder struct {
	io.Reader
}

// NewPacketDecoder returns a reader able to frame incoming byte streams from the reader into individual packets
func NewPacketDecoder(rd io.Reader) *PacketDecoder {
	p := PacketDecoder{
		bufio.NewReaderSize(rd, 4096),
	}
	return &p
}

// RecvPacket contains the numeric ID (normalized to match server version IDs) as well as
// a slice pointing to the binary data of the packet
type RecvPacket struct {
	PacketDecoder
	ID PacketID
}

func (p *PacketDecoder) ReadPacket() RecvPacket {
	length := p.ReadVar32()
	if length <= 0 {
		panic(fmt.Errorf("invalid packet sizing %d (must be > 0)", length))
	}

	payload := make([]byte, length)
	_, err := p.Read(payload)
	if err != nil {
		panic(fmt.Errorf("failed reading %d bytes for packet payload: %w", length, err))
	}
	payloadDecoder := PacketDecoder{bytes.NewReader(payload)}

	id := PacketID(payloadDecoder.ReadVar32())
	if id < 0 {
		panic(fmt.Errorf("invalid packet ID %d (must be >= 0)", id))
	}

	packet := RecvPacket{
		PacketDecoder: payloadDecoder,
		ID:            id,
	}
	return packet
}

func (p *PacketDecoder) nextVarint(max int) int64 {
	rd := p.Reader.(io.ByteReader) // always assume we are using a buffered reader
	var num int
	var res int64

	for {
		tmp, err := rd.ReadByte()
		if err != nil {
			panic(err)
		}
		res |= (int64(tmp) & 0x7F) << uint(num*7)

		if num++; num > max {
			panic("Invalid varint: value too big")
		}

		if tmp&0x80 != 0x80 {
			break
		}
	}

	return res
}

func (p *PacketDecoder) ReadVar32() int32 {
	return int32(p.nextVarint(5))
}

func (p *PacketDecoder) ReadVar64() int64 {
	return p.nextVarint(10)
}

func (p *PacketDecoder) ReadString() string {
	size := p.ReadVar32()
	strBytes := make([]byte, size)
	_, err := p.Read(strBytes)
	if err != nil {
		panic(fmt.Errorf("failed to read string from packet: %w", err))
	}
	return string(strBytes)
}

func (p *PacketDecoder) ReadPosition() pstn.Block {
	v := p.ReadI64()
	x := int32(v >> 38)
	y := int32(v & 0xFFF)
	z := int32(v << 26 >> 38)

	if x >= 1<<25 {
		x -= 1 << 26
	}
	if y >= 1<<11 {
		y -= 1 << 12
	}
	if z >= 1<<25 {
		z -= 1 << 26
	}
	return pstn.Block{X: x, Y: y, Z: z}
}

func (p *PacketDecoder) ReadUUID() uuid.UUID {
	bs := make([]byte, 16)
	_, err := p.Read(bs)
	if err != nil {
		panic(err)
	}
	return uuid.FromBytesOrNil(bs)
}

func (p *PacketDecoder) mustRead(data interface{}) {
	err := binary.Read(p, binary.BigEndian, data)
	if err != nil {
		panic(fmt.Errorf("unexpected error reading stream: %w", err))
	}
}

func (p *PacketDecoder) ReadBool() bool {
	return p.ReadU8() == 1
}

func (p *PacketDecoder) ReadU8() uint8 {
	var res uint8
	p.mustRead(&res)
	return res
}

func (p *PacketDecoder) ReadI8() int8 {
	var res int8
	p.mustRead(&res)
	return res
}

func (p *PacketDecoder) ReadU16() uint16 {
	var res uint16
	p.mustRead(&res)
	return res
}

func (p *PacketDecoder) ReadI16() int16 {
	var res int16
	p.mustRead(&res)
	return res
}

func (p *PacketDecoder) ReadU32() uint32 {
	var res uint32
	p.mustRead(&res)
	return res
}

func (p *PacketDecoder) ReadI32() int32 {
	var res int32
	p.mustRead(&res)
	return res
}

func (p *PacketDecoder) ReadU64() uint64 {
	var res uint64
	p.mustRead(&res)
	return res
}

func (p *PacketDecoder) ReadI64() int64 {
	var res int64
	p.mustRead(&res)
	return res
}

func (p *PacketDecoder) ReadFloat32() float32 {
	var res float32
	p.mustRead(&res)
	return res
}

func (p *PacketDecoder) ReadFloat64() float64 {
	var res float64
	p.mustRead(&res)
	return res
}
