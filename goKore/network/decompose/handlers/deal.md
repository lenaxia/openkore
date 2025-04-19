**Trade Deal Handlers:**

- deal_add_you (lines 7828-7862)
  - Handles adding items to trade deals
  - Features:
    - Processes various failure cases:
      - Overweight (code 1)
      - Untradeable item (code 2)
      - Other error codes
    - Updates deal state with item info
    - Maintains item counts
    - Triggers inventory updates
    - Calls plugin hooks for deal modifications