package protocol

import (
	"bufio"
	"io"
)

// PacketDecoder is a buffered reader that can frame and return individual packets from any Read interface
type PacketDecoder struct {
	rd io.Reader
}

// NewPacketReader returns a reader able to frame incoming byte streams from the reader into individual packets
func NewPacketReader(rd io.Reader) *PacketDecoder {
	p := PacketDecoder{
		rd: bufio.NewReaderSize(rd, 4096),
	}
	return &p
}

const (
	PacketHandshake = 0x00
)

// RawPacket contains the numeric ID (normalized to match server version IDs) as well as
// a slice pointing to the binary data of the packet
type RawPacket struct {
	ID   uint
	Data []byte
}

func (p *PacketDecoder) nextVarint(max int) int64 {
	rd := p.rd.(io.ByteReader) // always assume we are using a buffered reader
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

func (p *PacketDecoder) NextVar32() int32 {
	return int32(p.nextVarint(5))
}

func (p *PacketDecoder) NextVar64() int64 {
	return p.nextVarint(10)
}

func (p *PacketDecoder) ReadPacket() RawPacket {
	len := p.NextVar32()
	if len <= 0 {
		panic("invalid packet sizing (must be > 0)")
	}

	id := p.NextVar32()
	if id < 0 {
		panic("invalid packet ID (must be >= 0)")
	}

	packet := RawPacket {
		ID: uint(id),
		Data: make([]byte, len),
	}

	_, err := io.ReadFull(p.rd, packet.Data)
	if err != nil {
		panic(err)
	}
	return packet
}
