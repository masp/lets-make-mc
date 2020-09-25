package chunks

import (
	"github.com/masp/mcgo/blocks"
	"github.com/masp/mcgo/pstn"
	"testing"
)

func TestChunkSection_SetBlockAt(t *testing.T) {
	c := NewChunk(pstn.Chunk{})
	for s, section := range c.Sections {
		if !section.IsAir() {
			t.Errorf("got section (%d) IsAir() = false, want IsAir() = true", s)
		}
	}

	c.SetBlockAt(0, 0, 0, blocks.Stone)
	if c.Sections[0].IsAir() {
		t.Errorf("got section 0 IsAir() = true, want IsAir() = false")
	}
	c.SetBlockAt(0, 0, 0, blocks.Air)
	if !c.Sections[0].IsAir() {
		t.Errorf("got section 0 IsAir() = false after block change, want IsAir() = true")
	}
}
