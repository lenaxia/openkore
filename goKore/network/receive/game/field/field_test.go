package field

import (
	"testing"
)

// TestFieldCreation tests the creation of a Field
func TestFieldCreation(t *testing.T) {
	// Create a simple test field
	width := 10
	height := 10
	rawMap := make([]byte, width*height)

	// Set some cells as walkable
	for i := 0; i < width*height; i++ {
		if i%2 == 0 {
			rawMap[i] = TileWalk
		}
	}

	field := &Field{
		name:     "test_field",
		baseName: "test_field",
		width:    width,
		height:   height,
		rawMap:   rawMap,
	}

	// Test field properties
	if field.Name() != "test_field" {
		t.Errorf("Expected name to be 'test_field', got '%s'", field.Name())
	}

	if field.Width() != width {
		t.Errorf("Expected width to be %d, got %d", width, field.Width())
	}

	if field.Height() != height {
		t.Errorf("Expected height to be %d, got %d", height, field.Height())
	}
}

// TestWalkableChecks tests the walkable cell checks
func TestWalkableChecks(t *testing.T) {
	// Create a simple test field
	width := 10
	height := 10
	rawMap := make([]byte, width*height)

	// Set some cells as walkable
	for i := 0; i < width*height; i++ {
		if i%2 == 0 {
			rawMap[i] = TileWalk
		}
	}

	field := &Field{
		name:     "test_field",
		baseName: "test_field",
		width:    width,
		height:   height,
		rawMap:   rawMap,
	}

	// Test walkable cells
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			offset := y*width + x
			expected := rawMap[offset] == TileWalk
			if field.IsWalkable(x, y) != expected {
				t.Errorf("IsWalkable(%d, %d) = %v, expected %v", x, y, field.IsWalkable(x, y), expected)
			}
		}
	}
}

// TestOffMapChecks tests the off-map checks
func TestOffMapChecks(t *testing.T) {
	// Create a simple test field
	width := 10
	height := 10
	rawMap := make([]byte, width*height)

	field := &Field{
		name:     "test_field",
		baseName: "test_field",
		width:    width,
		height:   height,
		rawMap:   rawMap,
	}

	// Test off-map coordinates
	testCases := []struct {
		x, y     int
		expected bool
	}{
		{-1, 5, true},  // Left edge
		{10, 5, true},  // Right edge
		{5, -1, true},  // Top edge
		{5, 10, true},  // Bottom edge
		{5, 5, false},  // Inside
		{0, 0, false},  // Top-left corner
		{9, 9, false},  // Bottom-right corner
		{-1, -1, true}, // Outside top-left
		{10, 10, true}, // Outside bottom-right
	}

	for _, tc := range testCases {
		if field.IsOffMap(tc.x, tc.y) != tc.expected {
			t.Errorf("IsOffMap(%d, %d) = %v, expected %v", tc.x, tc.y, field.IsOffMap(tc.x, tc.y), tc.expected)
		}
	}
}

// TestCellCreation tests the creation of Cell objects
func TestCellCreation(t *testing.T) {
	// Create a simple test field
	width := 10
	height := 10
	rawMap := make([]byte, width*height)

	// Set some cells as walkable
	for i := 0; i < width*height; i++ {
		if i%2 == 0 {
			rawMap[i] = TileWalk
		}
	}

	field := &Field{
		name:     "test_field",
		baseName: "test_field",
		width:    width,
		height:   height,
		rawMap:   rawMap,
	}

	// Test cell creation
	cell := NewCell(field, 5, 5)
	if cell == nil {
		t.Fatal("Expected cell to be created, got nil")
	}

	if cell.X() != 5 {
		t.Errorf("Expected X to be 5, got %d", cell.X())
	}

	if cell.Y() != 5 {
		t.Errorf("Expected Y to be 5, got %d", cell.Y())
	}

	// Test cell properties
	expected := (rawMap[5*width+5] & TileWalk) != 0
	if cell.IsWalkable() != expected {
		t.Errorf("IsWalkable() = %v, expected %v", cell.IsWalkable(), expected)
	}
}

// TestCellGrid tests the CellGrid functionality
func TestCellGrid(t *testing.T) {
	// Create a simple test field
	width := 10
	height := 10
	rawMap := make([]byte, width*height)

	// Set some cells as walkable
	for i := 0; i < width*height; i++ {
		if i%2 == 0 {
			rawMap[i] = TileWalk
		}
	}

	field := &Field{
		name:     "test_field",
		baseName: "test_field",
		width:    width,
		height:   height,
		rawMap:   rawMap,
	}

	// Create a cell grid
	grid := NewCellGrid(field)
	if grid == nil {
		t.Fatal("Expected grid to be created, got nil")
	}

	if grid.Width() != width {
		t.Errorf("Expected width to be %d, got %d", width, grid.Width())
	}

	if grid.Height() != height {
		t.Errorf("Expected height to be %d, got %d", height, grid.Height())
	}

	// Test getting cells
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			cell := grid.GetCell(x, y)
			if cell == nil {
				t.Errorf("Expected cell at (%d, %d) to exist, got nil", x, y)
				continue
			}

			if cell.X() != x {
				t.Errorf("Expected X to be %d, got %d", x, cell.X())
			}

			if cell.Y() != y {
				t.Errorf("Expected Y to be %d, got %d", y, cell.Y())
			}
		}
	}
}

