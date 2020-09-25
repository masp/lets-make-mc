package chunks

import (
	"bytes"
	"reflect"
	"testing"
)

func TestNewLightingArray(t *testing.T) {
	var data [2048]byte
	copy(data[:], bytes.Repeat([]byte{238}, 2048))

	type args struct {
		values []byte
	}
	tests := []struct {
		name string
		args args
		want LightingArray
	}{
		{"allset", args{bytes.Repeat([]byte{14}, 4096)}, LightingArray{Data: data}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewLightingArray(tt.args.values); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewLightingArray() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewPaddedBlockArray(t *testing.T) {
	type args struct {
		bitsPerItem int
		values      []BlockState
	}
	tests := []struct {
		name string
		args args
		want []int64
	}{
		{"5bits",
			args{5, []BlockState{1, 2, 2, 3, 4, 4, 5, 6, 6, 4, 8, 0, 7, 4, 3, 13, 15, 16, 9, 14, 10, 12, 0, 2}},
			[]int64{0x0020863148418841, 0x01018A7260F68C87},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPaddedBlockArray(tt.args.bitsPerItem, tt.args.values); !reflect.DeepEqual(got.Data, tt.want) {
				t.Errorf("CreatePackedArray() = %064b, want %064b", got.Data, tt.want)
			}
		})
	}
	// 0000000000100000100001100011000101001000010000011000100001000001
	// 0000000000100000100001100011000101001000010000011000100001000001

	// 0000000100000001100010100111001001100000111101101000110010000111
	// 0000000100000001100010100111001001100000111101101000110010000111
}
