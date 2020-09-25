package proto

import (
	"bytes"
	"fmt"
	"testing"
)

type varintTestCase struct {
	bytes []byte
	value int64
}

var varint32Cases = []varintTestCase{
	{[]byte{0x00}, 0},
	{[]byte{0x01}, 1},
	{[]byte{0x02}, 2},
	{[]byte{0x7f}, 127},
	{[]byte{0x80, 0x01}, 128},
	{[]byte{0xff, 0x01}, 255},
	{[]byte{0xff, 0xff, 0xff, 0xff, 0x07}, 2147483647},
	{[]byte{0xff, 0xff, 0xff, 0xff, 0x0f}, -1},
	{[]byte{0x80, 0x80, 0x80, 0x80, 0x08}, -2147483648},
}

func TestPacketReader_ReadVar32(t *testing.T) {
	for _, cs := range varint32Cases {
		pd := NewPacketDecoder(bytes.NewBuffer(cs.bytes))
		if actual := pd.ReadVar32(); actual != int32(cs.value) {
			t.Errorf("ReadVar32(%v): Expected %d, got %d", cs.bytes, actual, cs.value)
		}
	}
}

var varint64Cases = []varintTestCase{
	{[]byte{0x00}, 0},
	{[]byte{0x01}, 1},
	{[]byte{0x02}, 2},
	{[]byte{0x7f}, 127},
	{[]byte{0x80, 0x01}, 128},
	{[]byte{0xff, 0x01}, 255},
	{[]byte{0xff, 0xff, 0xff, 0xff, 0x07}, 2147483647},
	{[]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x7f}, 9223372036854775807},
	{[]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x01}, -1},
	{[]byte{0x80, 0x80, 0x80, 0x80, 0xf8, 0xff, 0xff, 0xff, 0xff, 0x01}, -2147483648},
	{[]byte{0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x01}, -9223372036854775808},
}

func TestPacketReader_ReadVar64(t *testing.T) {
	for _, cs := range varint64Cases {
		pd := NewPacketDecoder(bytes.NewBuffer(cs.bytes))
		if actual := pd.ReadVar64(); actual != cs.value {
			t.Errorf("ReadVar64(%v): Expected %d, got %d", cs.bytes, actual, cs.value)
		}
	}
}

func TestPacketDecoder_ReadPacket(t *testing.T) {
	packetData := []byte{0x03, 0x01, 0x00, 0x50}
	d := NewPacketDecoder(bytes.NewBuffer(packetData))
	packet := d.ReadPacket()
	if packet.ID != 1 {
		t.Errorf("RecvPacket ID: Expected 1, got %d", packet.ID)
	}

	if value := packet.ReadU16(); value != 0x50 {
		t.Errorf("RecvPacket Value: Expected 80, got %d", value)
	}
}

func TestPacketDecoder_ReadString(t *testing.T) {
	tests := []struct {
		want string
	}{
		{"test♠"},
		{"another random string"},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("ReadString(%s)", tt.want), func(t *testing.T) {
			var b bytes.Buffer
			buf := NewEncoder(&b)
			binStr := []byte(tt.want)
			buf.WriteVar32(int32(len(binStr)))
			buf.Write(binStr)
			p := &PacketDecoder{&b}
			if got := p.ReadString(); got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
