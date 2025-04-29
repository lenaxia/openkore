// Package actor provides actor-related packet sending functionality.
package actor

// ParseActorMove parses the coordinates from a packet.
// This is equivalent to the parse_actor_move function in Send.pm.
func ParseActorMove(args map[string]interface{}) {
	// Call MakeCoordsDir to parse the coordinates
	MakeCoordsDir(args)
}

// ReconstructActorMove reconstructs the coordinates for a packet.
// This is equivalent to the reconstruct_actor_move function in Send.pm.
func ReconstructActorMove(args map[string]interface{}) {
	// Check if no_padding is set
	noPadding, ok := args["no_padding"].(bool)
	if !ok {
		// Default to false if not set
		noPadding = false
	}

	// Get the x and y coordinates
	x, ok := args["x"].(int)
	if !ok {
		x = 0
	}

	y, ok := args["y"].(int)
	if !ok {
		y = 0
	}

	// Generate the coordinate string
	args["coords"] = GetCoordString(x, y, noPadding)
}

// GetCoordString generates a coordinate string from x and y coordinates.
// This is equivalent to the getCoordString function in Utils.pm.
func GetCoordString(x, y int, noPadding bool) []byte {
	// Convert x and y to bytes
	xByte := byte(x & 0xFF)
	yByte := byte(y & 0xFF)

	// Calculate the direction
	dir := 0

	// Create the coordinate string
	coords := []byte{xByte, yByte, byte(dir)}

	return coords
}

// MakeCoordsDir parses the coordinates from a packet and sets the x, y, and dir values.
// This is equivalent to the makeCoordsDir function in Utils.pm.
func MakeCoordsDir(args map[string]interface{}) {
	// Get the coordinate bytes
	coords, ok := args["coords"].([]byte)
	if !ok || len(coords) < 3 {
		// Default values if coords is not valid
		args["x"] = 0
		args["y"] = 0
		args["dir"] = 0
		return
	}

	// Parse the coordinates
	args["x"] = int(coords[0])
	args["y"] = int(coords[1])
	args["dir"] = int(coords[2])
}
