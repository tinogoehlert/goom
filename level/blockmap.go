package level

import (
	"github.com/tinogoehlert/goom/utils"
)

// BlockList list of LineDefs within the Block
type BlockList struct {
}

// BlockMap is simply a grid of "blocks"' each 128×128 units
type BlockMap struct {
	Origin utils.Vec2
}
