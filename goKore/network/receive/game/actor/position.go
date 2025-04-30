package actor

import (
	"github.com/lenaxia/goKore/network/receive/game/field"
)

// Position is an alias for field.Position
type Position = field.Position

// NewPosition creates a new Position with the given coordinates
func NewPosition(x, y int) *Position {
	return &Position{X: x, Y: y}
}
