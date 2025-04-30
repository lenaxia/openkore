package minimap

import (
	"fmt"

	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// Color represents an RGBA color
type Color struct {
	Red   uint8
	Green uint8
	Blue  uint8
	Alpha uint8
}

// String returns a string representation of the color
func (c Color) String() string {
	return fmt.Sprintf("[R:%d, G:%d, B:%d, A:%d]", c.Red, c.Green, c.Blue, c.Alpha)
}

// MinimapIndicator represents a minimap indicator
type MinimapIndicator struct {
	NpcID   uint32
	X       uint16
	Y       uint16
	Type    uint8
	Effect  uint16
	Color   Color
	Show    bool
	ActorID uint32
}

// MinimapManager handles minimap-related packet handling
type MinimapManager struct {
	parser      *core.CoreParser
	hookManager *hooks.HookManager
	logger      core.Logger

	// Map of quest types to colors
	questTypeColors map[uint8]Color
}

// NewMinimapManager creates a new minimap manager
func NewMinimapManager(parser *core.CoreParser, hookManager *hooks.HookManager, logger core.Logger) *MinimapManager {
	// Initialize quest type colors
	questTypeColors := make(map[uint8]Color)
	questTypeColors[0] = Color{255, 255, 255, 0} // Default color

	return &MinimapManager{
		parser:          parser,
		hookManager:     hookManager,
		logger:          logger,
		questTypeColors: questTypeColors,
	}
}

// RegisterHandlers registers all minimap-related packet handlers
func (mm *MinimapManager) RegisterHandlers() {
	// Register minimap indicator handler
	mm.parser.RegisterHandlerFunc("0144", "minimap_indicator", "V v2 B2 B4",
		[]string{"npcID", "x", "y", "type", "effect", "red", "green", "blue", "alpha"}, mm.HandleMinimapIndicator)
}

// GetActorName gets the name of an actor by its ID
func (mm *MinimapManager) GetActorName(actorID uint32) string {
	// In a real implementation, this would look up the actor name from the actor list
	// For now, we'll just return a placeholder
	return fmt.Sprintf("Actor#%d", actorID)
}

// HandleMinimapIndicator handles the minimap_indicator packet (lines 3071-3099)
func (mm *MinimapManager) HandleMinimapIndicator(args map[string]interface{}) error {
	// Extract packet data
	npcID, ok := args["npcID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid npcID in minimap_indicator packet")
	}

	x, ok := args["x"].(uint16)
	if !ok {
		return fmt.Errorf("invalid x in minimap_indicator packet")
	}

	y, ok := args["y"].(uint16)
	if !ok {
		return fmt.Errorf("invalid y in minimap_indicator packet")
	}

	// Type is optional
	var indicatorType uint8
	if typeVal, ok := args["type"].(uint8); ok {
		indicatorType = typeVal
	}

	// Effect is optional
	var effect uint16
	if effectVal, ok := args["effect"].(uint16); ok {
		effect = effectVal
	}

	// Extract color data
	red, ok := args["red"].(uint8)
	if !ok {
		red = 255 // Default to white
	}

	green, ok := args["green"].(uint8)
	if !ok {
		green = 255 // Default to white
	}

	blue, ok := args["blue"].(uint8)
	if !ok {
		blue = 255 // Default to white
	}

	alpha, ok := args["alpha"].(uint8)
	if !ok {
		alpha = 0 // Default to transparent
	}

	// Create color
	color := Color{
		Red:   red,
		Green: green,
		Blue:  blue,
		Alpha: alpha,
	}

	// Determine if the indicator is being shown or cleared
	show := indicatorType != 2

	// Get actor name
	actorName := mm.GetActorName(npcID)

	// We'll use these values directly in the hook call

	// Determine indicator type string
	indicatorTypeStr := "minimap indicator"
	if indicatorType != 0 && indicatorType != 1 && indicatorType != 2 {
		indicatorTypeStr = fmt.Sprintf("minimap indicator (unknown type %d)", indicatorType)
	} else if effect == 1 {
		indicatorTypeStr = "*Quest!*"
	} else if effect == 9999 {
		// Special case, don't log anything
		return nil
	} else if effect > 0 {
		indicatorTypeStr = fmt.Sprintf("unknown effect %d", effect)
	}

	// Log message
	if show {
		mm.logger.Info("%s shown %s at location %d, %d with the color %s",
			actorName, indicatorTypeStr, x, y, color.String())
	} else {
		mm.logger.Info("%s cleared %s at location %d, %d with the color %s",
			actorName, indicatorTypeStr, x, y, color.String())
	}

	// Call hook
	mm.hookManager.CallHook("minimap_indicator", map[string]interface{}{
		"npcID":   npcID,
		"x":       x,
		"y":       y,
		"type":    indicatorType,
		"effect":  effect,
		"color":   color,
		"show":    show,
		"actorID": npcID,
	})

	return nil
}
