package proto

import (
	"bytes"
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
	"github.com/masp/mcgo/pstn"
	"math"
	"reflect"
	"testing"
)

func initOps(b *bytes.Buffer) (*PacketEncoder, *PacketDecoder) {
	return NewEncoder(b), NewPacketDecoder(b)
}

func TestVarints(t *testing.T) {
	properties := gopter.NewProperties(nil)

	properties.Property("decode_var64(encode_var64(v)) = v", prop.ForAll(
		func(v int64) bool {
			var b bytes.Buffer
			enc, dec := initOps(&b)

			enc.WriteVar64(v)
			return dec.ReadVar64() == v
		},
		gen.Int64Range(math.MinInt64, math.MaxInt64),
	))

	properties.Property("decode_var32(encode_var32(v)) = v", prop.ForAll(
		func(v int32) bool {
			var b bytes.Buffer
			enc, dec := initOps(&b)

			enc.WriteVar32(v)
			return dec.ReadVar32() == v
		},
		gen.Int32Range(math.MinInt32, math.MaxInt32),
	))

	properties.TestingRun(t)
}

func TestStrings(t *testing.T) {
	properties := gopter.NewProperties(nil)

	properties.Property("decode_string(encode_string(v)) = v", prop.ForAll(
		func(s string) bool {
			var b bytes.Buffer
			enc, dec := initOps(&b)

			enc.WriteString(s)
			return dec.ReadString() == s
		},
		gen.AnyString(),
	))

	properties.TestingRun(t)
}

func genPosition() gopter.Gen {
	return gen.Struct(reflect.TypeOf(&pstn.Block{}), map[string]gopter.Gen{
		"X": gen.Int32Range(-30000000, 30000000),
		"Y": gen.Int32Range(-256, 256),
		"Z": gen.Int32Range(-30000000, 30000000),
	})
}

func TestPositions(t *testing.T) {
	properties := gopter.NewProperties(nil)

	properties.Property("decode_pos(encode_pos(v)) = v", prop.ForAll(
		func(p pstn.Block) bool {
			var b bytes.Buffer
			enc, dec := initOps(&b)

			enc.WritePosition(p)
			return dec.ReadPosition() == p
		},
		genPosition(),
	))

	properties.TestingRun(t)
}
