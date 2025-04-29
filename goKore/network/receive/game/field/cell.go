package field

import (
	"fmt"
)

// Cell represents a single cell in a field
type Cell struct {
	// The field this cell belongs to
	field *Field

	// X coordinate
	x int

	// Y coordinate
	y int

	// Cell type (walkable, water, etc.)
	cellType byte
}

// NewCell creates a new Cell instance
func NewCell(field *Field, x, y int) *Cell {
	if field == nil || field.IsOffMap(x, y) {
		return nil
	}

	return &Cell{
		field:    field,
		x:        x,
		y:        y,
		cellType: field.GetBlock(x, y),
	}
}

// X returns the X coordinate of the cell
func (c *Cell) X() int {
	return c.x
}

// Y returns the Y coordinate of the cell
func (c *Cell) Y() int {
	return c.y
}

// Position returns the position of the cell
func (c *Cell) Position() Position {
	return Position{X: c.x, Y: c.y}
}

// Type returns the cell type
func (c *Cell) Type() byte {
	return c.cellType
}

// IsWalkable checks if the cell is walkable
func (c *Cell) IsWalkable() bool {
	return (c.cellType & TileWalk) != 0
}

// IsSnipable checks if the cell is snipable
func (c *Cell) IsSnipable() bool {
	return (c.cellType & TileSnipe) != 0
}

// IsWater checks if the cell is water
func (c *Cell) IsWater() bool {
	return (c.cellType & TileWater) != 0
}

// IsCliff checks if the cell is a cliff
func (c *Cell) IsCliff() bool {
	return (c.cellType & TileCliff) != 0
}

// Weight returns the weight of the cell for pathfinding
func (c *Cell) Weight() byte {
	return c.field.GetBlockWeight(c.x, c.y)
}

// DistanceTo calculates the Manhattan distance to another cell
func (c *Cell) DistanceTo(other *Cell) int {
	return abs(c.x-other.x) + abs(c.y-other.y)
}

// DistanceToPosition calculates the Manhattan distance to a position
func (c *Cell) DistanceToPosition(pos Position) int {
	return abs(c.x-pos.X) + abs(c.y-pos.Y)
}

// Neighbors returns the walkable neighboring cells
func (c *Cell) Neighbors() []*Cell {
	var neighbors []*Cell

	// Check the four adjacent cells (up, right, down, left)
	directions := []struct{ dx, dy int }{
		{0, -1}, // Up
		{1, 0},  // Right
		{0, 1},  // Down
		{-1, 0}, // Left
	}

	for _, dir := range directions {
		nx, ny := c.x+dir.dx, c.y+dir.dy

		// Skip if out of bounds
		if c.field.IsOffMap(nx, ny) {
			continue
		}

		// Skip if not walkable
		if !c.field.IsWalkable(nx, ny) {
			continue
		}

		neighbors = append(neighbors, NewCell(c.field, nx, ny))
	}

	return neighbors
}

// AllNeighbors returns all neighboring cells (including diagonals)
func (c *Cell) AllNeighbors() []*Cell {
	var neighbors []*Cell

	// Check all eight surrounding cells
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			// Skip the cell itself
			if dx == 0 && dy == 0 {
				continue
			}

			nx, ny := c.x+dx, c.y+dy

			// Skip if out of bounds
			if c.field.IsOffMap(nx, ny) {
				continue
			}

			neighbors = append(neighbors, NewCell(c.field, nx, ny))
		}
	}

	return neighbors
}

