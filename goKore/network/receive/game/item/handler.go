package item

import (
	"github.com/lenaxia/goKore/network/hooks"
	"github.com/lenaxia/goKore/network/receive/core"
)

// Handler is the main handler for item-related packets
type Handler struct {
	inventoryManager *InventoryManager
	storageManager   *StorageManager
	equipmentManager *EquipmentManager
	itemUsageManager *ItemUsageManager
	cartManager      *CartManager
	logger           core.Logger
}

// NewHandler creates a new item handler
func NewHandler(baseParse Parser, hooks *hooks.HookManager, logger core.Logger) *Handler {
	return &Handler{
		inventoryManager: NewInventoryManager(baseParse, hooks, logger),
		storageManager:   NewStorageManager(baseParse, hooks, logger),
		equipmentManager: NewEquipmentManager(baseParse, hooks, logger),
		itemUsageManager: NewItemUsageManager(baseParse, hooks, logger),
		cartManager:      NewCartManager(baseParse, hooks, logger),
		logger:           logger,
	}
}

// RegisterHandlers registers all item-related packet handlers
func (h *Handler) RegisterHandlers() {
	// Register inventory handlers
	h.inventoryManager.RegisterHandlers()

	// Register storage handlers
	h.storageManager.RegisterHandlers()

	// Register equipment handlers
	h.equipmentManager.RegisterHandlers()

	// Register item usage handlers
	h.itemUsageManager.RegisterHandlers()

	// Register cart handlers
	h.cartManager.RegisterHandlers()
}

// GetInventoryManager returns the inventory manager
func (h *Handler) GetInventoryManager() *InventoryManager {
	return h.inventoryManager
}

// GetStorageManager returns the storage manager
func (h *Handler) GetStorageManager() *StorageManager {
	return h.storageManager
}

// GetEquipmentManager returns the equipment manager
func (h *Handler) GetEquipmentManager() *EquipmentManager {
	return h.equipmentManager
}

// GetCartManager returns the cart manager
func (h *Handler) GetCartManager() *CartManager {
	return h.cartManager
}

// GetItemUsageManager returns the item usage manager
func (h *Handler) GetItemUsageManager() *ItemUsageManager {
	return h.itemUsageManager
}
