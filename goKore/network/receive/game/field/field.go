package field

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Tile types
const (
	TileNoWalk = 0
	TileWalk   = 1
	TileSnipe  = 2
	TileWater  = 4
	TileCliff  = 8
)

// Field represents a game map
type Field struct {
	// Name of the field (e.g., "prontera", "new_1-2", "0021@cata")
	name string

	// Base name of the field (e.g., "prontera", "new_1-2", "1@cata")
	baseName string

	// Instance ID (e.g., "002" in "0021@cata")
	instanceID string

	// Width of the field in cells
	width int

	// Height of the field in cells
	height int

	// Raw map data containing information about walkable/non-walkable cells
	rawMap []byte

	// Weight map data used for pathfinding
	weightMap []byte
}

// SetName sets the field's name
func (f *Field) SetName(name string) {
	f.name = name
}

// SetBaseName sets the field's base name
func (f *Field) SetBaseName(baseName string) {
	f.baseName = baseName
}

// SetInstanceID sets the field's instance ID
func (f *Field) SetInstanceID(instanceID string) {
	f.instanceID = instanceID
}

// SetWidth sets the field's width
func (f *Field) SetWidth(width int) {
	f.width = width
}

// SetHeight sets the field's height
func (f *Field) SetHeight(height int) {
	f.height = height
}

// SetRawMap sets the field's raw map data
func (f *Field) SetRawMap(rawMap []byte) {
	f.rawMap = rawMap
}

// SetWeightMap sets the field's weight map data
func (f *Field) SetWeightMap(weightMap []byte) {
	f.weightMap = weightMap
}

// New creates a new Field instance
func New(options map[string]interface{}) (*Field, error) {
	field := &Field{}

	if filename, ok := options["file"].(string); ok {
		if err := field.LoadFile(filename, true); err != nil {
			return nil, err
		}
	} else if name, ok := options["name"].(string); ok {
		if err := field.LoadByName(name, true); err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("no field name or filename specified")
	}

	return field, nil
}

// Name returns the field's name
func (f *Field) Name() string {
	return f.name
}

// BaseName returns the field's base name
func (f *Field) BaseName() string {
	return f.baseName
}

// InstanceID returns the field's instance ID
func (f *Field) InstanceID() string {
	return f.instanceID
}

// Width returns the field's width in cells
func (f *Field) Width() int {
	return f.width
}

// Height returns the field's height in cells
func (f *Field) Height() int {
	return f.height
}

// GetOffset calculates the offset in the map data for a given x,y coordinate
func (f *Field) GetOffset(x, y int) int {
	return (y * f.width) + x
}

// IsOffMap checks if coordinates are outside the map boundaries
func (f *Field) IsOffMap(x, y int) bool {
	return x < 0 || x >= f.width || y < 0 || y >= f.height
}

// GetBlock returns the block type at the specified coordinates
func (f *Field) GetBlock(x, y int) byte {
	if f.IsOffMap(x, y) {
		return 0
	}
	offset := f.GetOffset(x, y)
	return f.rawMap[offset]
}

// IsWalkable checks if the cell at (x,y) is walkable
func (f *Field) IsWalkable(x, y int) bool {
	if f.IsOffMap(x, y) {
		return false
	}
	return (f.GetBlock(x, y) & TileWalk) != 0
}

// IsSnipable checks if the cell at (x,y) is snipable
func (f *Field) IsSnipable(x, y int) bool {
	if f.IsOffMap(x, y) {
		return false
	}
	return (f.GetBlock(x, y) & TileSnipe) != 0
}

// IsWater checks if the cell at (x,y) is water
func (f *Field) IsWater(x, y int) bool {
	if f.IsOffMap(x, y) {
		return false
	}
	return (f.GetBlock(x, y) & TileWater) != 0
}

// IsCliff checks if the cell at (x,y) is a cliff
func (f *Field) IsCliff(x, y int) bool {
	if f.IsOffMap(x, y) {
		return false
	}
	return (f.GetBlock(x, y) & TileCliff) != 0
}

// GetBlockWeight returns the weight of the block at (x,y)
func (f *Field) GetBlockWeight(x, y int) byte {
	if f.IsOffMap(x, y) || f.weightMap == nil {
		return 0
	}
	offset := f.GetOffset(x, y)
	return f.weightMap[offset]
}

