// Package game provides handlers for game-related packets.
package game

import (
	"fmt"

	"github.com/lenaxia/goKore/network/receive/types"
)

// RegisterActorHandlers registers all actor-related handlers with the receive component.
func RegisterActorHandlers(receive types.Receive) {
	receive.RegisterHandler("actor_move", handleActorMove)
	receive.RegisterHandler("actor_info", handleActorInfo)
	receive.RegisterHandler("actor_status", handleActorStatus)
	// More actor handlers would be registered here
}

// handleActorMove handles the actor_move packet.
func handleActorMove(args map[string]interface{}) error {
	// Extract fields from args
	actorID, ok := args["actorID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid actorID")
	}

	fromX, ok := args["fromX"].(uint16)
	if !ok {
		return fmt.Errorf("invalid fromX")
	}

	fromY, ok := args["fromY"].(uint16)
	if !ok {
		return fmt.Errorf("invalid fromY")
	}

	toX, ok := args["toX"].(uint16)
	if !ok {
		return fmt.Errorf("invalid toX")
	}

	toY, ok := args["toY"].(uint16)
	if !ok {
		return fmt.Errorf("invalid toY")
	}

	// Process the actor move
	fmt.Printf("Actor %d moved from (%d,%d) to (%d,%d)\n", actorID, fromX, fromY, toX, toY)

	// In a real implementation, this would update the game state
	// and trigger appropriate actions

	return nil
}

// handleActorInfo handles the actor_info packet.
func handleActorInfo(args map[string]interface{}) error {
	// Extract fields from args
	actorID, ok := args["actorID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid actorID")
	}

	name, ok := args["name"].(string)
	if !ok {
		return fmt.Errorf("invalid name")
	}

	// Process the actor info
	fmt.Printf("Actor info: ID=%d, name=%s\n", actorID, name)

	// In a real implementation, this would update the game state
	// and trigger appropriate actions

	return nil
}

// handleActorStatus handles the actor_status packet.
func handleActorStatus(args map[string]interface{}) error {
	// Extract fields from args
	actorID, ok := args["actorID"].(uint32)
	if !ok {
		return fmt.Errorf("invalid actorID")
	}

	status, ok := args["status"].(uint32)
	if !ok {
		return fmt.Errorf("invalid status")
	}

	// Process the actor status
	fmt.Printf("Actor %d status changed to %d\n", actorID, status)

	// In a real implementation, this would update the game state
	// and trigger appropriate actions

	return nil
}
