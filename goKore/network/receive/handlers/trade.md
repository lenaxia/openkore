# Trade Related Handlers

**Method Implementations:**
- deal_add_you - Trade item addition (self) handler (lines 7828-7862)
  - Processes items added by player to trade
  - Handles various failure cases:
    * 1: Target overweight - "That person is overweight"
    * 2: Untradeable item - "This item cannot be traded"
    * 192: Unknown success status (debug message)
    * Other: Generic failure with code
  - For successful additions:
    * Gets item ID from packet
    * Returns early if ID <= 0
    * Gets item reference from inventory
    * Increments currentDeal{you_items} counter
    * Updates currentDeal{you} hash with item details
    * Displays message with item name and amount
    * Calls inventoryItemRemoved to update inventory
    * Triggers deal_you_added hook with item info
  - Contains FIXME comments about potential inventory corruption
  - Uses "deal" message category for all outputs
- deal_request - Trade request handler (lines 5913-5927)
  - Processes incoming trade requests
  - Extracts player information:
    * Level (or "Unknown" if not provided)
    * Name (converted from bytes to string)
  - Sets up incomingDeal with player name
  - Starts auto-cancel timeout
  - Displays request message with player name and level
  - Shows instructions for accepting/denying the deal
  - Triggers incoming_deal hook with:
    * Player name
    * Player level
    * Player ID
  - Uses "deal" message category for all outputs
- deal_finalize - Trade finalization handler (lines 5898-5911)
  - Processes trade finalization notifications
  - Handles two scenarios:
    * Other player finalizes (type 1):
      - Sets currentDeal{other_finalize} flag
      - Displays "[player] finalized the Deal" message
      - Triggers finalized_deal hook with player name
    * Current player finalizes (type != 1):
      - Sets currentDeal{you_finalize} flag
      - Subtracts trade zeny from character
      - Displays "You finalized the Deal" message
      - Contains FIXME comment about zeny subtraction timing
  - Uses "deal" message category for all outputs
- deal_complete - Trade completion handler (lines 5890-5896)
  - Processes trade completion notifications
  - Cleans up all trade-related data structures:
    * Clears outgoingDeal hash
    * Clears incomingDeal hash
    * Clears currentDeal hash
  - Displays "Deal Complete" message
  - Triggers complete_deal hook
  - Uses "deal" message category
  - Similar to deal_cancelled but with different message and hook
- deal_cancelled - Trade cancellation handler (lines 5882-5888)
  - Processes trade cancellation notifications
  - Cleans up all trade-related data structures:
    * Clears incomingDeal hash
    * Clears outgoingDeal hash
    * Clears currentDeal hash
  - Displays "Deal Cancelled" message
  - Triggers cancelled_deal hook
  - Uses "deal" message category
  - Simple implementation with minimal processing
- deal_begin - Trade initiation handler (lines 5846-5875)
  - Processes trade begin results
  - Handles different error types:
    * 0: Target too far away
    * 2: Target in another deal
    * 5: Target opening storage
  - Handles successful trade start (type 3):
    * For incoming deals: Sets name from incomingDeal
    * For outgoing deals:
      - Gets player by ID from playersList
      - Sets name from player or generates "Unknown #ID"
    * Clears incomingDeal/outgoingDeal as appropriate
    * Displays "Engaged Deal" message
    * Triggers engaged_deal hook
  - Triggers error_deal hook for error cases
  - Uses "deal" message category for all outputs
- deal_add_other - Trade item addition handler (lines 5825-5844)
  - Processes items added by other player to trade
  - Handles two scenarios:
    * Item addition (nameID > 0):
      - Creates or updates item in currentDeal{other} hash
      - Updates item properties (amount, nameID, identified, etc.)
      - Sets item name using itemName function
      - Displays message with item name and amount
    * Zeny addition (amount > 0):
      - Updates currentDeal{other_zeny} value
      - Formats zeny amount for display
      - Displays message with amount
  - Uses "deal" message category for all outputs