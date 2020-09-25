package proto

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/Tnze/go-mc/nbt"
	"github.com/masp/mcgo/pstn"
	uuid "github.com/satori/go.uuid"
	"io"
)

type PacketEncoder struct {
	io.Writer
}

func NewEncoder(w io.Writer) *PacketEncoder {
	if _, ok := w.(io.ByteWriter); !ok { // not a buffered writer
		w = bufio.NewWriter(w)
	}
	return &PacketEncoder{w}
}

type EncodableAsPacket interface {
	// Encode writes the packet contents to the packet encoder. If packet is invalid,
	// a receiver should panic with an error.
	EncodeTo(e *PacketEncoder)
}

type EncodeFunc func(*PacketEncoder)
type SimplePacket struct {
	encoderFunc EncodeFunc
}

func (s SimplePacket) EncodeTo(e *PacketEncoder) { s.encoderFunc(e) }

func NewPacket(encoderFunc EncodeFunc) SimplePacket {
	return SimplePacket{
		encoderFunc: encoderFunc,
	}
}

func mustWrite(_ int, err error) {
	must(err)
}

func must(err error) {
	if err != nil {
		panic(fmt.Errorf("packet encoding fialed: %w", err))
	}
}

// EncodePacket takes an encodable packet struct and creates a id + packet data byte array that can then be sent.
// It's good to encode the byte array before writing to the actual socket, so we can treat it like arbitrary data
// and not have to worry about the lifetime of the original source data at all.
func EncodePacket(id PacketID, p EncodableAsPacket) []byte {
	// We need to get the length of the packet before writing it to the buffer
	var b bytes.Buffer
	tmpEnc := NewEncoder(&b)
	tmpEnc.WriteVar32(int32(id))
	p.EncodeTo(tmpEnc)
	return b.Bytes()
}

// Flush forces the underlying buffer that the encoder is using to be forced to the writer. If there's no underlying
// writer, this is a no-op.
func (e *PacketEncoder) Flush() {
	if writer, ok := e.Writer.(*bufio.Writer); ok {
		must(writer.Flush())
	}
}

func (e *PacketEncoder) WriteVar32(v int32) {
	wr := e.Writer.(io.ByteWriter)
	for {
		temp := v & 0x7F
		v = int32(uint32(v) >> 7)

		if v != 0 {
			temp |= 0x80
		}

		must(wr.WriteByte(byte(temp)))

		if v == 0 {
			break
		}
	}
}

func (e *PacketEncoder) WriteVar64(v int64) {
	wr := e.Writer.(io.ByteWriter)
	for {
		temp := v & 0x7F
		v = int64(uint64(v) >> 7)

		if v != 0 {
			temp |= 0x80
		}

		must(wr.WriteByte(byte(temp)))

		if v == 0 {
			break
		}
	}
}

func (e *PacketEncoder) WriteString(s string) {
	b := []byte(s)
	e.WriteVar32(int32(len(b)))
	mustWrite(e.Write(b))
}

func (e *PacketEncoder) WriteBytes(p []byte) {
	mustWrite(e.Write(p))
}

func (e *PacketEncoder) WritePosition(p pstn.Block) {
	position := ((uint64(p.X) & 0x3FFFFFF) << 38) | ((uint64(p.Z) & 0x3FFFFFF) << 12) | (uint64(p.Y) & 0xFFF)
	e.WriteI64(int64(position))
}

func (e *PacketEncoder) WriteUUID(uuid uuid.UUID) {
	mustWrite(e.Write(uuid.Bytes()))
}

func (e *PacketEncoder) WriteNBT(data interface{}) {
	must(nbt.Marshal(e.Writer, data))
}

func (e *PacketEncoder) mustWriteBE(data interface{}) {
	_ = binary.Write(e.Writer, binary.BigEndian, data)
}

func (e *PacketEncoder) WriteBool(v bool) {
	if v {
		e.mustWriteBE(byte(0x01))
	} else {
		e.mustWriteBE(byte(0x00))
	}
}
func (e *PacketEncoder) WriteU8(v uint8)        { e.mustWriteBE(v) }
func (e *PacketEncoder) WriteU16(v uint16)      { e.mustWriteBE(v) }
func (e *PacketEncoder) WriteU32(v uint32)      { e.mustWriteBE(v) }
func (e *PacketEncoder) WriteU64(v uint64)      { e.mustWriteBE(v) }
func (e *PacketEncoder) WriteI8(v int8)         { e.mustWriteBE(v) }
func (e *PacketEncoder) WriteI16(v int16)       { e.mustWriteBE(v) }
func (e *PacketEncoder) WriteI32(v int32)       { e.mustWriteBE(v) }
func (e *PacketEncoder) WriteI64(v int64)       { e.mustWriteBE(v) }
func (e *PacketEncoder) WriteFloat32(v float32) { e.mustWriteBE(v) }
func (e *PacketEncoder) WriteFloat64(v float64) { e.mustWriteBE(v) }

func (e *PacketEncoder) WritePacket(packet []byte) {
	e.WriteVar32(int32(len(packet)))
	mustWrite(e.Write(packet))
}