// GetCellInfo returns information about the cell at (x,y)
func (f *Field) GetCellInfo(x, y int) string {
	if f.IsOffMap(x, y) {
		return fmt.Sprintf("Cell %d %d is off the map.", x, y)
	}

	var info strings.Builder

	if f.IsWalkable(x, y) {
		info.WriteString(fmt.Sprintf("Cell %d %d is walkable.\n", x, y))
		weight := f.GetBlockWeight(x, y)
		info.WriteString(fmt.Sprintf("Cell %d %d has weight %d.\n", x, y, weight))
	} else {
		info.WriteString(fmt.Sprintf("Cell %d %d is not walkable.\n", x, y))
	}

	if f.IsSnipable(x, y) {
		info.WriteString(fmt.Sprintf("Cell %d %d is snipable.\n", x, y))
	} else {
		info.WriteString(fmt.Sprintf("Cell %d %d is not snipable.\n", x, y))
	}

	if f.IsWater(x, y) {
		info.WriteString(fmt.Sprintf("Cell %d %d is water.\n", x, y))
	} else {
		info.WriteString(fmt.Sprintf("Cell %d %d is not water.\n", x, y))
	}

	if f.IsCliff(x, y) {
		info.WriteString(fmt.Sprintf("Cell %d %d is a Cliff.\n", x, y))
	} else {
		info.WriteString(fmt.Sprintf("Cell %d %d is not a Cliff.\n", x, y))
	}

	return info.String()
}

// LoadFile loads a field file (.fld2)
func (f *Field) LoadFile(filename string, loadWeightMap bool) error {
	var fieldData []byte
	var err error

	// Check if the file is gzipped
	if strings.HasSuffix(filename, ".gz") {
		file, err := os.Open(filename)
		if err != nil {
			return fmt.Errorf("cannot open %s for reading: %w", filename, err)
		}
		defer file.Close()

		gz, err := gzip.NewReader(file)
		if err != nil {
			return fmt.Errorf("error creating gzip reader for %s: %w", filename, err)
		}
		defer gz.Close()

		fieldData, err = io.ReadAll(gz)
		if err != nil {
			return fmt.Errorf("error decompressing %s: %w", filename, err)
		}
	} else {
		fieldData, err = os.ReadFile(filename)
		if err != nil {
			return fmt.Errorf("cannot open %s for reading: %w", filename, err)
		}
	}

	// Extract width and height from the first 4 bytes
	if len(fieldData) < 4 {
		return errors.New("field file is too small")
	}

	f.width = int(binary.LittleEndian.Uint16(fieldData[0:2]))
	f.height = int(binary.LittleEndian.Uint16(fieldData[2:4]))
	f.rawMap = fieldData[4:]

	// Load the associated weight map if requested
	if loadWeightMap {
		weightFile := strings.TrimSuffix(filename, filepath.Ext(filename))
		weightFile = strings.TrimSuffix(weightFile, ".fld2")
		weightFile += ".weight"

		// Try to load the weight map
		if err := f.loadWeightMap(weightFile, f.width, f.height); err != nil {
			// If weight map loading fails, we should generate it
			// For now, we'll just create an empty weight map
			f.weightMap = make([]byte, f.width*f.height)
		}
	}

	// Set the base name from the filename
	_, file := filepath.Split(filename)
	f.baseName = strings.TrimSuffix(file, filepath.Ext(file))
	f.baseName = strings.TrimSuffix(f.baseName, ".fld2")
	f.name = f.baseName

	return nil
}

// loadWeightMap loads a weight map file (.weight)
func (f *Field) loadWeightMap(filename string, width, height int) error {
	var weightData []byte
	var err error

	// Check if the file exists with .gz extension
	gzFilename := filename + ".gz"
	if _, err := os.Stat(gzFilename); err == nil {
		filename = gzFilename
	}

	// Check if the file is gzipped
	if strings.HasSuffix(filename, ".gz") {
		file, err := os.Open(filename)
		if err != nil {
			return fmt.Errorf("cannot open %s for reading: %w", filename, err)
		}
		defer file.Close()

		gz, err := gzip.NewReader(file)
		if err != nil {
			return fmt.Errorf("error creating gzip reader for %s: %w", filename, err)
		}
		defer gz.Close()

		weightData, err = io.ReadAll(gz)
		if err != nil {
			return fmt.Errorf("error decompressing %s: %w", filename, err)
		}
	} else {
		weightData, err = os.ReadFile(filename)
		if err != nil {
			return fmt.Errorf("cannot open %s for reading: %w", filename, err)
		}
	}

	// Check file version
	if len(weightData) < 8 {
		return errors.New("weight map file is too small")
	}

	var version uint16 = 0
	if bytes.Equal(weightData[0:2], []byte("V#")) {
		version = binary.LittleEndian.Uint16(weightData[2:4])
		weightData = weightData[4:]
	}

	// Get map width and height from the weight file
	fileWidth := int(binary.LittleEndian.Uint16(weightData[0:2]))
	fileHeight := int(binary.LittleEndian.Uint16(weightData[2:4]))
	weightData = weightData[4:]

	// Version 1 is the current version
	if version >= 1 && width == fileWidth && height == fileHeight {
		f.weightMap = weightData
		return nil
	}

	return errors.New("invalid weight map file")
}