// TestLineOfSight tests the line of sight functionality
func TestLineOfSight(t *testing.T) {
	// Create a simple test field
	width := 10
	height := 10
	rawMap := make([]byte, width*height)

	// Make all cells walkable
	for i := 0; i < width*height; i++ {
		rawMap[i] = TileWalk
	}

	// Add a wall in the middle
	for y := 3; y < 7; y++ {
		rawMap[y*width+5] = 0 // Not walkable
	}

	field := &Field{
		name:     "test_field",
		baseName: "test_field",
		width:    width,
		height:   height,
		rawMap:   rawMap,
	}

	// Test line of sight
	testCases := []struct {
		fromX, fromY, toX, toY int
		canSnipe               bool
		expected               bool
	}{
		{1, 5, 8, 5, false, false}, // Blocked by wall
		{1, 1, 8, 8, false, true},  // Diagonal, not blocked
		{1, 5, 8, 5, true, false},  // Blocked by wall, even with snipe
		{5, 1, 5, 8, false, false}, // Vertical, blocked by wall
	}

	for _, tc := range testCases {
		from := Position{X: tc.fromX, Y: tc.fromY}
		to := Position{X: tc.toX, Y: tc.toY}
		result := field.CheckLOS(from, to, tc.canSnipe)
		if result != tc.expected {
			t.Errorf("CheckLOS(%v, %v, %v) = %v, expected %v", from, to, tc.canSnipe, result, tc.expected)
		}
	}
}

// TestClosestWalkableSpot tests finding the closest walkable spot
func TestClosestWalkableSpot(t *testing.T) {
	// Create a simple test field
	width := 10
	height := 10
	rawMap := make([]byte, width*height)

	// Make all cells walkable
	for i := 0; i < width*height; i++ {
		rawMap[i] = TileWalk
	}

	// Make a 3x3 non-walkable area
	for y := 4; y <= 6; y++ {
		for x := 4; x <= 6; x++ {
			rawMap[y*width+x] = 0 // Not walkable
		}
	}

	field := &Field{
		name:     "test_field",
		baseName: "test_field",
		width:    width,
		height:   height,
		rawMap:   rawMap,
	}

	// Test finding closest walkable spot
	testCases := []struct {
		x, y         int
		maxDistance  int
		expectedX    int
		expectedY    int
		expectResult bool
	}{
		{5, 5, 2, 3, 5, true},  // Center of non-walkable area, should find edge
		{5, 5, 0, 0, 0, false}, // Center of non-walkable area, no distance allowed
		{2, 2, 5, 2, 2, true},  // Already walkable spot
		{9, 9, 1, 9, 9, true},  // Edge of map, already walkable
	}

	for _, tc := range testCases {
		pos := Position{X: tc.x, Y: tc.y}
		result := field.ClosestWalkableSpot(pos, tc.maxDistance)

		if tc.expectResult {
			if result == nil {
				t.Errorf("ClosestWalkableSpot(%v, %d) returned nil, expected result", pos, tc.maxDistance)
				continue
			}

			if result.X != tc.expectedX || result.Y != tc.expectedY {
				t.Errorf("ClosestWalkableSpot(%v, %d) = (%d, %d), expected (%d, %d)",
					pos, tc.maxDistance, result.X, result.Y, tc.expectedX, tc.expectedY)
			}
		} else {
			if result != nil {
				t.Errorf("ClosestWalkableSpot(%v, %d) = %v, expected nil", pos, tc.maxDistance, result)
			}
		}
	}
}

// TestPathFinding tests the A* pathfinding algorithm
func TestPathFinding(t *testing.T) {
	// Create a simple test field
	width := 10
	height := 10
	rawMap := make([]byte, width*height)

	// Make all cells walkable
	for i := 0; i < width*height; i++ {
		rawMap[i] = TileWalk
	}

	// Add a wall with a gap
	for y := 3; y < 8; y++ {
		if y != 5 {
			rawMap[y*width+5] = 0 // Not walkable
		}
	}

	field := &Field{
		name:     "test_field",
		baseName: "test_field",
		width:    width,
		height:   height,
		rawMap:   rawMap,
	}

	// Create a cell grid and pathfinder
	grid := NewCellGrid(field)
	pathFinder := NewPathFinder(grid)

	// Test pathfinding
	testCases := []struct {
		fromX, fromY, toX, toY int
		expectedLength         int // Expected path length, or 0 if no path expected
	}{
		{1, 5, 8, 5, 9},   // Path through the gap
		{1, 1, 8, 8, 8},   // Diagonal path
		{1, 1, 1, 1, 1},   // Same start and end
		{1, 1, -1, -1, 0}, // Invalid end position
	}

	for _, tc := range testCases {
		from := Position{X: tc.fromX, Y: tc.fromY}
		to := Position{X: tc.toX, Y: tc.toY}
		path := pathFinder.FindPath(from, to)

		if tc.expectedLength > 0 {
			if path == nil {
				t.Errorf("FindPath(%v, %v) returned nil, expected path", from, to)
				continue
			}

			if path.Length() != tc.expectedLength {
				t.Errorf("FindPath(%v, %v) path length = %d, expected %d",
					from, to, path.Length(), tc.expectedLength)
			}

			// Check that the path starts and ends at the right places
			if path.Start().X() != tc.fromX || path.Start().Y() != tc.fromY {
				t.Errorf("Path starts at (%d, %d), expected (%d, %d)",
					path.Start().X(), path.Start().Y(), tc.fromX, tc.fromY)
			}

			if path.End().X() != tc.toX || path.End().Y() != tc.toY {
				t.Errorf("Path ends at (%d, %d), expected (%d, %d)",
					path.End().X(), path.End().Y(), tc.toX, tc.toY)
			}
		} else {
			if path != nil && path.Length() > 0 {
				t.Errorf("FindPath(%v, %v) returned path of length %d, expected no path", from, to, path.Length())
			}
		}
	}
}
