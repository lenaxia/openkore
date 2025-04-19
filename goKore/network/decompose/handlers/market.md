**Market System Handlers:**

- sell_result() - Item selling results (lines 9519-9531)
  - Processes item selling results:
    * Failed sale (displays error)
    * Successful sale (displays count and success message)
  - Clears sell list
  - Updates AI state if in sellAuto mode

**Search Store/Vendor System:**
- search_store_open() - Opens vendor search (lines 9209-9219)
  - Initializes universal catalog
  - Tracks search type (gold/silver)
  - Displays remaining searches

- search_store_fail() - Failed search (lines 9221-9239)
  - Handles different failure cases:
    * No results (0)
    * Too many requests (1)
    * Invalid search (2)
    * Server error (3)
    * No permission (4)

- search_store_result() - Search results (lines 9241-9278)
  - Processes vendor listings:
    * Store/account IDs
    * Shop name
    * Item details (name, price, amount)
    * Refine level
    * Cards
  - Maintains pagination
  - Triggers search_store hook

- search_store_pos() - Vendor location (lines 9280-9284)
  - Displays vendor coordinates