**Auction System Handlers:**
- auction_list (lines 6277-6290)
  - Manages auction item listings
  - Processes:
    - Item IDs
    - Seller information
    - Price data
    - Time remaining
  - Maintains:
    - @auctionItems array
    - %auctionDetails hash

- auction_search_result (lines 6292-6300)
  - Handles auction search results

- auction_bid_result (lines 6302-6315)
  - Processes bid outcomes
  - Handles:
    - Success/failure status
    - Bid amount validation

- auction_purchase_result (lines 6317-6330)
  - Processes instant purchase outcomes
  - Handles:
    - Purchase status codes
    - Item transfer verification
  - Features:
    - Updates inventory on success
    - Triggers 'purchase_complete' hook
    - Provides transaction details
    - Current highest bid
  - Features:
    - Updates auction status
    - Triggers 'bid_response' hook
    - Provides user feedback
  - Features:
    - Processes search criteria matches
    - Filters and sorts results
    - Updates auction display
    - Triggers 'auction_results' hook