// LoadByName loads a field by its name
func (f *Field) LoadByName(name string, loadWeightMap bool) error {
	// Extract instance ID if present
	baseName, instanceID := f.nameToBaseName(name)
	f.baseName = baseName
	f.instanceID = instanceID

	// Determine the file path
	// In a real implementation, this would use configuration settings
	// For now, we'll just use a simple approach
	filename := f.baseName + ".fld2"

	// Check if the file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// Try with .gz extension
		filename += ".gz"
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			return fmt.Errorf("no corresponding field file found for field '%s'", name)
		}
	}

	// Load the file
	if err := f.LoadFile(filename, loadWeightMap); err != nil {
		return err
	}

	// Set the name
	f.name = name

	return nil
}

// nameToBaseName extracts the base name and instance ID from a field name
func (f *Field) nameToBaseName(name string) (string, string) {
	// Check for instance maps like "0021@cata"
	var instanceID string

	// Simple regex-like check for instance maps
	if len(name) >= 5 && name[3] == '@' {
		instanceID = name[0:3]
		name = name[3:]
	}

	return name, instanceID
}

// Position represents a 2D position in the game world
type Position struct {
	X int
	Y int
}

// ClosestWalkableSpot finds the closest walkable spot to the given position
func (f *Field) ClosestWalkableSpot(pos Position, maxDistance int) *Position {
	// If the position is already walkable, return it
	if f.IsWalkable(pos.X, pos.Y) {
		return &Position{X: pos.X, Y: pos.Y}
	}

	// If no max distance is specified, return nil
	if maxDistance <= 0 {
		return nil
	}

	// Check in increasing distance
	for distance := 1; distance <= maxDistance; distance++ {
		// Check all cells at this distance
		blocks := f.calcRectArea(pos.X, pos.Y, distance)
		for _, block := range blocks {
			if f.IsWalkable(block.X, block.Y) {
				return &block
			}
		}
	}

	return nil
}

// calcRectArea calculates a rectangle area around a point
func (f *Field) calcRectArea(x, y, radius int) []Position {
	var result []Position

	// Calculate the bounds of the rectangle
	minX := x - radius
	if minX < 0 {
		minX = 0
	}

	maxX := x + radius
	if maxX >= f.width {
		maxX = f.width - 1
	}

	minY := y - radius
	if minY < 0 {
		minY = 0
	}

	maxY := y + radius
	if maxY >= f.height {
		maxY = f.height - 1
	}

	// Add the top and bottom edges
	for i := minX; i <= maxX; i++ {
		// Top edge
		if minY >= 0 && f.IsWalkable(i, minY) {
			result = append(result, Position{X: i, Y: minY})
		}

		// Bottom edge
		if maxY < f.height && minY != maxY && f.IsWalkable(i, maxY) {
			result = append(result, Position{X: i, Y: maxY})
		}
	}

	// Add the left and right edges
	for i := minY + 1; i < maxY; i++ {
		// Left edge
		if minX >= 0 && f.IsWalkable(minX, i) {
			result = append(result, Position{X: minX, Y: i})
		}

		// Right edge
		if maxX < f.width && minX != maxX && f.IsWalkable(maxX, i) {
			result = append(result, Position{X: maxX, Y: i})
		}
	}

	return result
}

// CheckLOS checks if there is a line of sight between two positions
func (f *Field) CheckLOS(from, to Position, canSnipe bool) bool {
	// Implementation of Bresenham's line algorithm
	x0, y0 := from.X, from.Y
	x1, y1 := to.X, to.Y

	dx := abs(x1 - x0)
	dy := abs(y1 - y0)

	var sx, sy int
	if x0 < x1 {
		sx = 1
	} else {
		sx = -1
	}

	if y0 < y1 {
		sy = 1
	} else {
		sy = -1
	}

	err := dx - dy

	for {
		// Check if the current cell is walkable or snipable
		if canSnipe {
			if !f.IsWalkable(x0, y0) && !f.IsSnipable(x0, y0) {
				return false
			}
		} else {
			if !f.IsWalkable(x0, y0) {
				return false
			}
		}

		if x0 == x1 && y0 == y1 {
			break
		}

		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x0 += sx
		}

		if e2 < dx {
			err += dx
			y0 += sy
		}
	}

	return true
}

// Helper function for absolute value
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
