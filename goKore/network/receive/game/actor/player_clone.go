package actor

// Add clone field to Player struct
func init() {
	// This is just a placeholder to ensure the file is compiled
	// The actual clone field is added to the Player struct in player.go
}

// IsClone returns whether the player is a clone (offline shop)
func (p *Player) IsClone() bool {
	// We'll use a special field in the Player struct to track this
	// For now, we can check if the player has a specific flag set
	return p.clone
}

// SetClone sets whether the player is a clone (offline shop)
func (p *Player) SetClone(isClone bool) {
	p.clone = isClone
}
