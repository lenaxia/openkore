// Package send provides functionality for sending packets to the server.
package send

import (
	"github.com/lenaxia/goKore/network/common"
	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/send/core"
	"github.com/lenaxia/goKore/network/send/types"
)

// Re-export types from the common package
type PacketConstruction = common.PacketConstruction

// Re-export types from the types package
type SendHandler = types.SendHandler
type Send = types.Send

// Re-export types from the core package
type BaseSend = core.BaseSend

// PacketConstructionProvider is a function that returns a map of packet constructions
type PacketConstructionProvider func() map[string]PacketConstruction

// NewBaseSend creates a new BaseSend instance.
func NewBaseSend(hookManager *hooks.HookManager) *BaseSend {
	return core.NewBaseSend(hookManager)
}
