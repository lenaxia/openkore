package main

import (
	"fmt"

	"github.com/lenaxia/goKore/network/receive/game/field"
)

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	fmt.Println("Field test application")
	fmt.Println("=====================")
	fmt.Println("This is a simple demonstration of the field package.")
	fmt.Println()

	// Create a simple test field
	width := 20
	height := 20
	rawMap := make([]byte, width*height)

	// Make all cells walkable
	for i := 0; i < width*height; i++ {
		rawMap[i] = field.TileWalk
	}

	// Add a wall with a gap
	for y := 5; y < 15; y++ {
		if y != 10 {
			rawMap[y*width+10] = 0 // Not walkable
		}
	}

	// Create the field
	fmt.Println("Creating test field...")
	testField := &field.Field{}

	// Use exported methods to set field properties
	// In a real application, you would use New() or LoadFile()
	testField.SetName("test_field")
	testField.SetWidth(width)
	testField.SetHeight(height)
	testField.SetRawMap(rawMap)

	// Print information about the field
	fmt.Printf("Field name: %s\n", testField.Name())
	fmt.Printf("Field dimensions: %d x %d\n", testField.Width(), testField.Height())
	fmt.Println()

	// Print information about some cells
	fmt.Println("Cell information:")
	fmt.Println("----------------")

	// Check walkable cells
	fmt.Printf("Cell (5,5) is walkable: %v\n", testField.IsWalkable(5, 5))
	fmt.Printf("Cell (10,10) is walkable: %v\n", testField.IsWalkable(10, 10))
	fmt.Printf("Cell (10,5) is walkable: %v\n", testField.IsWalkable(10, 5))

	// Create a cell grid
	fmt.Println("\nCreating cell grid...")
	grid := field.NewCellGrid(testField)

	// Get some cells
	cell1 := grid.GetCell(5, 5)
	cell2 := grid.GetCell(15, 15)

	if cell1 != nil && cell2 != nil {
		fmt.Printf("Distance from cell (5,5) to cell (15,15): %d\n", cell1.DistanceTo(cell2))
	}

	// Find a path
	fmt.Println("\nPathfinding:")
	fmt.Println("------------")

	// Create a pathfinder
	pathFinder := field.NewPathFinder(grid)

	// Find a path from (5,5) to (15,15)
	start := field.Position{X: 5, Y: 5}
	end := field.Position{X: 15, Y: 15}

	path := pathFinder.FindPath(start, end)

	if path != nil {
		fmt.Printf("Found path from (%d,%d) to (%d,%d) with %d steps\n",
			start.X, start.Y, end.X, end.Y, path.Length())

		// Print the first few steps
		positions := path.Positions()
		fmt.Println("First few steps:")
		for i := 0; i < min(5, len(positions)); i++ {
			fmt.Printf("  Step %d: (%d,%d)\n", i+1, positions[i].X, positions[i].Y)
		}
	} else {
		fmt.Printf("No path found from (%d,%d) to (%d,%d)\n", start.X, start.Y, end.X, end.Y)
	}

	fmt.Println()

	fmt.Println("Press Enter to exit...")
	fmt.Scanln()
}
