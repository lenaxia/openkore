# Auction Related Handlers

**Method Implementations:**
- auction_result - Auction result handler (lines 10325-10352)
  - Processes auction result notifications
  - Handles multiple flag values:
    * 0: Bid failure - "You have failed to bid into the auction"
    * 1: Bid success - "You have successfully bid in the auction"
    * 2: Auction canceled - "The auction has been canceled"
    * 3: Cannot cancel - "An auction with at least one bidder cannot be canceled"
    * 4: Too many items - "You cannot register more than 5 items in an auction at a time"
    * 5: Insufficient zeny - "You do not have enough Zeny to pay the Auction Fee"
    * 6: Auction won - "You have won the auction"
    * 7: Auction lost - "You have failed to win the auction"
    * 8: Insufficient zeny - "You do not have enough Zeny"
    * 9: Too many bids - "You cannot place more than 5 bids at a time"
    * Other: Unknown result with flag value
  - Uses "info" message category for known results
  - Uses "warning" message category for unknown results
- auction_add_item - Auction item addition handler (lines 10264-10272)
  - Processes auction item addition notifications
  - Handles two fail values:
    * 1: Failure - "Failed (note: usable items can't be auctioned) to add item with index: X"
    * 0: Success - "Succeeded to add item with index: X"
  - Displays message with item ID
  - Uses "info" message category for all messages
  - Simple implementation focused on notification
- auction_windows - Auction window status handler (lines 10254-10262)
  - Processes auction window status notifications
  - Handles two flag values:
    * 1: Closed - "Auction window is now closed"
    * 0: Opened - "Auction window is now opened"
  - Uses "info" message category for all messages
  - Simple implementation focused on notification
- auction_my_sell_stop - Auction end handler (lines 10239-10252)
  - Processes auction end notifications
  - Handles multiple flag values:
    * 0: Success - "You have ended the auction"
    * 1: Failure - "You cannot end the auction"
    * 2: Invalid bid - "Bid number is incorrect"
    * Other: Unknown result with flag value
  - Uses "info" message category for known results
  - Uses "warning" message category for unknown results