// WalkableNeighbors returns all walkable neighboring cells (including diagonals)
func (c *Cell) WalkableNeighbors() []*Cell {
	var neighbors []*Cell

	// Check all eight surrounding cells
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			// Skip the cell itself
			if dx == 0 && dy == 0 {
				continue
			}

			nx, ny := c.x+dx, c.y+dy

			// Skip if out of bounds
			if c.field.IsOffMap(nx, ny) {
				continue
			}

			// Skip if not walkable
			if !c.field.IsWalkable(nx, ny) {
				continue
			}

			// For diagonal movement, check that the path is not blocked
			if dx != 0 && dy != 0 {
				// Check that we can move horizontally and vertically
				if !c.field.IsWalkable(c.x+dx, c.y) || !c.field.IsWalkable(c.x, c.y+dy) {
					continue
				}
			}

			neighbors = append(neighbors, NewCell(c.field, nx, ny))
		}
	}

	return neighbors
}

// CanMoveTo checks if it's possible to move from this cell to another cell
func (c *Cell) CanMoveTo(other *Cell) bool {
	// Check if the cells are in the same field
	if c.field != other.field {
		return false
	}

	// Check if both cells are walkable
	if !c.IsWalkable() || !other.IsWalkable() {
		return false
	}

	// Calculate Manhattan distance
	dx := abs(c.x - other.x)
	dy := abs(c.y - other.y)

	// If the cells are adjacent, we can move directly
	if (dx == 1 && dy == 0) || (dx == 0 && dy == 1) {
		return true
	}

	// If the cells are diagonal, check that the path is not blocked
	if dx == 1 && dy == 1 {
		return c.field.IsWalkable(c.x+dx, c.y) && c.field.IsWalkable(c.x, c.y+dy)
	}

	// For longer distances, we need to check the path
	return c.field.CheckLOS(c.Position(), other.Position(), false)
}

// HasLineOfSight checks if there is a line of sight to another cell
func (c *Cell) HasLineOfSight(other *Cell, canSnipe bool) bool {
	// Check if the cells are in the same field
	if c.field != other.field {
		return false
	}

	return c.field.CheckLOS(c.Position(), other.Position(), canSnipe)
}

// String returns a string representation of the cell
func (c *Cell) String() string {
	var typeStr string

	if c.IsWalkable() {
		typeStr = "walkable"
	} else {
		typeStr = "not walkable"
	}

	if c.IsWater() {
		typeStr += ", water"
	}

	if c.IsSnipable() {
		typeStr += ", snipable"
	}

	if c.IsCliff() {
		typeStr += ", cliff"
	}

	return fmt.Sprintf("Cell(%d,%d): %s, weight=%d", c.x, c.y, typeStr, c.Weight())
}

// CellGrid represents a grid of cells
type CellGrid struct {
	field  *Field
	cells  [][]*Cell
	width  int
	height int
}

// NewCellGrid creates a new CellGrid for a field
func NewCellGrid(field *Field) *CellGrid {
	width := field.Width()
	height := field.Height()

	// Initialize the grid
	cells := make([][]*Cell, height)
	for y := 0; y < height; y++ {
		cells[y] = make([]*Cell, width)
		for x := 0; x < width; x++ {
			cells[y][x] = NewCell(field, x, y)
		}
	}

	return &CellGrid{
		field:  field,
		cells:  cells,
		width:  width,
		height: height,
	}
}

// GetCell returns the cell at the specified coordinates
func (g *CellGrid) GetCell(x, y int) *Cell {
	if x < 0 || x >= g.width || y < 0 || y >= g.height {
		return nil
	}

	return g.cells[y][x]
}

// GetCellAtPosition returns the cell at the specified position
func (g *CellGrid) GetCellAtPosition(pos Position) *Cell {
	return g.GetCell(pos.X, pos.Y)
}

// Width returns the width of the grid
func (g *CellGrid) Width() int {
	return g.width
}

// Height returns the height of the grid
func (g *CellGrid) Height() int {
	return g.height
}

// Field returns the field associated with this grid
func (g *CellGrid) Field() *Field {
	return g.field
}

// FindPath finds a path between two cells using A* algorithm
// This is a placeholder - the actual implementation will be in path.go
func (g *CellGrid) FindPath(start, end *Cell) []*Cell {
	// This will be implemented in path.go
	return nil
}
