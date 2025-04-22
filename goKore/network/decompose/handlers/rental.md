# Item Rental Related Handlers

**Method Implementations:**
- rental_time - Rental time notification handler (lines 10607-10610)
  - Processes rental time notifications
  - Gets item name using itemNameSimple function
  - Displays message with item name and remaining time in minutes
  - Uses "info" message category
  - Contains TODO comment: "can we use itemName($actor)?"
  - Contains comment: "don't think so because it seems that this packet is received before the inventory list"
  - Simple implementation focused on notification