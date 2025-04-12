package entity

type EquipSlot int

const (
	SlotHeadTop EquipSlot = iota + 1
	SlotHeadMid
	SlotHeadLow
	SlotBody
	SlotRightHand
	SlotLeftHand
	// ... other slots from equipSlot_lut
)

func (s EquipSlot) IsValid() bool {
	return s >= SlotHeadTop && s <= SlotLastValid
}

// Slot mapping from Globals.pm %equipSlot_lut
var slotNames = map[EquipSlot]string{
	SlotHeadTop:    "Top Headgear",
	SlotHeadMid:    "Mid Headgear",
	SlotHeadLow:    "Low Headgear",
	SlotBody:       "Body Armor",
	SlotRightHand:  "Right Hand",
	SlotLeftHand:   "Left Hand",
	// ... other slots
}

func SlotName(slot EquipSlot) string {
	name, ok := slotNames[slot]
	if !ok {
		return "Unknown"
	}
	return name
}
